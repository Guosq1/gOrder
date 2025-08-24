package main

import (
	"context"

	"github.com/Hypocrite/gorder/common/config"
	"github.com/Hypocrite/gorder/common/discovery"
	"github.com/Hypocrite/gorder/common/genproto/stockpb"
	"github.com/Hypocrite/gorder/common/logging"
	"github.com/Hypocrite/gorder/common/server"
	"github.com/Hypocrite/gorder/common/tracing"
	"github.com/Hypocrite/gorder/stock/ports"
	"github.com/Hypocrite/gorder/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	logging.Init()
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatal(err)
	}
}

func main() {

	serviceName := viper.GetString("stock.service-name")
	serverType := viper.GetString("stock.server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application := service.NewApplication(ctx)
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}

	defer func() {
		_ = deregisterFunc()
	}()

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer(application)
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
		//ToDo
	default:
		panic("unexpected server type: ")
	}

}

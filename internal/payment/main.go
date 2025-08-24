package main

import (
	"context"

	"github.com/Hypocrite/gorder/common/broker"
	_ "github.com/Hypocrite/gorder/common/config"
	"github.com/Hypocrite/gorder/common/logging"
	"github.com/Hypocrite/gorder/common/server"
	"github.com/Hypocrite/gorder/common/tracing"
	"github.com/Hypocrite/gorder/payment/infrastructure/consumer"
	"github.com/Hypocrite/gorder/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	logging.Init()

}

func main() {
	serviceName := viper.GetString("payment.service-name")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serverType := viper.GetString("payment.server-to-run")

	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer shutdown(ctx)

	application, cleanup := service.NewApplication(ctx)
	defer cleanup()

	ch, closeCH := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)
	defer func() {
		_ = ch.Close()
		_ = closeCH()
	}()
	go consumer.NewConsumer(application).Listen(ch)

	paymentHandler := NewPaymentHandler(ch)
	switch serverType {
	case "http":
		server.RunHTTpServer(serviceName, paymentHandler.RegisterRoutes)
	case "grpc":
		panic("unsupported type: grpc")
	default:
		logrus.Panic("unsupported server type")
	}

}

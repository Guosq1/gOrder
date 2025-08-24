package service

import (
	"context"

	"github.com/Hypocrite/gorder/common/broker"
	grpcClient "github.com/Hypocrite/gorder/common/client"
	"github.com/Hypocrite/gorder/common/metric"
	"github.com/Hypocrite/gorder/order/adapters"
	Grpc "github.com/Hypocrite/gorder/order/adapters/gprc"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"

	"github.com/Hypocrite/gorder/order/app"
	"github.com/Hypocrite/gorder/order/app/command"
	"github.com/Hypocrite/gorder/order/app/query"

	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	ch, closeCH := broker.Connect(
		viper.GetString("rabbitmq.user"),
		viper.GetString("rabbitmq.password"),
		viper.GetString("rabbitmq.host"),
		viper.GetString("rabbitmq.port"),
	)

	stockGRPC := Grpc.NewStockGPRC(stockClient)
	return newApplication(ctx, stockGRPC, ch), func() {
		_ = closeStockClient()
		_ = closeCH()
		_ = ch.Close()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
	orderRepo := adapters.NewMemoryOrderRepository()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metric.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logger, metricClient),
			UpdateOrder: command.NewUpdateOrderHandler(orderRepo, logger, metricClient),
		},
		Queries: app.Queries{
			GetCustomerOrder: query.NewGetCustomerOrderHandler(orderRepo, logger, metricClient),
		},
	}
}

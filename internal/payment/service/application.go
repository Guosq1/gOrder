package service

import (
	"context"

	grpcClient "github.com/Hypocrite/gorder/common/client"
	"github.com/Hypocrite/gorder/common/metric"
	"github.com/Hypocrite/gorder/payment/adapters"
	"github.com/Hypocrite/gorder/payment/app"
	"github.com/Hypocrite/gorder/payment/app/command"
	"github.com/Hypocrite/gorder/payment/domain"
	"github.com/Hypocrite/gorder/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	//memProcessor := processor.NewInmemProcessor()
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(ctx context.Context, orderGRPC command.OrderService, processor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricClient := metric.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(processor, orderGRPC, logger, metricClient),
		},
	}
}

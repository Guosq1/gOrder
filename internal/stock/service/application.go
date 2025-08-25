package service

import (
	"context"

	"github.com/Hypocrite/gorder/common/metric"
	"github.com/Hypocrite/gorder/stock/adapters"
	"github.com/Hypocrite/gorder/stock/app"
	"github.com/Hypocrite/gorder/stock/app/query"
	"github.com/Hypocrite/gorder/stock/infrastructure/integration"
	"github.com/Hypocrite/gorder/stock/infrastructure/persistent"
	"github.com/sirupsen/logrus"
)

func NewApplication(_ context.Context) app.Application {

	db := persistent.NewMySQL()
	stockRepo := adapters.NewMySQLStockRepository(db)
	logger := logrus.NewEntry(logrus.StandardLogger())
	stripeApi := integration.NewStripeApi()
	metricClient := metric.TodoMetrics{}
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			CheckIfItemsInStock: query.NewCheckIfItemsInStockHandler(stockRepo, stripeApi, logger, metricClient),
			GetItems:            query.NewGetItemsHandler(stockRepo, logger, metricClient),
		},
	}
}

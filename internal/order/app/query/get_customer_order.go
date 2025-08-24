package query

import (
	"context"

	"github.com/Hypocrite/gorder/common/decorator"
	"github.com/Hypocrite/gorder/common/tracing"
	domain "github.com/Hypocrite/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
)

type GetCustomerOrder struct {
	CustomerID string
	OrderID    string
}

type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

type getCustomerOrderHandler struct {
	orderRepo domain.Repository
}

func NewGetCustomerOrderHandler(
	orderRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyQueryDecorator[GetCustomerOrder, *domain.Order](
		getCustomerOrderHandler{orderRepo: orderRepo},
		logger,
		metricClient,
	)
}

func (g getCustomerOrderHandler) Handle(ctx context.Context, q GetCustomerOrder) (*domain.Order, error) {

	_, span := tracing.Start(ctx, "getCustomerOrderHandler.Handle")

	o, err := g.orderRepo.Get(ctx, q.OrderID, q.CustomerID)
	if err != nil {
		return nil, err
	}

	span.AddEvent("get customer order success")
	span.End()
	return o, nil

}

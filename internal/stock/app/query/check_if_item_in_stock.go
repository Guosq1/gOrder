package query

import (
	"context"

	"github.com/Hypocrite/gorder/common/decorator"
	"github.com/Hypocrite/gorder/common/genproto/orderpb"
	domain "github.com/Hypocrite/gorder/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*orderpb.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorator[CheckIfItemsInStock, []*orderpb.Item](
		checkIfItemsInStockHandler{
			stockRepo: stockRepo,
		},
		logger,
		metricClient,
	)
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, q CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var res []*orderpb.Item
	for _, item := range q.Items {
		res = append(res, &orderpb.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
		})
	}
	return res, nil
}

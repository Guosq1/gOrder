package query

import (
	"context"

	"github.com/Hypocrite/gorder/common/decorator"
	domain "github.com/Hypocrite/gorder/stock/domain/stock"
	"github.com/Hypocrite/gorder/stock/entity"
	"github.com/Hypocrite/gorder/stock/infrastructure/integration"
	"github.com/sirupsen/logrus"
)

type CheckIfItemsInStock struct {
	Items []*entity.ItemWithQuantity
}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

type checkIfItemsInStockHandler struct {
	stockRepo domain.Repository
	stripeApi *integration.StripeApi
}

func NewCheckIfItemsInStockHandler(
	stockRepo domain.Repository,
	stripeApi *integration.StripeApi,
	logger *logrus.Entry,
	metricClient decorator.MetricsClient,
) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	if stripeApi == nil {
		panic("stripeApi is nil")
	}
	return decorator.ApplyQueryDecorator[CheckIfItemsInStock, []*entity.Item](
		checkIfItemsInStockHandler{
			stockRepo: stockRepo,
			stripeApi: stripeApi,
		},
		logger,
		metricClient,
	)
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, q CheckIfItemsInStock) ([]*entity.Item, error) {
	var res []*entity.Item
	for _, item := range q.Items {
		priceID, err := c.stripeApi.GetPriceByProductID(ctx, item.ID)
		if err != nil || priceID == "" {
			return nil, err
		}
		res = append(res, &entity.Item{
			ID:       item.ID,
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}
	return res, nil
}

package query

import (
	"context"
	"strings"
	"time"

	"github.com/Hypocrite/gorder/common/decorator"
	"github.com/Hypocrite/gorder/common/handler/redis"
	domain "github.com/Hypocrite/gorder/stock/domain/stock"
	"github.com/Hypocrite/gorder/stock/entity"
	"github.com/Hypocrite/gorder/stock/infrastructure/integration"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	redisLockPrefix = "check_stock_"
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

	if err := lock(ctx, getLockKey(q)); err != nil {
		return nil, errors.Wrapf(err, "redis lock error: key=%s", getLockKey(q))
	}
	defer func() {
		if err := unlock(ctx, getLockKey(q)); err != nil {
			logrus.Warnf("redis fail to unlock , err=%v", err)
		}
	}()

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
	if err := c.checkStock(ctx, q.Items); err != nil {
		return nil, err
	}
	return res, nil
}

func getLockKey(query CheckIfItemsInStock) string {
	var ids []string
	for _, i := range query.Items {
		ids = append(ids, i.ID)
	}
	return redisLockPrefix + strings.Join(ids, "_")
}

func lock(ctx context.Context, key string) error {
	return redis.SetNX(ctx, redis.LocalClient(), key, "1", 5*time.Minute)
}

func unlock(ctx context.Context, key string) error {
	return redis.Del(ctx, redis.LocalClient(), key)
}

func (c checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
	var ids []string
	for _, item := range query {
		ids = append(ids, item.ID)
	}
	records, err := c.stockRepo.GetStock(ctx, ids)
	if err != nil {
		return err
	}
	idQuantity := make(map[string]int32)
	for _, record := range records {
		idQuantity[record.ID] += record.Quantity
	}

	var (
		ok       = true
		failedOn []struct {
			ID   string
			Want int32
			Have int32
		}
	)

	for _, record := range query {
		if record.Quantity > idQuantity[record.ID] {
			ok = false
			failedOn = append(failedOn, struct {
				ID   string
				Want int32
				Have int32
			}{ID: record.ID, Want: record.Quantity, Have: idQuantity[record.ID]})
		}
	}
	if ok {
		return c.stockRepo.UpdateStock(ctx, query, func(
			ctx context.Context,
			existing []*entity.ItemWithQuantity,
			query []*entity.ItemWithQuantity,
		) ([]*entity.ItemWithQuantity, error) {
			var newItems []*entity.ItemWithQuantity
			for _, e := range existing {
				for _, q := range query {
					if e.ID == q.ID {
						newItems = append(newItems, &entity.ItemWithQuantity{
							ID:       e.ID,
							Quantity: e.Quantity - q.Quantity,
						})
					}
				}
			}
			return newItems, nil
		})
	}
	return domain.ExceedStockError{FailedOn: failedOn}
}

package query

import (
	"context"

	"github.com/Hypocrite/gorder/common/genproto/orderpb"
	"github.com/Hypocrite/gorder/common/genproto/stockpb"
)

type StockService interface {
	CheckIfItemInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemInStockRes, error)
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}

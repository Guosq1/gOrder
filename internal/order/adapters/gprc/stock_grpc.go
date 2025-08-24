package gprc

import (
	"context"

	"github.com/Hypocrite/gorder/common/genproto/orderpb"
	"github.com/Hypocrite/gorder/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

type StockGPRC struct {
	client stockpb.StockServiceClient
}

func NewStockGPRC(client stockpb.StockServiceClient) *StockGPRC {
	return &StockGPRC{client: client}
}

func (s StockGPRC) CheckIfItemInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemInStockRes, error) {
	resp, err := s.client.CheckIfItemInStock(ctx, &stockpb.CheckIfItemInStockReq{
		Items: items,
	})
	logrus.Info("stock_grpc response: ", resp)
	return resp, err
}

func (s StockGPRC) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsReq{
		ItemIDs: itemIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

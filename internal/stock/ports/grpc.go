package ports

import (
	"context"
	"github.com/Hypocrite/gorder/common/genproto/stockpb"
	"github.com/Hypocrite/gorder/stock/app"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, req *stockpb.GetItemsReq) (*stockpb.GetItemsRes, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) CheckIfItemInStock(ctx context.Context, req *stockpb.CheckIfItemInStockReq) (*stockpb.CheckIfItemInStockRes, error) {
	//TODO implement me
	panic("implement me")
}

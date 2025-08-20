package ports

import (
	"context"
	"github.com/Hypocrite/gorder/common/genproto/orderpb"
	"github.com/Hypocrite/gorder/order/app"
	"github.com/golang/protobuf/ptypes/empty"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderReq) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) GetOrder(ctx context.Context, req *orderpb.GetOrderReq) (*orderpb.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (G GRPCServer) UpdateOrder(ctx context.Context, order *orderpb.Order) (*empty.Empty, error) {
	//TODO implement me
	panic("implement me")
}

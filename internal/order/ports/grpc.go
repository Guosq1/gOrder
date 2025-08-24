package ports

import (
	"context"

	"github.com/Hypocrite/gorder/common/genproto/orderpb"
	"github.com/Hypocrite/gorder/order/app"
	"github.com/Hypocrite/gorder/order/app/command"
	"github.com/Hypocrite/gorder/order/app/query"
	"github.com/Hypocrite/gorder/order/convertor"
	domain "github.com/Hypocrite/gorder/order/domain/order"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, req *orderpb.CreateOrderReq) (*emptypb.Empty, error) {
	_, err := G.app.Commands.CreateOrder.Handle(ctx, command.CrateOrder{
		CustomerID: req.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ProtosToEntities(req.Items),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (G GRPCServer) GetOrder(ctx context.Context, req *orderpb.GetOrderReq) (*orderpb.Order, error) {
	o, err := G.app.Queries.GetCustomerOrder.Handle(ctx, query.GetCustomerOrder{
		CustomerID: req.CustomerID,
		OrderID:    req.OrderID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return convertor.NewOrderConvertor().EntityToProto(o), nil
}

func (G GRPCServer) UpdateOrder(ctx context.Context, req *orderpb.Order) (_ *emptypb.Empty, err error) {
	logrus.Infof("Received Order: %v", req)
	order, err := domain.NewOrder(
		req.ID,
		req.CustomerID,
		req.Status,
		req.PaymentLink,
		convertor.NewItemConvertor().ProtosToEntities(req.Items))
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		return
	}

	_, err = G.app.Commands.UpdateOrder.Handle(ctx, command.UpdateOrder{
		Order: order,
		UpdateFn: func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		},
	})

	return
}

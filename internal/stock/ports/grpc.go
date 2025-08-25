package ports

import (
	"context"

	"github.com/Hypocrite/gorder/common/genproto/stockpb"
	"github.com/Hypocrite/gorder/common/tracing"
	"github.com/Hypocrite/gorder/stock/app"
	"github.com/Hypocrite/gorder/stock/app/query"
	"github.com/Hypocrite/gorder/stock/convertor"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, req *stockpb.GetItemsReq) (*stockpb.GetItemsRes, error) {

	_, span := tracing.Start(ctx, "GetItems")
	defer span.End()

	items, err := G.app.Queries.GetItems.Handle(ctx, query.GetItems{ItemIDs: req.ItemIDs})
	if err != nil {
		return nil, err
	}
	return &stockpb.GetItemsRes{Items: convertor.NewItemConvertor().EntitiesToProtos(items)}, nil
}

func (G GRPCServer) CheckIfItemInStock(ctx context.Context, req *stockpb.CheckIfItemInStockReq) (*stockpb.CheckIfItemInStockRes, error) {

	_, span := tracing.Start(ctx, "CheckIfItemInStock")
	defer span.End()

	items, err := G.app.Queries.CheckIfItemsInStock.Handle(ctx, query.CheckIfItemsInStock{
		Items: convertor.NewItemWithQuantityConvertor().ProtosToEntities(req.Items)})
	if err != nil {
		return nil, err
	}

	return &stockpb.CheckIfItemInStockRes{
		InStock: 1,
		Items:   convertor.NewItemConvertor().EntitiesToProtos(items),
	}, nil
}

package convertor

import (
	"github.com/Hypocrite/gorder/common/genproto/orderpb"
	"github.com/Hypocrite/gorder/stock/entity"
)

type OrderConvertor struct{}

type ItemConvertor struct{}

type ItemWithQuantityConvertor struct{}

func (c *OrderConvertor) EntityToProto(eo *entity.Order) *orderpb.Order {
	return &orderpb.Order{
		ID:          eo.ID,
		CustomerID:  eo.CustomerID,
		Status:      eo.Status,
		Items:       NewItemConvertor().EntitiesToProtos(eo.Items),
		PaymentLink: eo.PaymentLink,
	}
}

func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *entity.Order {
	return &entity.Order{
		ID:          o.ID,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().ProtosToEntities(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *ItemConvertor) EntitiesToProtos(items []*entity.Item) (res []*orderpb.Item) {
	for _, item := range items {
		res = append(res, c.EntityToProto(item))
	}
	return
}

func (c *ItemConvertor) ProtosToEntities(items []*orderpb.Item) (res []*entity.Item) {
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return
}

func (c *ItemConvertor) EntityToProto(item *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemConvertor) ProtoToEntity(item *orderpb.Item) *entity.Item {
	return &entity.Item{
		ID:       item.ID,
		Name:     item.Name,
		Quantity: item.Quantity,
		PriceID:  item.PriceID,
	}
}

func (c *ItemWithQuantityConvertor) EntitiesToProtos(items []*entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.EntityToProto(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) EntityToProto(i *entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
	return &orderpb.ItemWithQuantity{
		ID:       i.ID,
		Quantity: i.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ProtosToEntities(items []*orderpb.ItemWithQuantity) (res []*entity.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.ProtoToEntity(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) ProtoToEntity(i *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       i.ID,
		Quantity: i.Quantity,
	}
}

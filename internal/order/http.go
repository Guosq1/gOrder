package main

import (
	"fmt"

	"github.com/Hypocrite/gorder/common"
	client "github.com/Hypocrite/gorder/common/client/order"
	"github.com/Hypocrite/gorder/order/app"
	"github.com/Hypocrite/gorder/order/app/command"
	"github.com/Hypocrite/gorder/order/app/dto"
	"github.com/Hypocrite/gorder/order/app/query"
	"github.com/Hypocrite/gorder/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID string) {

	var (
		req  client.CreateOrderReq
		err  error
		resp dto.CreateOrderResp
	)

	defer func() {
		H.Response(c, err, &resp)
	}()

	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}

	if err = H.validate(req); err != nil {
		//err = errors.NewWithError(consts.ErrnoRequestValidateError, err)
		return
	}

	r, err := H.app.Commands.CreateOrder.Handle(c.Request.Context(), command.CrateOrder{
		CustomerID: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientsToEntities(req.Items),
	})
	if err != nil {
		return
	}
	resp = dto.CreateOrderResp{
		OrderID:     r.OrderID,
		CustomerID:  req.CustomerId,
		RedirectURL: fmt.Sprintf("http://localhost:7777/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderID),
	}
}

func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerID string, orderID string) {

	var (
		err  error
		resp struct {
			Order *client.Order
		}
	)

	defer func() {
		H.Response(c, err, resp)
	}()

	o, err := H.app.Queries.GetCustomerOrder.Handle(c.Request.Context(), query.GetCustomerOrder{
		OrderID:    orderID,
		CustomerID: customerID,
	})

	if err != nil {
		return
	}

	resp.Order = convertor.NewOrderConvertor().EntityToClient(o)
}

func (H HTTPServer) validate(req client.CreateOrderReq) error {
	for _, v := range req.Items {
		if v.Quantity <= 0 {
			return fmt.Errorf("quantity must be positive, got %d from %s", v.Quantity, v.Id)
		}
	}
	return nil
}

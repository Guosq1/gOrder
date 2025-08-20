package service

import (
	"context"

	"github.com/Hypocrite/gorder/order/app"
)

func NewApplication(ctx context.Context) app.Application {
	//orderRepo := adapters.NewMemoryOrderRepository()

	return app.Application{
		//Queries: New(orderRepo),
	}
}

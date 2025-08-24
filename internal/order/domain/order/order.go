package order

import (
	"fmt"

	"github.com/Hypocrite/gorder/order/entity"
	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v82"
)

type Order struct {
	ID          string
	CustomerID  string
	Status      string
	PaymentLink string
	Items       []*entity.Item
}

func NewOrder(id, customerID, status, paymentLink string, items []*entity.Item) (*Order, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}
	if customerID == "" {
		return nil, errors.New("customerID is empty")
	}
	if status == "" {
		return nil, errors.New("status is empty")
	}
	if items == nil {
		return nil, errors.New("items is empty")
	}
	return &Order{
		ID:          id,
		CustomerID:  customerID,
		Status:      status,
		PaymentLink: paymentLink,
		Items:       items,
	}, nil
}

func (o *Order) IsPaid() error {
	if o.Status == string(stripe.CheckoutSessionPaymentStatusPaid) {
		return nil
	}
	return fmt.Errorf("order %s is not paid, status = %s", o.ID, o.Status)
}

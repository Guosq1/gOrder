package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Hypocrite/gorder/common/genproto/orderpb"
	"github.com/Hypocrite/gorder/common/tracing"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

const (
	successURL = "http://localhost:7777/success"
)

type StripeProcessor struct {
	apiKey string
}

func NewStripeProcessor(apiKey string) *StripeProcessor {
	if apiKey == "" {
		panic("api_key is empty")
	}
	stripe.Key = apiKey
	return &StripeProcessor{apiKey: apiKey}
}

func (s StripeProcessor) CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error) {
	_, span := tracing.Start(ctx, "stripe_processor.create_payment_link")
	defer span.End()

	var items []*stripe.CheckoutSessionLineItemParams
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			//price_1RytaBG0qlOoZ1d5HgarxICR
			Price:    stripe.String("price_1RyY3S8dBxTVAlmEOiBdniHt"),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}
	marshalledItems, _ := json.Marshal(order.Items)
	metadata := map[string]string{
		"orderID":     order.ID,
		"customerID":  order.CustomerID,
		"status":      order.Status,
		"items":       string(marshalledItems),
		"paymentLink": order.PaymentLink,
	}

	params := &stripe.CheckoutSessionParams{
		Metadata:   metadata,
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(fmt.Sprintf("%s?customerID=%s&orderID=%s", successURL, order.CustomerID, order.ID)),
	}
	result, err := session.New(params)
	if err != nil {
		return "", err
	}
	return result.URL, nil
}

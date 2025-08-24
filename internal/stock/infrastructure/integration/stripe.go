package integration

import (
	"context"

	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/product"
)

type StripeApi struct {
	apiKey string
}

func NewStripeApi() *StripeApi {
	return &StripeApi{apiKey: viper.GetString("stripe-key")}
}

func (s *StripeApi) GetPriceByProductID(ctx context.Context, pid string) (string, error) {
	stripe.Key = s.apiKey
	params := &stripe.ProductParams{}
	result, err := product.Get(pid, params)

	if err != nil {
		return "", err
	}
	return result.DefaultPrice.ID, err
}

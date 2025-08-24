package tests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	sw "github.com/Hypocrite/gorder/common/client/order"
	_ "github.com/Hypocrite/gorder/common/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	ctx    = context.Background()
	server = fmt.Sprintf("http://%s/api", viper.GetString("order.http-addr"))
)

func TestMain(m *testing.M) {
	before()
	m.Run()
}

func before() {
	log.Printf("server %s start", server)
}

func TestCreateOrder_success(t *testing.T) {

	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: "123",
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_SvUtpxV5be9KmM",
				Quantity: 1,
			},
			{
				Id:       "prod_SvUYSv739cprIh",
				Quantity: 1,
			},
		},
	})
	t.Logf("body=%s", string(response.Body))
	assert.Equal(t, http.StatusOK, response.StatusCode())
	assert.Equal(t, 0, response.JSON200.Errcode)
}

func TestCreateOrder_invalidParams(t *testing.T) {

	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: "123",
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_SvUtpxV5be9KmM",
				Quantity: 1,
			},
			{
				Id:       "prod_SvUYSv739cprIh",
				Quantity: 1,
			},
			{
				Id:       "prod_SvUYgRlyMkyFB1",
				Quantity: 1,
			},
		},
	})
	assert.Equal(t, 200, response.StatusCode())
	assert.Equal(t, 0, response.JSON200.Errcode)
}

func getResponse(t *testing.T, customerID string, body sw.PostCustomerCustomerIdOrdersJSONRequestBody) *sw.PostCustomerCustomerIdOrdersResponse {
	t.Helper()
	client, err := sw.NewClientWithResponses(server)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("getResponse body=%+v", body)
	response, err := client.PostCustomerCustomerIdOrdersWithResponse(ctx, customerID, body)
	if err != nil {
		t.Fatal(err)
	}
	return response
}

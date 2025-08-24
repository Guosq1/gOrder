package dto

type CreateOrderResp struct {
	OrderID     string `json:"order_id"`
	CustomerID  string `json:"customer_id"`
	RedirectURL string `json:"redirect_url"`
}

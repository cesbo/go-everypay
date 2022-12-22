package everypay

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Data to initiate payment
type OneOff struct {
	// Transaction amount.
	// The currency is taken from the specified processing account.
	// Can be also set as 0 for card verification (saving card for further token payments).
	Amount float32 `json:"amount"`

	// Unique order reference.
	OrderReference string `json:"order_reference"`

	// When description is provided, it will be used with open banking payment methods.
	Description string `json:"payment_description,omitempty"`

	// Customer’s email for fraud prevention.
	CustomerEmail string `json:"email"`

	// Customer’s IP address for fraud prevention.
	CustomerIp string `json:"customer_ip"`

	// URL where the Customer should be redirected after completing the payment.
	CustomerUrl string `json:"customer_url"`

	// It must be sent when RequestToken defined.
	// It is the type of the agreement:
	// - `unscheduled`
	// - `recurring`
	TokenAgreement string `json:"token_agreement,omitempty"`

	// Token should be returned in the response for future usage
	RequestToken string `json:"request_token,omitempty"`
}

type oneOffRequest struct {
	*OneOff

	ApiUsername string    `json:"api_username"`
	AccountName string    `json:"account_name"`
	Nonce       string    `json:"nonce"`
	Timestamp   time.Time `json:"timestamp"`
}

type oneOffResponse struct {
	OrderReference   string `json:"order_reference"`
	PaymentReference string `json:"payment_reference"`
	PaymentLink      string `json:"payment_link"`
}

// InitialPayment returns the payment link
func (e *Everypay) InitialPayment(o *OneOff) (string, error) {
	requestData := &oneOffRequest{
		OneOff:      o,
		ApiUsername: e.username,
		AccountName: e.account,
		Nonce:       uuid.NewString(),
		Timestamp:   time.Now().UTC(),
	}

	responseData := &oneOffResponse{}

	if err := e.request("payments/oneoff", requestData, responseData); err != nil {
		return "", fmt.Errorf("oneoff: %w", err)
	}

	return responseData.PaymentLink, nil
}

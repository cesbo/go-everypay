package everypay

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Payment struct {
	// Payment standing amount.
	Amount float32 `json:"standing_amount"`

	// Unique order reference.
	OrderReference string `json:"order_reference"`

	// Current status of the payment. Possible values:
	// - `initial`
	// - `waiting_for_sca`
	// - `sent_for_processing`
	// - `waiting_for_3ds_response`
	// - `settled`
	// - `failed`
	// - `abandoned`
	// - `voided`
	// - `refunded`
	// - `chargebacked`
	PaymentState string `json:"payment_state"`

	// Time of the transaction
	TransactionTime time.Time `json:"transaction_time"`
}

// GetPayment returns the status of a payment with the given payment reference.
func (e *Everypay) GetPayment(reference string) (*Payment, error) {
	u := &url.URL{
		Path:     "payments/" + reference,
		RawQuery: "api_username=" + e.username,
	}

	responseData := &Payment{}

	if err := e.request(http.MethodGet, u, nil, responseData); err != nil {
		return nil, fmt.Errorf("payment status: %w", err)
	}

	return responseData, nil
}

package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/literate-octo-pancake/currency/protos/currency"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	log hclog.Logger
}

// NewCurrency creates a new Currency server
func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{l}
}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, req *currency.RateRequest) (*currency.RateResponse, error){
	c.log.Info("Handle GetRate", "base", req.GetBase(), "destination", req.GetDestination())

	return &currency.RateResponse{Rate: 0.5}, nil
}
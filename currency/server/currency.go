package server

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/literate-octo-pancake/currency/data"
	"github.com/jalexanderII/literate-octo-pancake/currency/protos/currency"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	log   hclog.Logger
	rates *data.ExchangeRates
}

// NewCurrency creates a new Currency server
func NewCurrency(l hclog.Logger, r *data.ExchangeRates) *Currency {
	return &Currency{l, r}
}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies.
func (c *Currency) GetRate(ctx context.Context, req *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", req.GetBase(), "destination", req.GetDestination())

	rate, err := c.rates.GetRates(req.GetBase().String(), req.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &currency.RateResponse{Rate: rate}, nil
}

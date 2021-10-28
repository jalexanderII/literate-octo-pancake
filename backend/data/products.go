package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/literate-octo-pancake/currency/protos/currency"
)

var ErrProductNotFound = fmt.Errorf("product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this product
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this product
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[0-9]+
	SKU string `json:"sku" validate:"required,sku"`
}

// Products is a collection of Product
type Products []*Product

func fxPrice(rate float32, p Products) Products {
	npl := make(Products, len(p))
	for idx, product := range p {
		np := *product
		np.Price *= rate
		npl[idx] = &np
	}
	return npl
}

type ProductsDB struct {
	log      hclog.Logger
	currency currency.CurrencyClient
}

func NewProductsDB(l hclog.Logger, c currency.CurrencyClient) *ProductsDB {
	return &ProductsDB{l, c}
}

func (pdb *ProductsDB) getRate(dest string) (float32, error) {
	// get exchange rate
	rr := &currency.RateRequest{
		Base:        currency.RateRequest_EUR,
		Destination: currency.RateRequest_Currencies(currency.RateRequest_Currencies_value[dest]),
	}

	resp, err := pdb.currency.GetRate(context.Background(), rr)
	return resp.Rate, err
}

// GetProducts returns a list of products
func (pdb *ProductsDB) GetProducts(dest string) (Products, error) {
	if dest == "" {
		return productList, nil
	}

	// get exchange rate
	rate, err := pdb.getRate(dest)
	if err != nil {
		pdb.log.Error("Error doing currency conversion", "destination", dest)
	}

	return fxPrice(rate, productList), nil
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func (pdb *ProductsDB) GetProductByID(id int, dest string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	if dest == "" {
		return productList[i], nil
	}

	// get exchange rate
	rate, err := pdb.getRate(dest)
	if err != nil {
		pdb.log.Error("Error doing currency conversion", "destination", dest)
	}

	// new productlist with only one product
	pl := []*Product{productList[i]}

	return fxPrice(rate, pl)[0], nil
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func (pdb *ProductsDB) UpdateProduct(p Product, dest string) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// get exchange rate
	rate, err := pdb.getRate(dest)
	if err != nil {
		pdb.log.Error("Error doing currency conversion", "destination", dest)
	}

	p.Price *= rate

	// update the product in the DB
	productList[i] = &p

	return nil
}

// AddProduct adds a new product to the database
func AddProduct(p *Product) {
	// get the next id in sequence
	maxID := productList[len(productList)-1].ID
	p.ID = maxID + 1
	productList = append(productList, p)
}

// DeleteProduct deletes a product from the database
func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       4.25,
		SKU:         "abc-123",
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "short and strong coffee without milk",
		Price:       2.00,
		SKU:         "fjk-123",
	},
}

package handlers

import (
	"context"
	"github.com/jalexanderII/literate-octo-pancake/currency/protos/currency"
	"net/http"

	"github.com/jalexanderII/literate-octo-pancake/backend/data"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(w http.ResponseWriter, _ *http.Request) {
	p.l.Println("[DEBUG] get all records")
	w.Header().Add("Content-Type", "application/json")

	prods := data.GetProducts()

	err := data.ToJSON(prods, w)
	if err != nil {
		// we should never be here but log the error just in case
		p.l.Println("[ERROR] serializing product", err)
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	id := getProductID(r)

	p.l.Println("[DEBUG] get record id", id)

	prod, err := data.GetProductByID(id)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)

		w.WriteHeader(http.StatusNotFound)
		err := data.ToJSON(&GenericError{Message: err.Error()}, w)
		if err != nil {
			return
		}
		return
	default:
		p.l.Println("[ERROR] fetching product", err)

		w.WriteHeader(http.StatusInternalServerError)
		err := data.ToJSON(&GenericError{Message: err.Error()}, w)
		if err != nil {
			return
		}
		return
	}

	// get exchange rate
	rr := &currency.RateRequest{
		Base:        currency.RateRequest_EUR,
		Destination: currency.RateRequest_GBP,
	}

	resp, err := p.c.GetRate(context.Background(), rr)
	if err != nil {
		p.l.Println("[Error] error getting new rate", err)
		err := data.ToJSON(&GenericError{Message: err.Error()}, w)
		if err != nil {
			return
		}
		return
	}

	p.l.Printf("Resp %#v", resp)

	prod.Price = prod.Price * resp.Rate

	err = data.ToJSON(prod, w)
	if err != nil {
		// we should never be here but log the error just in case
		p.l.Println("[ERROR] serializing product", err)
	}
}

// swagger:route POST /products products createProduct
// Create a new product
//
// responses:
//	200: productResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new products
func (p *Products) Create(_ http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	p.l.Printf("[DEBUG] Inserting product: %#v\n", prod)
	data.AddProduct(prod)
}

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update products
func (p *Products) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Println("[DEBUG] updating record id", prod.ID)

	err := data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] product not found", err)

		w.WriteHeader(http.StatusNotFound)
		err := data.ToJSON(&GenericError{Message: "Product not found in database"}, w)
		if err != nil {
			return
		}
		return
	}

	// write the no content success header
	w.WriteHeader(http.StatusNoContent)
}

// swagger:route DELETE /products/{id} products deleteProduct
// Update a products details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  501: errorResponse

// Delete handles DELETE requests and removes items from the database
func (p *Products) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := getProductID(r)

	p.l.Println("[DEBUG] deleting record id", id)

	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Println("[ERROR] deleting record id does not exist")

		w.WriteHeader(http.StatusNotFound)
		err := data.ToJSON(&GenericError{Message: err.Error()}, w)
		if err != nil {
			return
		}
		return
	}

	if err != nil {
		p.l.Println("[ERROR] deleting record", err)

		w.WriteHeader(http.StatusInternalServerError)
		err := data.ToJSON(&GenericError{Message: err.Error()}, w)
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

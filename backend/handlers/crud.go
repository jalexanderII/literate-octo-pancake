package handlers

import (
	"net/http"

	"github.com/jalexanderII/literate-octo-pancake/backend/data"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(w http.ResponseWriter, r *http.Request) {
	p.l.Debug("Get all records")
	w.Header().Add("Content-Type", "application/json")

	prods, err := p.pdb.GetProducts(r.URL.Query().Get("currency"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, w)
		return
	}

	err = data.ToJSON(prods, w)
	if err != nil {
		// we should never be here but log the error just in case
		p.l.Error("Unable to serialize product", "error", err)
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

	cur := r.URL.Query().Get("currency")
	p.l.Debug("Get record id", "id", id, "currency", cur)

	prod, err := p.pdb.GetProductByID(id, cur)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("Unable to fetch product", "error", err)

		w.WriteHeader(http.StatusNotFound)
		err := data.ToJSON(&GenericError{Message: err.Error()}, w)
		if err != nil {
			return
		}
		return
	default:
		p.l.Error("Unable to fetch product", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		err := data.ToJSON(&GenericError{Message: err.Error()}, w)
		if err != nil {
			return
		}
		return
	}

	err = data.ToJSON(prod, w)
	if err != nil {
		// we should never be here but log the error just in case
		p.l.Error("Unable to serialize product", "error", err)
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

	p.l.Debug("Inserting product %#v\n", prod)
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
	p.l.Debug("Updating record id", "id", prod.ID)

	err := p.pdb.UpdateProduct(prod, "")
	if err == data.ErrProductNotFound {
		p.l.Error("Product not found", "error", err)

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

	p.l.Debug("Deleting record id", "id", id)

	err := data.DeleteProduct(id)
	if err == data.ErrProductNotFound {
		p.l.Error("Unable to delete record, id does not exist", "error", err)

		w.WriteHeader(http.StatusNotFound)
		err := data.ToJSON(&GenericError{Message: err.Error()}, w)
		if err != nil {
			return
		}
		return
	}

	if err != nil {
		p.l.Error("Unable to delete record", "error", err)

		w.WriteHeader(http.StatusInternalServerError)
		err := data.ToJSON(&GenericError{Message: err.Error()}, w)
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

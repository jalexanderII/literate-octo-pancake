package handlers

import (
	"context"
	"net/http"

	"github.com/jalexanderII/literate-octo-pancake/backend/data"
)

// MiddlewareValidateProduct validates the product in the request and calls next if ok
func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		prod := &data.Product{}

		err := data.FromJSON(prod, r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)

			w.WriteHeader(http.StatusBadRequest)
			err := data.ToJSON(&GenericError{Message: err.Error()}, w)
			if err != nil {
				return
			}
			return
		}

		// validate the product
		errs := p.v.Validate(prod)
		if len(errs) != 0 {
			p.l.Println("[ERROR] validating product", errs)

			// return the validation messages as an array
			w.WriteHeader(http.StatusUnprocessableEntity)
			err := data.ToJSON(&ValidationError{Messages: errs.Errors()}, w)
			if err != nil {
				return
			}
			return
		}

		// add the product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

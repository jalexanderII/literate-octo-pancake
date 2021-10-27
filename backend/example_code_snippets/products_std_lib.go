package example_code_snippets
//
//import (
//	"fmt"
//	"log"
//	"net/http"
//	"regexp"
//	"strconv"
//
//	"github.com/Fudoshin2596/curly-telegram/backend/data"
//)
//
//// Products is a simple handler
//type Products struct {
//	l *log.Logger
//}
//
//// NewProducts creates a new Products handler with the given logger
//func NewProducts(l *log.Logger) *Products {
//	return &Products{l}
//}
//
//// ServeHTTP is the main entry point for the handler and implements the go http.Handler interface
//// https://golang.org/pkg/net/http/#Handler
//func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	switch method := r.Method; method {
//	case http.MethodGet:
//		// handle the request for a list of products
//		p.getProducts(w, r)
//		return
//	case http.MethodPost:
//		p.addProduct(w, r)
//		return
//	case http.MethodPut:
//
//		p.updateProduct(w, r)
//		return
//	default:
//		// catch all
//		// if no method is satisfied return an error
//		w.WriteHeader(http.StatusMethodNotAllowed)
//	}
//}
//
//func getIDFromURI(r *http.Request) (int, error) {
//	reg := regexp.MustCompile("/([0-9]+)")
//	g := reg.FindAllStringSubmatch(r.URL.Path, -1)
//	if len(g) != 1 || len(g[0]) != 2 {
//		return -1, fmt.Errorf("invalid URI")
//	}
//	return strconv.Atoi(g[0][1])
//}
//
//// getProducts returns the products from the data store
//func (p *Products) getProducts(w http.ResponseWriter, _ *http.Request) {
//	p.l.Println("Handle GET Products")
//
//	// fetch the products from the datastore
//	lp := data.GetProducts()
//
//	// serialize the list to JSON and write to client
//	err := lp.ToJSON(w)
//	if err != nil {
//		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
//	}
//}
//
//// addProducts accepts a json product fields from user and unmarshalls it into a Product and saves it
//func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
//	p.l.Println("Handle POST Products")
//
//	prod := &data.Product{}
//	err := prod.FromJSON(r.Body)
//	if err != nil {
//		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
//	}
//
//	p.l.Printf("Prod: %#v added to datastore", prod)
//	data.AddProduct(prod)
//}
//
//func (p *Products) updateProduct(w http.ResponseWriter, r *http.Request) {
//	p.l.Println("Handle PUT Product")
//
//	prod := &data.Product{}
//	err := prod.FromJSON(r.Body)
//	if err != nil {
//		http.Error(w, "Unable to unmarshal json", http.StatusBadRequest)
//	}
//
//	id, err := getIDFromURI(r)
//	if err != nil {
//		http.Error(w, "Invalid URI", http.StatusBadRequest)
//	}
//
//	err = data.UpdateProduct(id, prod)
//	if err == data.ErrProductNotFound {
//		http.Error(w, "Product not found", http.StatusNotFound)
//		return
//	}
//
//	if err != nil {
//		http.Error(w, "Product not found", http.StatusInternalServerError)
//		return
//	}
//
//}

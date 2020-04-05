// Package classification of Product API
//
// Documentation for Product API
//
// 		Schemes: http
// 		BasePath: /
// 		Version: 1.0.0
//
// 		Consumes:
// 		- application/json
//
// 		Produces:
// 		- application/json
//
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/kerembalci90/go-microservice-demo/data"
)

// A list of products returned in the response
// swagger:response productsResponse
type productsResponseSwaggerWrapper struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:response noContent
type productNoResponseSwaggerWrapper struct {
}

// swagger:parameters deleteProduct
type productIDParameterSwaggerWrapper struct {
	// ID of the product to be deleted from the database
	// in: path
	// required: true
	ID int `json:"id"`
}

// Products Struct that provides access to operating on Products
type Products struct {
	log *log.Logger
}

// NewProducts Factory function for returning a new instance of Products handler
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
// 	200: productsResponse

// GetProducts Responsible for handling the return of a list of products
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle GET Products request")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	// d, err := json.Marshal(lp)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// CreateProduct Responsible for handling a product creation
func (p *Products) CreateProduct(rw http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle POST Product request")

	product := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&product)
}

// UpdateProduct Responsible for handling the update of an existing product
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle PUT Product request")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

// swagger:route DELETE /products/{id} products deleteProducts
// Deletes a product from the list of products
// responses:
// 	204: noContent

func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle DELETE Product request")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	err = data.DeleteProduct(id)
	if err != nil {
		http.Error(rw, "Unable to delete product", http.StatusInternalServerError)
		return
	}
}

// KeyProduct context key
type KeyProduct struct{}

// ProductValidationMiddleware Middleware handler for constructing model from json data
func (p Products) ProductValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		product := data.Product{}
		err := product.FromJSON(r.Body)
		if err != nil {
			p.log.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Unable to unmarshall JSON", http.StatusBadRequest)
			return
		}

		err = product.Validate()
		if err != nil {
			p.log.Println("[ERROR] validating product", err)
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, product)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

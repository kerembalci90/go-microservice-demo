package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator"
)

// Product denotes a product item for our coffee shop
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// FromJSON Converts a product from JSON format to data format
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

// Validate Validates the product based on the validation rules set for that product model
func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku is of format abc-asdf-wefw
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

// Products denotes a list of product items
type Products []*Product

// ToJSON Converts the list of products to JSON format for output
func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

// GetProducts Returns the list of products
func GetProducts() Products {
	return productList
}

// AddProduct Adds a new product to the product list
func AddProduct(p *Product) {
	p.ID = getNextProductID()
	productList = append(productList, p)
}

// UpdateProduct Takes care of updating a product with the given id
func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProductByID(id)
	if err != nil {
		return err
	}

	p.ID = id
	productList[pos] = p
	return nil
}

// ErrorProductNotFound Custom error for denoting that a product is not found
var ErrorProductNotFound = fmt.Errorf("Product not found")

// findProductById retuns the product and the position of it in the product list
func findProductByID(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}

	return nil, -1, ErrorProductNotFound
}

func getNextProductID() int {
	p := productList[len(productList)-1]
	return p.ID + 1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

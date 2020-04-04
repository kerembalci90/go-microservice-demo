package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "john",
		Price: 1.00,
		SKU:   "abs-wefwe-wef",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}

package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

//product defines the structure for an api product
type Product struct {
	ID          int     `json:"id"`                       //remember no space between colon(:) and inverted comma allowed
	Name        string  `json:"name" validate:"required"` //validate: "required" is a way to write required variable, check docs for more info on this
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"gt=0"` //validate: gt=0 price failed during test
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"` //no data displayed in this case
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`

	// ID          int
	// Name        string
	// Description string
	// Price       float64
	// SKU         string
	// CreatedOn   string
	// UpdatedOn   string
	// DeletedOn   string
}

//io.Reader is used to read what's on the screen
func (p *Product) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r) //creating a new decoder, passing our data via r variable
	return e.Decode(p)      //returns error if any
}

//video 6 part
//using validation from github
//validation function implementing interface
func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSku) //custom validation of sku
	return validate.Struct(p)
}

//how to create our own validation
func validateSku(fl validator.FieldLevel) bool {

	//regexp
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

type Products []*Product //creating own type

//we use encode function to dont use memory to store and display data
//io writer is used to write something on screen or whatever the writer is
func (x *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w) //creating an encoder object passed with the io writer, in this case might be rw
	return e.Encode(x)      //returns error if any
}

// func GetProducts() []*Product {
// 	return productList
// }

//using our type for linking
func GetProducts() Products {
	return productList
}

//adding the user data received
func Addproduct(p *Product) {
	p.ID = getNextID()                   //replaces the id with getnextid's return
	productList = append(productList, p) //appending the value
}

func UpdateProduct(id int, p *Product) error {
	_, pos, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	productList[pos] = p

	return nil
}

func findProduct(id int) (*Product, int, error) {
	for i, p := range productList {
		if p.ID == id {
			return p, i, nil
		}
	}
	return nil, -1, ErrorProductNotFound
}

//The Errorf function lets us use formatting features to create descriptive error messages
// const name, id = "bueller", 17
// 	err := fmt.Errorf("user %q (id %d) not found", name, id)
// 	fmt.Println(err.Error())
//output: user "bueller" (id 17) not found
var ErrorProductNotFound = fmt.Errorf("product not found")

func getNextID() int {
	lp := productList[len(productList)-1] //getting last entry
	lp.ID = lp.ID + 1                     //incrementing the last value
	return lp.ID                          //returning it
}

//we made a few changes
//code looks different from instructor's code
var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "frothy milk shake",
		Price:       2.45,
		SKU:         "abc234",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fgh234",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

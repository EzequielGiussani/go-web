package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type ResponseGetAllProducts struct {
	Message string    `json:"message"`
	Data    []Product `json:"data"`
	Error   bool      `json:"error"`
}

type ProductsMap struct {
	st map[int]Product
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Quantity    uint      `json:"quantity"`
	CodeValue   string    `json:"code_value"`
	IsPublished bool      `json:"is_published"`
	Expiration  time.Time `json:"expiration"`
	Price       float64   `json:"price"`
}

func NewProductsMAP() *ProductsMap {
	return &ProductsMap{
		st: make(map[int]Product),
	}
}

// Custom unmarshaler for the 'Expiration' field
func (p *Product) UnmarshalJSON(data []byte) error {
	type Alias Product
	aux := &struct {
		Expiration string `json:"expiration"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Parse the 'Expiration' string to time.Time
	expirationTime, err := time.Parse("02/01/2006", aux.Expiration)
	if err != nil {
		return err
	}
	p.Expiration = expirationTime

	return nil
}

func (ps *ProductsMap) LoadProducts() (*ProductsMap, error) {
	file, err := os.Open("./products.json")

	if err != nil {
		return nil, errors.New(err.Error())
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	products := make([]Product, 0)

	if err := decoder.Decode(&products); err != nil {
		return nil, errors.New(err.Error())
	}

	for _, prod := range products {
		ps.st[prod.ID] = prod
	}

	return ps, nil

}

func (ps *ProductsMap) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		products := ps.st

		productsList := make([]Product, 0, len(products))

		for _, prod := range products {
			productsList = append(productsList, prod)
		}

		body := &ResponseGetAllProducts{
			Message: "success",
			Data:    productsList,
			Error:   false,
		}
		json.NewEncoder(w).Encode(body)
	}
}

func (ps *ProductsMap) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")

		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "application/json")
			body := &ResponseGetAllProducts{
				Message: "error converting id to int",
				Data:    nil,
				Error:   true,
			}
			json.NewEncoder(w).Encode(body)
			return
		}

		products := ps.st

		productsList := make([]Product, 0, len(products))

		product, ok := products[id]

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Add("Content-Type", "application/json")
			body := &ResponseGetAllProducts{
				Message: "error product not found",
				Data:    nil,
				Error:   true,
			}
			json.NewEncoder(w).Encode(body)
			return
		}

		productsList = append(productsList, product)

		body := &ResponseGetAllProducts{
			Message: "success",
			Data:    productsList,
			Error:   false,
		}
		json.NewEncoder(w).Encode(body)
	}
}

func (ps *ProductsMap) GetBySearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		products := ps.st

		priceGt, err := strconv.ParseFloat(r.URL.Query().Get("priceGt"), 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "application/json")
			body := &ResponseGetAllProducts{
				Message: "error converting priceGt to float64",
				Data:    nil,
				Error:   true,
			}
			json.NewEncoder(w).Encode(body)
			return
		}

		productsList := make([]Product, 0, len(products))

		for _, prod := range products {
			if prod.Price > priceGt {
				productsList = append(productsList, prod)
			}
		}

		body := &ResponseGetAllProducts{
			Message: "success",
			Data:    productsList,
			Error:   false,
		}
		json.NewEncoder(w).Encode(body)
	}
}

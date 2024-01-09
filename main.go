package main

import (
	"fmt"
	"net/http"

	"github.com/EzequielGiussani/go-web/internal/product/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("pong"))

		if err != nil {
			fmt.Println(err.Error())
			return
		}
	})

	handler := handlers.NewProductsMAP()

	handler, err := handler.LoadProducts()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	router.Route("/products", func(r chi.Router) {
		r.Get("/", handler.GetAll())
		r.Get("/{id}", handler.GetById())
		r.Get("/search", handler.GetBySearch())
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}

}

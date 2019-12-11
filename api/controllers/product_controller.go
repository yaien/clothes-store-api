package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/yaien/clothes-store-api/api/helpers/response"

	"github.com/gorilla/mux"
	"github.com/yaien/clothes-store-api/api/models"
	"github.com/yaien/clothes-store-api/api/services"
)

type ProductController struct {
	Products services.ProductService
}

func (p *ProductController) Create(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	err = p.Products.Create(&product)
	if err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	response.Send(w, product)

}
func (p *ProductController) Find(w http.ResponseWriter, r *http.Request) {
	products, err := p.Products.Find()
	if err != nil {
		response.Error(w, err, http.StatusInternalServerError)
		return
	}
	response.Send(w, products)
}

func (p *ProductController) Param(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id := mux.Vars(r)["product_id"]
	product, err := p.Products.Get(id)
	if err != nil {
		response.Error(w, errors.New("PRODUCT_NOT_FOUND"), http.StatusNotFound)
		return
	}
	ctx := context.WithValue(r.Context(), key("product"), product)
	next(w, r.WithContext(ctx))
}

func (p *ProductController) Show(w http.ResponseWriter, r *http.Request) {
	product := r.Context().Value(key("product")).(*models.Product)
	response.Send(w, product)
}

func (p *ProductController) Update(w http.ResponseWriter, r *http.Request) {
	var data models.Product
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}

	product := r.Context().Value(key("product")).(*models.Product)
	data.ID = product.ID
	if err := p.Products.Update(&data); err != nil {
		response.Error(w, err, http.StatusBadRequest)
		return
	}
	response.Send(w, data)
}

package handlers

import (
	"log"
	"net/http"

	"github.com/dalbqrq/nic-jackson/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle get
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	// handle put
	if r.Method == http.MethodPut {
		p.putProducts(rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) putProducts(rw http.ResponseWriter, r *http.Request) {
	http.Error(rw, "no put implemented", http.StatusInternalServerError)
}

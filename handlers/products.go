package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/dalbqrq/nic-jackson/data"
)

type Products struct {
	log *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// handle get
	if r.Method == http.MethodGet {
		p.log.Println("GET", r.URL.Path)
		p.getProducts(rw, r)
		return
	}

	// handle post
	if r.Method == http.MethodPost {
		p.log.Println("POST", r.URL.Path)
		p.postProducts(rw, r)
		return
	}

	// handle put
	if r.Method == http.MethodPut {
		p.log.Println("PUT", r.URL.Path)

		// expect id in the URI
		re := regexp.MustCompile(`/([0-9]+)`)
		g := re.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.log.Println("Invalid URI more than one id", g)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.log.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.log.Println("Got ID:", id)
		p.putProducts(id, rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle GET Products")

	prod := data.GetProducts()
	err := prod.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) postProducts(rw http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle POST Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusInternalServerError)
	}

	p.log.Printf("Prod: %#v", prod)

	data.AddProduct(prod)
}

func (p *Products) putProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.log.Println("Handle PUT Products")

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusInternalServerError)
	}

	p.log.Printf("Prod: %#v", prod)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

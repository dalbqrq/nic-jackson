package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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
		p.l.Println("GET", r.URL.Path)
		p.getProducts(rw, r)
		return
	}

	// handle post
	if r.Method == http.MethodPost {
		p.l.Println("POST", r.URL.Path)
		p.postProducts(rw, r)
		return
	}

	// handle put
	if r.Method == http.MethodPut {
		p.l.Println("PUT", r.URL.Path)

		// expect id in the URI
		re := regexp.MustCompile(`/([0-9]+)`)
		g := re.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id", g)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.l.Println("Got ID:", id)
		p.putProducts(id, rw, r)
		return
	}

	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")

	// fetch the products from the datastore
	lp := data.GetProducts()

	// serialize the list to JSON
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (pl *Products) postProducts(rw http.ResponseWriter, r *http.Request) {
	pl.l.Println("Handle POST Products")

	p := &data.Product{}
	err := p.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusInternalServerError)
	}

	pl.l.Printf("Prod: %#v", p)

	data.AddProduct(p)
}

func (pl *Products) putProducts(id int, rw http.ResponseWriter, r *http.Request) {
	pl.l.Println("Handle PUT Products")

	p := &data.Product{}
	err := p.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusInternalServerError)
	}

	pl.l.Printf("Prod: %#v", p)

	err = data.UpdateProduct(id, p)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

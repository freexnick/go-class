package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type dollars float32
type database struct {
	mu sync.Mutex
	db map[string]dollars
}

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

func (d *database) list(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for item, price := range d.db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (d *database) add(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if _, ok := d.db[item]; ok {
		msg := fmt.Sprintf("duplicate item: %q", item)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		msg := fmt.Sprintf("invalid price: %q", price)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	d.db[item] = dollars(p)

	fmt.Fprintf(w, "added %s with price %s\n", item, d.db[item])
}

func (d *database) update(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if _, ok := d.db[item]; !ok {
		msg := fmt.Sprintf("no such item: %q", item)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		msg := fmt.Sprintf("invalid price: %q", price)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	d.db[item] = dollars(p)

	fmt.Fprintf(w, "updated %s with price %s\n", d.db[item], item)
}

func (d *database) find(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	item := req.URL.Query().Get("item")

	if _, ok := d.db[item]; !ok {
		msg := fmt.Sprintf("no such item: %q", item)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "item %s has price %s\n", item, d.db[item])
}

func (d *database) delete(w http.ResponseWriter, req *http.Request) {
	d.mu.Lock()
	defer d.mu.Unlock()

	item := req.URL.Query().Get("item")

	if _, ok := d.db[item]; !ok {
		msg := fmt.Sprintf("no such item: %q", item)
		http.Error(w, msg, http.StatusNotFound)
		return
	}

	delete(d.db, item)

	fmt.Fprintf(w, "item %s has been deleted", item)
}

func runServer() {
	db := database{
		db: map[string]dollars{
			"shoes": 50,
			"socks": 5,
		},
	}

	http.HandleFunc("/list", db.list)
	http.HandleFunc("/create", db.add)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/find", db.find)
	http.HandleFunc("/delete", db.delete)

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	runServer()
}

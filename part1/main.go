package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// NOTE: don't do this in real life
type dollars float32

func (d dollars) String() string {
	return fmt.Sprintf("$%.2f", d)
}

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%-6s: %3s\n", item, price)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if _, ok := db[item]; ok {
		msg := fmt.Sprintf("duplicate item: %q", item)
		http.Error(w, msg, http.StatusBadRequest) // 400 - BAD REQUEST
		return
	}

	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		msg := fmt.Sprintf("invalid price: %q", price)
		http.Error(w, msg, http.StatusBadRequest) // 400 - BAD REQUEST
		return
	}

	db[item] = dollars(p)
	fmt.Fprintf(w, "item %s added with price %s\n", item, db[item])
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if _, ok := db[item]; !ok {
		msg := fmt.Sprintf("Invalid item: %q, It doesn't exist.", item)
		http.Error(w, msg, http.StatusNotFound) // 404 - ITEM NOT FOUND
		return
	}

	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		msg := fmt.Sprintf("invalid price: %q", price)
		http.Error(w, msg, http.StatusBadRequest) // 400 - BAD REQUEST
		return
	}

	db[item] = dollars(p)
	fmt.Fprintf(w, "item %s updated with price %s\n", item, db[item])
}

func (db database) read(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	if _, ok := db[item]; !ok {
		msg := fmt.Sprintf("Invalid item: %q, It doesn't exist.", item)
		http.Error(w, msg, http.StatusNotFound) // 404 - ITEM NOT FOUND
		return
	}

	fmt.Fprintf(w, "Found item\n%-6s: %3s\n", item, db[item])
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	if _, ok := db[item]; !ok {
		msg := fmt.Sprintf("Invalid item: %q, It doesn't exist.", item)
		http.Error(w, msg, http.StatusNotFound) // 404 - ITEM NOT FOUND
		return
	}

	delete(db, item)

	fmt.Fprintf(w, "item %s deleted\n", item)
}

func main() {
	db := database{
		"shoes": 50,
		"socks": 5,
	}

	// NOTE that these are all method values
	// (closing over the object "db")

	http.HandleFunc("/list", db.list)
	http.HandleFunc("/create", db.create)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)
	http.HandleFunc("/read", db.read)

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

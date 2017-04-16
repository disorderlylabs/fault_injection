package main

import (
	cart "fault_injection/E-Commerce_App/shopping_cart/model"
	"fmt"
	"net/http"
	"strconv"
)

var (
	location = 1339
	redisport = 1337
)




func main() {
	cart.Init(redisport)

	fmt.Printf("Listening on port: %d\n", location)

	http.HandleFunc("/cart/add", add)
	http.HandleFunc("/cart/delete", delete)
	http.HandleFunc("/cart/items", getItems)
	http.HandleFunc("/cart/create", create)
	http.ListenAndServe("localhost:"+ strconv.Itoa(location), nil)
}


//creates a shopping cart and returns the cartID
func create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cart create")
	//only accept POST requests
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	//get the userID from request body
	userID := r.PostFormValue("userID")
	if userID == "" {
		http.Error(w, "could not parse userID", 400)
		return
	}

	//call internal function to create cartID
	cartID, err := cart.CreateCart(userID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//write cartID back in response body
	fmt.Fprintf(w, "%d", cartID)
}


//adds an item to a cart
func add(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	cartID := r.PostFormValue("cartID")
	if cartID == "" {
		http.Error(w, "could not parse userID", 400)
		return
	}

	itemID := r.PostFormValue("itemID")
	if itemID == "" {
		http.Error(w, "could not parse userID", 400)
		return
	}

	err := cart.AddItem(cartID, itemID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}



//deletes an item from a cart
func delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.Header().Set("Allow", "DELETE")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	cartID := r.URL.Query().Get("cartID")
	if cartID == "" {
		http.Error(w, "could not parse cartID", 400)
		return
	}

	itemID := r.URL.Query().Get("itemID")
	if itemID == "" {
		http.Error(w, "could not parse itemID", 400)
		return
	}

	fmt.Println("cartID: " + cartID + "  itemID: " + itemID)
	err := cart.DeleteItem(cartID, itemID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}




//get the items in a cart
func getItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.Header().Set("Allow", "GET")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	cartID := r.URL.Query().Get("cartID")
	if cartID == "" {
		http.Error(w, "could not parse cartID", 400)
		return
	}

	items, err := cart.GetItems(cartID)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "%v", items)
}
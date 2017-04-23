package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var(
	userID1 = "130328"

	item1 = "item1"
	item2 = "item2"
	item3 = "item3"

	price1 = "100"
	price2 = "200"
	price3 = "300"

	cartID1 = "9000"
	cartID2 = "9001"

	orderID1 = "5000"
	orderID2 = "5001"

	addr1 = "221B Baker Street"

	payment1 = "1111 2222 3333 4444"
)

func main() {
	location := 8008

	fmt.Printf("Listening on port: %d\n", location)

	http.HandleFunc("/catalog/add", catalogAdd)
	http.HandleFunc("/catalog/delete", catalogDelete)
	http.HandleFunc("/catalog/items", catalogItems)
	http.HandleFunc("/catalog/update", catalogUpdate)
	http.HandleFunc("/catalog/get", catalogGet)
	http.HandleFunc("/catalog/batchGet", catalogBatchGet)

	http.HandleFunc("/cart/add", cartAdd)
	http.HandleFunc("/cart/delete", cartDelete)
	http.HandleFunc("/cart/create", cartCreate)
	http.HandleFunc("/cart/items", cartItems)

	http.HandleFunc("/orders/create", ordersCreate)
	http.HandleFunc("/orders/shipping", ordersShipping)
	http.HandleFunc("/orders/payment", ordersPayment)
	http.HandleFunc("/orders/summary", ordersSummary)



	http.ListenAndServe("localhost:" + strconv.Itoa(location), nil)
}

func catalogGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("catalog get")

	itemID := r.URL.Query().Get("itemID")
	if itemID == "" {
		http.Error(w, "could not parse itemID", 400)
		return
	}

	//write cartID back in response body
	fmt.Fprintf(w, "%s:%s", item1, price1)
}

func catalogBatchGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("catalog batch get")

	items := r.URL.Query().Get("items")
	if items == "" {
		http.Error(w, "could not parse items", 400)
		return
	}
	fmt.Println(items)

	fmt.Fprintf(w, "%s:%s,%s:%s", item1, price1, item2, price2)
}


func catalogAdd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("catalog add")

	itemID := r.PostFormValue("itemID")
	if itemID == "" {
		http.Error(w, "could not parse itemID", 400)
		return
	}

	title := r.PostFormValue("title")
	if title == "" {
		http.Error(w, "could not parse title", 400)
		return
	}

	price := r.PostFormValue("price")
	if price == "" {
		http.Error(w, "could not parse price", 400)
		return
	}

	shippingCost := r.PostFormValue("shippingCost")
	if shippingCost == "" {
		http.Error(w, "could not parse shippingCost", 400)
		return
	}

	//shouldn't return anything
}

func catalogUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("catalog update")

	itemID := r.PostFormValue("itemID")
	if itemID == "" {
		http.Error(w, "could not parse itemID", 400)
		return
	}

	price := r.PostFormValue("price")
	if price == "" {
		http.Error(w, "could not parse price", 400)
		return
	}

	fmt.Printf("Updating item: %s  with price %s\n", itemID, price)

	//shouldn't return anything
}


func catalogDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("catalog delete")

	itemID := r.URL.Query().Get("itemID")
	if itemID == "" {
		http.Error(w, "could not parse itemID", 400)
		return
	}
	fmt.Println(itemID)
}


func catalogItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("catalog items")

	//return IDs of random items
	fmt.Fprintf(w, "%s:%s:%s", item1, item2, item3)
}


func cartCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cart create")

	userID := r.PostFormValue("userID")
	if userID == "" {
		http.Error(w, "could not parse userID", 400)
		return
	}

	fmt.Println("userID: " + userID)


	fmt.Fprintf(w, "%s", cartID1)
}


func cartAdd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cart add")

	cartID := r.PostFormValue("cartID")
	if cartID == "" {
		http.Error(w, "could not parse cartID", 400)
		return
	}

	itemID := r.PostFormValue("itemID")
	if itemID == "" {
		http.Error(w, "could not parse itemID", 400)
		return
	}


	fmt.Println("cartID: " + cartID + "  itemID: " + itemID)

}

func cartDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cart delete")

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
}

func cartItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("cart items")

	cartID := r.URL.Query().Get("cartID")
	if cartID == "" {
		http.Error(w, "could not parse cartID", 400)
		return
	}


	fmt.Println("cartID: " + cartID)

	fmt.Fprintf(w, "%s:%s", item1, item2)
}

func ordersCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ordersCreate")

	if r.Method != "POST" {
		fmt.Println("not POST request")
		w.Header().Set("Allow", "POST")
		http.Error(w, http.StatusText(405), 405)
		return
	}

	formItems := r.PostFormValue("items")
	if formItems == "" {
		http.Error(w, "could not parse formItems", 400)
		return
	}
	fmt.Println("Form items: " + formItems)

	items := strings.Split(formItems, ",")

	for _, item := range items {
		details := strings.Split(item, ":")
		fmt.Println("title: " + details[0])
		fmt.Println("price: " + details[1])
	}

	userID := r.PostFormValue("userID")
	if userID == "" {
		http.Error(w, "could not parse userID", 400)
		return
	}


	fmt.Println("userID: " + userID)

	fmt.Fprintf(w, "%s", orderID1)
}

func ordersShipping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("orders shipping")

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "could not parse userID", 400)
		return
	}


	fmt.Println("userID: " + userID)

	fmt.Fprintf(w, "%s", addr1)
}

func ordersPayment(w http.ResponseWriter, r *http.Request) {
	fmt.Println("orders payment")

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "could not parse userID", 400)
		return
	}


	fmt.Println("userID: " + userID)

	fmt.Fprintf(w, "%s", payment1)
}

func ordersSummary(w http.ResponseWriter, r *http.Request) {
	fmt.Println("orders summary")

	orderID := r.URL.Query().Get("orderID")
	if orderID == "" {
		http.Error(w, "could not parse orderID", 400)
		return
	}


	fmt.Println("orderID: " + orderID)

	fmt.Fprintf(w, "%s:%s:%s", userID1, addr1, payment1)
}
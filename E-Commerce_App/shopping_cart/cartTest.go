package main

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"fmt"
	//"strings"
)

var location = "1339"  //location of the cart service
var userID = "54321"

func main() {
	//testCreate()   First cart: 2596996162  Second cart: 4039455774

	/**
		Test adding to cart
		cartID := "2596996162"
		itemID := "1112"
		testAdd(cartID, itemID)
	*/


	/**
		Test getting items in a cart
		cartID := "2596996162"
		testGetItems(cartID)
	*/

	/**
		Test delete item
		cartID := "2596996162"
		itemID := "1111"
		testDelete(cartID, itemID)
	*/
}


func testCreate() {
	req := "http://localhost:" + location + "/cart/create"
	response, err := http.PostForm(req, url.Values{"userID": {userID}})
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseData))
}


func testAdd(cartID string, itemID string) {
	req := "http://localhost:" + location + "/cart/add"
	response, err := http.PostForm(req, url.Values{"cartID": {cartID}, "itemID": {itemID}})
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseData))
}


func testGetItems(cartID string) {
	req := "http://localhost:" + location + "/cart/items"
	req += "?cartID=" + cartID
	response, err := http.Get(req)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseData))
}

func testDelete(cartID string, itemID string) {
	client := &http.Client{}
	req := "http://localhost:" + location + "/cart/delete"
	req += "?cartID=" + cartID
	req += "&itemID=" + itemID

	request, _ := http.NewRequest("DELETE", req, nil)
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseData))
}
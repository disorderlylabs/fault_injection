package model

import (
	"github.com/mediocregopher/radix.v2/pool"
	"strconv"
	"log"
	"math/rand"
)

var (
	db *pool.Pool
)


//initialize the pool connection to Redis
func Init(redisPort int) {
	var err error
	portNum := strconv.Itoa(redisPort)
	db, err = pool.New("tcp", "localhost:"+portNum, 10)
	if err != nil {
		log.Panic(err)
	}
}


//creates a shopping cart and saves it into Redis
//Return: cartID and error (if any)
func CreateCart(userID string) (uint32, error){
	//get a connection from the pool
	conn, err := db.Get()
	if err != nil {
		return 0, err
	}
	defer db.Put(conn)

	//generate cartID
	cartID := rand.Uint32()

	//start a transaction
	err = conn.Cmd("MULTI").Err
	if err != nil {
		return 0, err
	}

	//add cartID to set of active carts
	err = conn.Cmd("SADD", "carts", cartID).Err
	if err != nil {
		return 0, err
	}

	//each cartID will be associated with a number of items and the userID
	err = conn.Cmd("HSET", cartID, "userID", userID).Err
	if err != nil {
		return 0, err
	}

	err = conn.Cmd("HSET", cartID, "numItems", 0).Err
	if err != nil {
		return 0, err
	}

	//execute the transaction
	err = conn.Cmd("EXEC").Err
	if err != nil {
		return 0, err
	}

	return cartID, nil
}


func AddItem(cartID string, itemID string) error {
	conn, err := db.Get()
	if err != nil {
		return err
	}
	defer db.Put(conn)


	err = conn.Cmd("MULTI").Err
	if err != nil {
		return err
	}

	//the set [cartID]:items contains the items associated with the cartID
	err = conn.Cmd("SADD", cartID + ":items", itemID).Err
	if err != nil {
		return err
	}

	//increment the number of items in the cart by 1
	err = conn.Cmd("HINCRBY", cartID, "numItems", 1).Err
	if err != nil {
		return err
	}

	err = conn.Cmd("EXEC").Err
	if err != nil {
		return err
	}

	return nil
}


func DeleteItem(cartID string, itemID string) error {
	conn, err := db.Get()
	if err != nil {
		return err
	}
	defer db.Put(conn)

	err = conn.Cmd("MULTI").Err
	if err != nil {
		return err
	}

	//remove item from set
	err = conn.Cmd("SREM", cartID + ":items", itemID).Err
	if err != nil {
		return err
	}

	err = conn.Cmd("HINCRBY", cartID, "numItems", -1).Err
	if err != nil {
		return err
	}

	err = conn.Cmd("EXEC").Err
	if err != nil {
		return err
	}

	return nil
}



func GetItems(cartID string) ([]string, error){
	conn, err := db.Get()
	if err != nil {
		return nil, err
	}
	defer db.Put(conn)


	//TODO: maybe error checking on whether the cartID is valid


	//get items in the cart
	items, err := conn.Cmd("SMEMBERS", cartID + ":items").List()
	if err != nil {
		return nil, err
	}

	return items, nil
}




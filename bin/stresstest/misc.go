package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/function61/pyramid-exampleapp-go/types"
	"io/ioutil"
	"math/rand"
	"net/http"
)

func randBetween(min int, max int) int {
	return min + (rand.Int() % (max - min + 1))
}

var products = []string{
	"Regular paper",
	"Premium copy paper",
	"Butt paper",
}

var userIds = []string{
	"6de34313",
	"10a3a454",
	"45206e23",
	"291dd548",
	"bb9bb9e7",
	"e0cae7fd",
	"e1dd2e26",
	"247901e4",
	"0a9d6485",
	"b7d1fbc7",
	"a00ba373",
	"1f535b6e",
	"57dcf2a7",
	"a915400b",
	"5748e781",
	"d530f396",
	"890bbfd4",
}

func placeOrderHttpRequest(baseUrl string) (orderId string, instanceId string, err error) {
	lineItems := []types.OrderPlacementItem{}

	numProducts := randBetween(1, 4)

	for i := 0; i < numProducts; i++ {
		product := products[randBetween(0, len(products)-1)]

		amount := randBetween(1, 20)

		lineItems = append(lineItems, types.OrderPlacementItem{product, amount})
	}

	orderPlacementReq := types.OrderPlacement{
		User:      userIds[randBetween(0, len(userIds)-1)],
		LineItems: lineItems,
	}

	jsonStr, _ := json.Marshal(orderPlacementReq)
	req, err := http.NewRequest("POST", baseUrl+"/command/place_order", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", "", errors.New(fmt.Sprintf("HTTP %s: %s", resp.Status, body))
	}

	// order id as body
	return string(body), resp.Header.Get("X-Instance"), nil
}

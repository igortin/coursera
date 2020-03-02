package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// structure Cart
type Cart struct {
	PaymentApiURL string
}

// method of structure Cart
func (c *Cart) Checkout(id string) (*CheckoutResult, error) {
	// get url for send request
	url := c.PaymentApiURL + "?id=" + id

	// send request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// get response answer
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// create instance  based on structure
	result := &CheckoutResult{}
	// from data to result instance

	err = json.Unmarshal(data, result)

	if err != nil {
		return nil, err
	}
	// retrun result
	return result, nil
}

type CheckoutResult struct {
	Status  int    `json:"status"`
	Balance int    `json:"balance"`
	Err     string `json:"err"`
}

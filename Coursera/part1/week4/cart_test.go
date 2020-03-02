package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type TestCase struct {
	ID      string
	Result  *CheckoutResult
	IsError bool
}

// func test include test cases
func TestCartCheckout(t *testing.T) {
	// list of the test cases
	cases := []TestCase{
		TestCase{
			ID: "42",
			Result: &CheckoutResult{
				Status:  200,
				Balance: 100500,
				Err:     "",
			},
			IsError: false,
		},
		TestCase{
			ID: "100500",
			Result: &CheckoutResult{
				Status:  400,
				Balance: 0,
				Err:     "bad_balance",
			},
			IsError: false,
		},
		TestCase{
			ID:      "__broken_json",
			Result:  nil,
			IsError: true,
		},
		TestCase{
			ID:      "__internal_error",
			Result:  nil,
			IsError: true,
		},
	}
	// htttp test server

	ts := httptest.NewServer(http.HandlerFunc(CheckoutDummy))

	// loop requests

	for caseNum, item := range cases {
		// create in each iteration instance
		c := &Cart{
			PaymentApiURL: ts.URL,
		}
		// request to ts server

		result, err := c.Checkout(item.ID)

		// fmt.Println(result, item.IsError, err)

		// Check variable "err" that is not nil, if it is not nil raise error test
		// due to expected "err" is nill, because item.IsError set false
		if err != nil && !item.IsError {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		// Check variable "err" that is nil, if it is nil raise error test
		// due to expected "err" is not nill, because item.IsError
		if err == nil && item.IsError {
			t.Errorf("[%d] expected error, got nil", caseNum)
		}
		// Check equals
		if !reflect.DeepEqual(item.Result, result) {
			t.Errorf("[%d] wrong result, expected %#v, got %#v",
				caseNum, item.Result, result)
		}
	}
	ts.Close()
}

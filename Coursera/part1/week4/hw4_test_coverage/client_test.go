package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/bradfitz/slice"
)

var (
	ErrorOrderField = "Not correct order_field"
)

// код писать тут
func TestClient(t *testing.T) {
	testCases := []*testCase{
		&testCase{
			Request: &SearchRequest{
				Limit:      2,
				Offset:     1,
				Query:      "Brooks",
				OrderField: "Id",
				OrderBy:    OrderByAsc,
			},
			Response: &SearchResponse{
				Users:    []User{},
				NextPage: false,
			},
			isErr: false,
		},
		&testCase{
			Request: &SearchRequest{
				Limit:      2,
				Offset:     1,
				Query:      "Dickson",
				OrderField: "Name",
				OrderBy:    OrderByAsc,
			},
			Response: &SearchResponse{
				Users:    []User{},
				NextPage: false,
			},
			isErr: false,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))

	client := &SearchClient{
		AccessToken: "Token1235777",
		URL:         ts.URL,
	}

	for _, item := range testCases {
		request := item.Request
		resp, err := client.FindUsers(*request)

		if err != nil && !item.isErr {
			fmt.Println("unexpected error", resp, err.Error())
		}
		if err == nil && item.isErr {
			fmt.Println("unexpected error", resp, err.Error())
		}
	}

	ts.Close()
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	token := r.Header["Accesstoken"][0]

	if !authorise(token) {
		w.WriteHeader(http.StatusForbidden)
	}

	// get request url parameters
	limit := r.FormValue("limit")
	limitInt, _ := strconv.Atoi(limit)
	offset := r.FormValue("offset")
	offsetInt, _ := strconv.Atoi(offset)
	query := r.FormValue("query")
	orderBy := r.FormValue("order_by")

	orderByInt, _ := strconv.Atoi(orderBy)
	orderField := r.FormValue("order_field")

	xmlFile, err := os.Open("dataset.xml")
	if err != nil {
		panic(err.Error())
	}
	defer xmlFile.Close()

	byteArr, err := ioutil.ReadAll(xmlFile)
	xmlData := &Root{}
	err = xml.Unmarshal(byteArr, &xmlData)

	if err != nil {
		panic(err.Error())
	}
	// get data from xml
	res := &SearchResponse{}
	for i := offsetInt; i < len(xmlData.Row); i++ {
		Name := xmlData.Row[i].FirstName + " " + xmlData.Row[i].LastName
		if query != "" {
			if strings.Contains(Name, query) || strings.Contains(xmlData.Row[i].About, query) {
				Id, _ := strconv.Atoi(xmlData.Row[i].ID)
				Age, _ := strconv.Atoi(xmlData.Row[i].Age)
				res.Users = append(res.Users, User{
					Id:     Id,
					Name:   Name,
					Age:    Age,
					About:  xmlData.Row[i].About,
					Gender: xmlData.Row[i].Gender,
				})
			}
		} else {
			Id, _ := strconv.Atoi(xmlData.Row[i].ID)
			Age, _ := strconv.Atoi(xmlData.Row[i].Age)
			res.Users = append(res.Users, User{
				Id:     Id,
				Name:   Name,
				Age:    Age,
				About:  xmlData.Row[i].About,
				Gender: xmlData.Row[i].Gender,
			})
		}
	}
	// check limit of slice and set res.NextPage
	if limitInt < len(res.Users) {
		res.Users = res.Users[:limitInt]
		res.NextPage = true
	} else {
		res.NextPage = false
	}
	// sort []Users by orderField
	err = sortSearchResponse(orderField, res, orderByInt)
	if err != nil {
		fmt.Println(err.Error())
	}

	// preapare response
	responseToClient := []byte{}
	// from structure to json
	responseToClient, err = json.Marshal(res)
	// send response
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(responseToClient))
}

func sortSearchResponse(param string, s *SearchResponse, orderByInt int) error {
	if orderByInt == 0 {
		switch param {
		case "Id":
			slice.Sort(s.Users[:], func(i, j int) bool {
				return s.Users[i].Id < s.Users[j].Id
			})
		case "Age":
			slice.Sort(s.Users[:], func(i, j int) bool {
				return s.Users[i].Age < s.Users[j].Age
			})
		case "":
			fallthrough
		case "Name":
			slice.Sort(s.Users[:], func(i, j int) bool {
				return s.Users[i].Age < s.Users[j].Age
			})
		default:
			return errors.New(ErrorOrderField)
		}
	} else if orderByInt == 1 {
		switch param {
		case "Id":
			slice.Sort(s.Users[:], func(i, j int) bool {
				return s.Users[i].Id > s.Users[j].Id
			})
		case "Age":
			slice.Sort(s.Users[:], func(i, j int) bool {
				return s.Users[i].Age > s.Users[j].Age
			})
		case "":
			fallthrough
		case "Name":
			slice.Sort(s.Users[:], func(i, j int) bool {
				return s.Users[i].Age > s.Users[j].Age
			})
		default:
			return errors.New(ErrorOrderField)
		}
	}
	return nil
}

func authorise(token string) bool {
	if token == "Token1235777" {
		return true
	}
	return false
}

type testCase struct {
	Request  *SearchRequest
	Response *SearchResponse
	isErr    bool
}

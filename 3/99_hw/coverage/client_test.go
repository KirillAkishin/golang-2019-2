package main

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type TestCaseSearchResponse struct {
	ID          string
	AccessToken string
	Request     SearchRequest
	Result      *CheckoutResultClient
}

type TestCaseErrors struct {
	URL         string
	AccessToken string
	Request     SearchRequest
	Result      *CheckoutResultClient
}

type CheckoutResultClient struct {
	clientResult *SearchResponse
	ErrorName    string
}

func SearchServerDummy(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("id")
	switch key {
	case "Dillard":
		result := []byte(`[{"Id":17,"Name":"Dillard Mccoy","Age":36,"About":"Laborum voluptate sit ipsum tempor dolore. Adipisicing reprehenderit minim aliqua est. Consectetur enim deserunt incididunt elit non consectetur nisi esse ut dolore officia do ipsum.\n","Gender":"male"},{"Id":3,"Name":"Everett Dillard","Age":27,"About":"Sint eu id sint irure officia amet cillum. Amet consectetur enim mollit culpa laborum ipsum adipisicing est laboris. Adipisicing fugiat esse dolore aliquip quis laborum aliquip dolore. Pariatur do elit eu nostrud occaecat.\n","Gender":"male"}]`)
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	case "BigLimit":
		result := []byte(`[{"Id":17,"Name":"Dillard Mccoy","Age":36,"About":"Laborum voluptate sit ipsum tempor dolore. Adipisicing reprehenderit minim aliqua est. Consectetur enim deserunt incididunt elit non consectetur nisi esse ut dolore officia do ipsum.\n","Gender":"male"},{"Id":3,"Name":"Everett Dillard","Age":27,"About":"Sint eu id sint irure officia amet cillum. Amet consectetur enim mollit culpa laborum ipsum adipisicing est laboris. Adipisicing fugiat esse dolore aliquip quis laborum aliquip dolore. Pariatur do elit eu nostrud occaecat.\n","Gender":"male"}]`)
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	case "TimeOut":
		time.Sleep(1001 * time.Millisecond)
	case "401":
		w.WriteHeader(http.StatusUnauthorized)
	case "500":
		w.WriteHeader(http.StatusInternalServerError)
	case "400errorInUnmarshal":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(""))
	case "400ErrorBadOrderField":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error": "ErrorBadOrderField"}`))
	case "400UnknownBadRequest":
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"Error": "BadRequest"}`))
	case "BrokenJson":
		result := []byte(``)
		w.WriteHeader(http.StatusOK)
		w.Write(result)
	}
}

func TestFindUserSearchResponse(t *testing.T) {
	cases := []TestCaseSearchResponse{
		TestCaseSearchResponse{ //тест на поиск Dillard
			ID:          "Dillard",
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: &SearchResponse{
					Users: []User{
						User{
							Id:     17,
							Name:   "Dillard Mccoy",
							Age:    36,
							About:  "Laborum voluptate sit ipsum tempor dolore. Adipisicing reprehenderit minim aliqua est. Consectetur enim deserunt incididunt elit non consectetur nisi esse ut dolore officia do ipsum.\n",
							Gender: "male",
						},
					},
					NextPage: true,
				},
				ErrorName: "",
			},
		},
		TestCaseSearchResponse{ // большой Limit
			ID:          "BigLimit",
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      999,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: &SearchResponse{
					Users: []User{
						User{
							Id:     17,
							Name:   "Dillard Mccoy",
							Age:    36,
							About:  "Laborum voluptate sit ipsum tempor dolore. Adipisicing reprehenderit minim aliqua est. Consectetur enim deserunt incididunt elit non consectetur nisi esse ut dolore officia do ipsum.\n",
							Gender: "male",
						},
						User{
							Id:     3,
							Name:   "Everett Dillard",
							Age:    27,
							About:  "Sint eu id sint irure officia amet cillum. Amet consectetur enim mollit culpa laborum ipsum adipisicing est laboris. Adipisicing fugiat esse dolore aliquip quis laborum aliquip dolore. Pariatur do elit eu nostrud occaecat.\n",
							Gender: "male",
						},
					},
					NextPage: false,
				},
				ErrorName: "",
			},
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServerDummy))
	filename = "dataset.xml"

	for caseNum, item := range cases {
		srv := &SearchClient{
			AccessToken: item.AccessToken,
			URL:         ts.URL + "?id=" + item.ID + "&",
		}

		result, err := srv.FindUsers(item.Request)

		if err != nil {
			t.Errorf("[%d] unexpected error: %#v", caseNum, err)
		}
		if !reflect.DeepEqual(item.Result.clientResult, result) {
			t.Errorf("[%d = %#v] wrong result, expected \n%#v,\n got \n%#v", caseNum, item.ID, item.Result.clientResult, result)
		}
	}
	ts.Close()
}

func TestFindUserErrors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServerDummy))
	filename = "dataset.xml"

	cases := []TestCaseErrors{
		TestCaseErrors{ //отрицательный Offset
			URL:         ts.URL,
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     -1,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: nil,
				ErrorName:    "offset must be > 0",
			},
		},
		TestCaseErrors{ //отрицательный Limit
			URL:         ts.URL,
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      -1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: nil,
				ErrorName:    "limit must be > 0",
			},
		},
		TestCaseErrors{ //долгий запрос
			URL:         ts.URL + "?id=TimeOut&",
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: nil,
				ErrorName:    "timeout for limit=2&offset=0&order_by=1&order_field=Name&query=Dillard",
			},
		},
		TestCaseErrors{ //неправильный url
			URL:         "",
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: nil,
				ErrorName:    "unknown error Get ?limit=2&offset=0&order_by=1&order_field=Name&query=Dillard: unsupported protocol scheme \"\"",
			},
		},
		TestCaseErrors{ //неправильный AccessToken
			URL:         ts.URL + "?id=401&",
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: nil,
				ErrorName:    "Bad AccessToken",
			},
		},
		TestCaseErrors{ // ошибка на сервере
			URL:         ts.URL + "?id=500&",
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: nil,
				ErrorName:    "SearchServer fatal error",
			},
		},
		TestCaseErrors{ // плохой запрос и плохая запись ошибки в body
			URL:         ts.URL + "?id=400errorInUnmarshal&",
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: nil,
				ErrorName:    "cant unpack error json: unexpected end of JSON input",
			},
		},
		TestCaseErrors{ // плохой запрос и ErrorBadOrderField
			URL:         ts.URL + "?id=400ErrorBadOrderField&",
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: nil,
				ErrorName:    "OrderFeld Name invalid",
			},
		},
		TestCaseErrors{ // плохой запрос и неизвестная ошибка
			URL:         ts.URL + "?id=400UnknownBadRequest&",
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: nil,
				ErrorName:    "unknown bad request error: BadRequest",
			},
		},
		TestCaseErrors{ // BrokenJson
			URL:         ts.URL + "?id=BrokenJson&",
			AccessToken: correctAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "Dillard",
				OrderField: "Name",
				OrderBy:    1,
			},
			Result: &CheckoutResultClient{
				clientResult: nil,
				ErrorName:    "cant unpack result json: unexpected end of JSON input",
			},
		},
	}

	for caseNum, item := range cases {
		srv := &SearchClient{
			AccessToken: item.AccessToken,
			URL:         item.URL,
		}

		result, err := srv.FindUsers(item.Request)

		if item.Result.ErrorName != err.Error() {
			t.Errorf("[%d] wrong error, expected \n%#v,\n got \n%#v", caseNum, item.Result.ErrorName, err)
		}
		if !reflect.DeepEqual(item.Result.clientResult, result) {
			t.Errorf("[%d] wrong result, expected \n%#v,\n got \n%#v", caseNum, item.Result.clientResult, result)
		}

	}

	ts.Close()
}

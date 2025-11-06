package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

type TestCaseServerSearch struct {
	AccessToken string
	Limit       int
	Offset      int
	Query       string
	OrderField  string
	OrderBy     int
	Result      *CheckoutResult
}

type TestCaseServerErrors struct {
	dataset string
	Result  *CheckoutResult
}

type CheckoutResult struct {
	Status       int
	SearchResult string
}

func TestSearchServer(t *testing.T) {
	searchCases := []TestCaseServerSearch{
		TestCaseServerSearch{ //тест на поиск Dillard
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "Dillard",
			OrderField:  "Name",
			OrderBy:     1,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":17,"Name":"Dillard Mccoy","Age":36,"About":"Laborum voluptate sit ipsum tempor dolore. Adipisicing reprehenderit minim aliqua est. Consectetur enim deserunt incididunt elit non consectetur nisi esse ut dolore officia do ipsum.\n","Gender":"male"}]`,
			},
		},
		TestCaseServerSearch{ //Пустое "Query" и пустое "Name"
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "",
			OrderBy:     1,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":15,"Name":"Allison Valdez","Age":21,"About":"Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n","Gender":"male"}]`,
			},
		},
		TestCaseServerSearch{ //incorrectAccessToken
			AccessToken: incorrectAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "Dillard",
			OrderField:  "Name",
			OrderBy:     1,
			Result: &CheckoutResult{
				Status:       http.StatusUnauthorized,
				SearchResult: "",
			},
		},
		TestCaseServerSearch{ //тест на прямую сортировку по Name
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "Name",
			OrderBy:     1,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":15,"Name":"Allison Valdez","Age":21,"About":"Labore excepteur voluptate velit occaecat est nisi minim. Laborum ea et irure nostrud enim sit incididunt reprehenderit id est nostrud eu. Ullamco sint nisi voluptate cillum nostrud aliquip et minim. Enim duis esse do aute qui officia ipsum ut occaecat deserunt. Pariatur pariatur nisi do ad dolore reprehenderit et et enim esse dolor qui. Excepteur ullamco adipisicing qui adipisicing tempor minim aliquip.\n","Gender":"male"}]`,
			},
		},
		TestCaseServerSearch{ //тест на отсутствие сортировки по Name
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "Name",
			OrderBy:     0,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":0,"Name":"Boyd Wolf","Age":22,"About":"Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n","Gender":"male"}]`,
			},
		},
		TestCaseServerSearch{ //тест на обратную сортировку по Name
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "Name",
			OrderBy:     -1,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":13,"Name":"Whitley Davidson","Age":40,"About":"Consectetur dolore anim veniam aliqua deserunt officia eu. Et ullamco commodo ad officia duis ex incididunt proident consequat nostrud proident quis tempor. Sunt magna ad excepteur eu sint aliqua eiusmod deserunt proident. Do labore est dolore voluptate ullamco est dolore excepteur magna duis quis. Quis laborum deserunt ipsum velit occaecat est laborum enim aute. Officia dolore sit voluptate quis mollit veniam. Laborum nisi ullamco nisi sit nulla cillum et id nisi.\n","Gender":"male"}]`,
			},
		},
		TestCaseServerSearch{ //тест на прямую сортировку по Id
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "Id",
			OrderBy:     1,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":0,"Name":"Boyd Wolf","Age":22,"About":"Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n","Gender":"male"}]`,
			},
		},
		TestCaseServerSearch{ //тест на отсутствие сортировки по Id
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "Id",
			OrderBy:     0,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":0,"Name":"Boyd Wolf","Age":22,"About":"Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n","Gender":"male"}]`,
			},
		},
		TestCaseServerSearch{ //тест на обратную сортировку по Id
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "Id",
			OrderBy:     -1,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":34,"Name":"Kane Sharp","Age":34,"About":"Lorem proident sint minim anim commodo cillum. Eiusmod velit culpa commodo anim consectetur consectetur sint sint labore. Mollit consequat consectetur magna nulla veniam commodo eu ut et. Ut adipisicing qui ex consectetur officia sint ut fugiat ex velit cupidatat fugiat nisi non. Dolor minim mollit aliquip veniam nostrud. Magna eu aliqua Lorem aliquip.\n","Gender":"male"}]`,
			},
		},
		TestCaseServerSearch{ //тест на прямую сортировку по Age
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "Age",
			OrderBy:     1,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":1,"Name":"Hilda Mayer","Age":21,"About":"Sit commodo consectetur minim amet ex. Elit aute mollit fugiat labore sint ipsum dolor cupidatat qui reprehenderit. Eu nisi in exercitation culpa sint aliqua nulla nulla proident eu. Nisi reprehenderit anim cupidatat dolor incididunt laboris mollit magna commodo ex. Cupidatat sit id aliqua amet nisi et voluptate voluptate commodo ex eiusmod et nulla velit.\n","Gender":"female"}]`,
			},
		},
		TestCaseServerSearch{ //тест на отсутствие сортировки по Age
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "Age",
			OrderBy:     0,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":0,"Name":"Boyd Wolf","Age":22,"About":"Nulla cillum enim voluptate consequat laborum esse excepteur occaecat commodo nostrud excepteur ut cupidatat. Occaecat minim incididunt ut proident ad sint nostrud ad laborum sint pariatur. Ut nulla commodo dolore officia. Consequat anim eiusmod amet commodo eiusmod deserunt culpa. Ea sit dolore nostrud cillum proident nisi mollit est Lorem pariatur. Lorem aute officia deserunt dolor nisi aliqua consequat nulla nostrud ipsum irure id deserunt dolore. Minim reprehenderit nulla exercitation labore ipsum.\n","Gender":"male"}]`,
			},
		},
		TestCaseServerSearch{ //тест на обратную сортировку по Age
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "Age",
			OrderBy:     -1,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: `[{"Id":32,"Name":"Christy Knapp","Age":40,"About":"Incididunt culpa dolore laborum cupidatat consequat. Aliquip cupidatat pariatur sit consectetur laboris labore anim labore. Est sint ut ipsum dolor ipsum nisi tempor in tempor aliqua. Aliquip labore cillum est consequat anim officia non reprehenderit ex duis elit. Amet aliqua eu ad velit incididunt ad ut magna. Culpa dolore qui anim consequat commodo aute.\n","Gender":"female"}]`,
			},
		},
		TestCaseServerSearch{ //тест на неправильное поле OrderBy
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "Age",
			OrderBy:     -99999999,
			Result: &CheckoutResult{
				Status:       http.StatusBadRequest,
				SearchResult: `{"Error": "BadRequest"}`,
			},
		},
		TestCaseServerSearch{ //тест на неправильное поле OrderField
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      0,
			Query:       "",
			OrderField:  "About",
			OrderBy:     -1,
			Result: &CheckoutResult{
				Status:       http.StatusBadRequest,
				SearchResult: `{"Error": "ErrorBadOrderField"}`,
			},
		},
		TestCaseServerSearch{ //тест на слишком большой Offset
			AccessToken: correctAccessToken,
			Limit:       1,
			Offset:      9999999,
			Query:       "",
			OrderField:  "Name",
			OrderBy:     -1,
			Result: &CheckoutResult{
				Status:       http.StatusOK,
				SearchResult: "",
			},
		},
	}

	filename = "dataset.xml"

	for caseNum, item := range searchCases {
		searcherParams := url.Values{}
		searcherParams.Add("limit", strconv.Itoa(item.Limit))
		searcherParams.Add("offset", strconv.Itoa(item.Offset))
		searcherParams.Add("query", item.Query)
		searcherParams.Add("order_field", item.OrderField)
		searcherParams.Add("order_by", strconv.Itoa(item.OrderBy))
		url := "http://example.com/api/user" + "?" + searcherParams.Encode()
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Add("AccessToken", item.AccessToken)

		SearchServer(w, req)

		if w.Code != item.Result.Status {
			t.Errorf("[search-%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.Result.Status)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)

		if string(body) != item.Result.SearchResult {
			t.Errorf("[search-%d] wrong Response: got \n%+v, \nexpected \n%+v",
				caseNum, bodyStr, item.Result.SearchResult)
		}
	}

	errorsCases := []TestCaseServerErrors{
		TestCaseServerErrors{
			dataset: "nonexist_file",
			Result: &CheckoutResult{
				Status:       http.StatusInternalServerError,
				SearchResult: "",
			},
		},
		TestCaseServerErrors{
			dataset: "broken_xml.xml",
			Result: &CheckoutResult{
				Status:       http.StatusInternalServerError,
				SearchResult: "",
			},
		},
	}

	searcherParams := url.Values{}
	searcherParams.Add("limit", "1")
	searcherParams.Add("offset", "0")
	searcherParams.Add("query", "")
	searcherParams.Add("order_field", "Name")
	searcherParams.Add("order_by", "0")
	url := "http://example.com/api/user" + "?" + searcherParams.Encode()

	for caseNum, item := range errorsCases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Add("AccessToken", correctAccessToken)

		filename = item.dataset

		SearchServer(w, req)

		if w.Code != item.Result.Status {
			t.Errorf("[errors-%d] wrong StatusCode: got %d, expected %d",
				caseNum, w.Code, item.Result.Status)
		}

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		bodyStr := string(body)

		if string(body) != item.Result.SearchResult {
			t.Errorf("[errors-%d] wrong Response: got \n%+v, \nexpected \n%+v",
				caseNum, bodyStr, item.Result.SearchResult)
		}
	}
}

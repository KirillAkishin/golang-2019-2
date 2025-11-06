package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

var filename string

const correctAccessToken = "GOOD TOKEN"
const incorrectAccessToken = "BAD TOKEN"

type DataForUser struct {
	ID        int    `xml:"id"`
	FirstName string `xml:"first_name"` // json:"name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

type DataForUsers struct {
	Version string        `xml:"version,attr"`
	List    []DataForUser `xml:"row"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.FormValue("limit"))
	offset, err := strconv.Atoi(r.FormValue("offset"))
	query := r.FormValue("query")
	orderField := r.FormValue("order_field")
	orderBy, err := strconv.Atoi(r.FormValue("order_by"))

	if (err != nil) ||
		((orderBy != -1) && (orderBy != 0) && (orderBy != 1)) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("BadRequest")
		w.Write([]byte(`{"Error": "BadRequest"}`))
		return
	}

	if (orderField != "") && (orderField != "Name") && (orderField != "Id") && (orderField != "Age") {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("ErrorBadOrderField")
		w.Write([]byte(`{"Error": "ErrorBadOrderField"}`))
		return
	}

	AccessToken := r.Header.Get("AccessToken")
	if AccessToken != correctAccessToken {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Incorrect AccessToken")
		return
	}

	dataForUsers := new(DataForUsers)

	xmlData, err := ioutil.ReadFile(filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Can't open", filename)
		return
	}
	err = xml.Unmarshal(xmlData, &dataForUsers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Can't unmarshal xml")
		return
	}

	var desiredUsers []User
	var desiredUser User
	flag := false
	if query == "" {
		flag = true
	}
	for i := range dataForUsers.List {
		if (flag) || (dataForUsers.List[i].FirstName == query) || (dataForUsers.List[i].LastName == query) || (dataForUsers.List[i].About == query) {
			desiredUser.About = dataForUsers.List[i].About
			desiredUser.Age = dataForUsers.List[i].Age
			desiredUser.Gender = dataForUsers.List[i].Gender
			desiredUser.Id = dataForUsers.List[i].ID
			desiredUser.Name = dataForUsers.List[i].FirstName + " " + dataForUsers.List[i].LastName
			desiredUsers = append(desiredUsers, desiredUser)
		}
	}

	switch orderBy {
	case -1:
		switch orderField {
		case "Id":
			sort.Slice(desiredUsers, func(i, j int) bool { return desiredUsers[i].Id > desiredUsers[j].Id })
		case "Age":
			sort.Slice(desiredUsers, func(i, j int) bool { return desiredUsers[i].Age > desiredUsers[j].Age })
		case "Name", "":
			sort.Slice(desiredUsers, func(i, j int) bool { return desiredUsers[i].Name > desiredUsers[j].Name })
		}
	case 1:
		switch orderField {
		case "Id":
			sort.Slice(desiredUsers, func(i, j int) bool { return desiredUsers[i].Id < desiredUsers[j].Id })
		case "Age":
			sort.Slice(desiredUsers, func(i, j int) bool { return desiredUsers[i].Age < desiredUsers[j].Age })
		case "Name", "":
			sort.Slice(desiredUsers, func(i, j int) bool { return desiredUsers[i].Name < desiredUsers[j].Name })
		}
	}

	if offset <= len(desiredUsers) {
		desiredUsers = desiredUsers[offset:]
	} else {
		return
	}
	if limit <= len(desiredUsers) {
		desiredUsers = desiredUsers[:limit]
	}

	result, _ := json.Marshal(desiredUsers)
	w.Write(result)

	return

}

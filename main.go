package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Person struct {
	Id        int    `json:"id"`
	FirstName string `json: "first_name"`
	LastName  string `json: "last_name"`
}

type Response struct {
	Persons []Person `json: "persons"`
}

func main() {
	fmt.Println("Hello")
	router := mux.NewRouter()
	router.HandleFunc("/health-check", HealthCheck).Methods("GET")
	router.HandleFunc("/persons", Persons).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(":8080", router)
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Api is up and running")
}

func prepareResponse() []Person {
	var persons []Person

	var person Person
	person.Id = 1
	person.FirstName = "Issac"
	person.LastName = "Newton"
	persons = append(persons, person)

	person.Id = 2
	person.FirstName = "Albert"
	person.LastName = "Einstein"
	persons = append(persons, person)
	return persons
}

func Persons(w http.ResponseWriter, r *http.Request) {
	var response Response
	persons := prepareResponse()
	response.Persons = persons
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}

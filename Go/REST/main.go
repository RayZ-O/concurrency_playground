package main

import (
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
)

type Person struct {
    Name    string    `json:"name"`
    Age     float64   `json:"age"`
    College string    `json:"college"`
}

type Success struct {
    Successful bool   `json:"successful"`
}

var persons map[string]Person = make(map[string]Person)

func RootHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Welcome!\n"))
}

func OwnerHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    owner := Person{"Rui", 24, "UFL"}
    if err := json.NewEncoder(w).Encode(owner); err != nil {
        panic(err)
    }
}

func PostPersonHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    decoder := json.NewDecoder(r.Body)
    var p Person
    if err := decoder.Decode(&p); err != nil {
        panic(err)
    }
    persons[p.Name] = p
    if err := json.NewEncoder(w).Encode(Success{true}); err != nil {
        panic(err)
    }
}

func GetPersonHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    if p, ok := persons[name]; ok {
        if err := json.NewEncoder(w).Encode(p); err != nil {
            panic(err)
        }
    } else {
        if err := json.NewEncoder(w).Encode(Success{false}); err != nil {
            panic(err)
        }
    }
}

func DeletePersonHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    if _, ok := persons[name]; ok {
        delete(persons, name)
    }
    if err := json.NewEncoder(w).Encode(Success{true}); err != nil {
        panic(err)
    }
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", RootHandler).Methods("GET").Headers("Content-Type", "application/text")
    r.HandleFunc("/owner", OwnerHandler).Methods("GET").Headers("Content-Type", "application/json")
    r.HandleFunc("/person", PostPersonHandler).Methods("POST").Headers("Content-Type", "application/json")
    r.HandleFunc("/person/{name}", GetPersonHandler).Methods("GET").Headers("Content-Type", "application/json")
    r.HandleFunc("/person/{name}", DeletePersonHandler).Methods("DELETE").Headers("Content-Type", "application/json")
    http.ListenAndServe(":8000", r)
}

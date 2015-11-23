package main

import (
    "net/http"
    "encoding/json"
    "encoding/hex"
    "crypto/sha256"
    "github.com/gorilla/mux"
    "math/rand"
)

type Person struct {
    Name    string    `json:"name"`
    Age     float64   `json:"age"`
    College string    `json:"college"`
}

type Credential struct {
    Name string
    Password string
}

type Success struct {
    Successful bool   `json:"successful"`
}

type Token struct {
    Token string      `json:"token"`
}


const alphanum = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var persons map[string]Person = make(map[string]Person)
var credentials map[string]Credential = make(map[string]Credential)
var tokens map[string]bool = make(map[string]bool)

func CalSha256(input string) string {
    hash := sha256.New()
    hash.Write([]byte(input))
    return hex.EncodeToString(hash.Sum(nil))
}

func RandString(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = alphanum[rand.Intn(len(alphanum))]
    }
    return string(b)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome!\n"))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    password := vars["password"]
    c := Credential{name, CalSha256(password)}
    credentials[name] = c
    if err := json.NewEncoder(w).Encode(Success{true}); err != nil {
        panic(err)
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    password := vars["password"]
    if c, ok := credentials[name]; ok {
        if CalSha256(password) == c.Password {
            token := RandString(10)
            tokens[token] = true
            if err := json.NewEncoder(w).Encode(Token{token}); err != nil {
                panic(err)
            }
        } else {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("401 Unauthorized\n"))
        }
    } else {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("404 User not found\n"))
    }
}

func Authentication(token string) bool {
    if _, ok := tokens[token]; ok {
        return true
    } else {
        return false
    }
}

func PostPersonHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    if Authentication(vars["token"]) {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        decoder := json.NewDecoder(r.Body)
        var p Person
        if err := decoder.Decode(&p); err != nil {
            panic(err)
        }
        persons[p.Name] = p
        if err := json.NewEncoder(w).Encode(Success{true}); err != nil {
            panic(err)
        }
    } else {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("401 Unauthorized\n"))
    }
}

func GetPersonHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    if Authentication(vars["token"]) {
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
    } else {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("401 Unauthorized\n"))
    }
}

func DeletePersonHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    if Authentication(vars["token"]) {
        name := vars["name"]
        if _, ok := persons[name]; ok {
            delete(persons, name)
        }
        if err := json.NewEncoder(w).Encode(Success{true}); err != nil {
            panic(err)
        }
    } else {
        w.WriteHeader(http.StatusUnauthorized)
        w.Write([]byte("401 Unauthorized\n"))
    }
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", RootHandler).Methods("GET")
    r.HandleFunc("/register/{name}/{password}", RegisterHandler).Methods("POST")
    r.HandleFunc("/login/{name}/{password}", LoginHandler).Methods("GET")
    r.HandleFunc("/person/{token}", PostPersonHandler).Methods("POST").Headers("Content-Type", "application/json")
    r.HandleFunc("/person/{token}/{name}", GetPersonHandler).Methods("GET")
    r.HandleFunc("/person/{token}/{name}", DeletePersonHandler).Methods("DELETE")
    http.ListenAndServe(":8000", r)
}

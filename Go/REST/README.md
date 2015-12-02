# Simple REST server

### Library:
gorilla/mux
website: http://www.gorillatoolkit.org/pkg/mux

### Usage:
(1) Start the server
cd src/rest
go run main.go

(2) Open another terminal
Use curl to test the API

2.1 Get content from root:
curl -X GET -H "Content-Type: application/text"  http://127.0.0.1:8000
Response: Welcome!

2.2 Get info of the owner
curl -X GET -H "Content-Type: application/json"  http://127.0.0.1:8000/owner
Response: {"name":"Rui","age":24,"college":"UFL"}

2.3 Post new person
curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"Bob\",\"age\":21,\"college\":\"UFL\"}"  http://127.0.0.1:8000/person
Response: {"successful":true}

2.4 Get info of posted person
curl -X GET -H "Content-Type: application/json"  http://127.0.0.1:8000/person/Bob
Response: {"name":"Bob","age":21,"college":"UFL"}

2.5 Delete a posted person
curl -X DELETE -H "Content-Type: application/json"  http://127.0.0.1:8000/person/Bob
Response: {"successful":true}

Try to get info again
curl -X GET -H "Content-Type: application/json"  http://127.0.0.1:8000/person/Bob
{"successful":false}

COP 5615 FALL 2015
Rui Zhang 

# REST server with authentication

### Library:  
gorilla/mux  
website: http://www.gorillatoolkit.org/pkg/mux  

### Usage:
(1) Start the server
cd src/
go run main.go

(2) Open another terminal
Use curl to test the API

2.1 Get content from root:
API: GET /
e.g. curl -X GET http://127.0.0.1:8000
Response: Welcome!

2.2 Register new user
API: POST /register/{name}/{password}
e.g. curl -X POST http://127.0.0.1:8000/register/rui/1234
Response: {"successful":true}

2.3 Login and get access token
API: GET /login/{name}/{password}
e.g. curl -X GET http://127.0.0.1:8000/login/rui/1234
Response: {"token":"6FBBD7P95O"}

2.4 Post user info
API: POST /person/{token}
e.g. 
(1) curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"Bob\",\"age\":21,\"college\":\"UFL\"}"  http://127.0.0.1:8000/person/6FBBD7P95O
Response: {"successful":true}
(2) curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"Bob\",\"age\":21,\"college\":\"UFL\"}"  http://127.0.0.1:8000/person/1111111111
Response: 401 Unauthorized

2.5 Get info of posted person
API: POST /person/{token}/{name}
e.g. curl -X GET http://127.0.0.1:8000/person/6FBBD7P95O/Bob
Response: {"name":"Bob","age":21,"college":"UFL"}

2.6 Delete a posted person
API: DELETE /person/{token}/{name}
e.g. curl -X DELETE http://127.0.0.1:8000/person/6FBBD7P95O/Bob
Response: {"successful":true}


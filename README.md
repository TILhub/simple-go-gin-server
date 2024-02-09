# simple-go-gin-server
### open main.go and run the application

#### Sample cURL
```
curl --location 'http://localhost:8080/api/v1/calculator' \
--header 'Content-Type: application/json' \
--data-raw '{
    "data" : {
    "operand1": "1110",
    "operand2": "5",
    "email": "example@gmail.com",
    "operand" : 3
}
}'
```

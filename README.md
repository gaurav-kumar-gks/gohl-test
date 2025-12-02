To run the server

go run cmd/api/main.go

Adding user:

curl --location 'localhost:8080/users' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "gaurav",
    "email": "gaurav.kumar1@gmail.com"
}'

Adding txn:

curl --location --request POST 'localhost:8080/transactions' \
--header 'Content-Type: application/json' \
--data '{
    "user_id": "fb95fbcd-4943-4462-a323-fa00b715eb1a",
    "type": "credit",
    "amount": 12.12,
    "description": "something"
}'

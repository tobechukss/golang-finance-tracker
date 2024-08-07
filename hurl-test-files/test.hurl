POST http://localhost:8085/api/v1/login
```json
{
    "email":"anothertest@email.com",
    "password": "434534"
}
```
HTTP 202
[Captures]
login_token: jsonpath "$['token']"

GET http://localhost:8085/api/v1/record
Authorization:{{login_token}}
HTTP 202
[Asserts]
jsonpath "$.data" exists


POST http://localhost:8085/api/v1/record
Authorization:{{login_token}}
```json
{
    "description": "test description from hurl tests",
    "category": "new category",
    "amount": 23434
}
```
HTTP 201
[Captures]
record_id: jsonpath "$['record_id']"
[Asserts]
jsonpath "$.record_id" exists

GET http://localhost:8085/api/v1/record/{{record_id}}
Authorization:{{login_token}}
HTTP 202
[Asserts]
jsonpath "$.data" exists

curl -d '{"userId": "asdf"}' http://localhost:8080/v1/demo/is-new-user

curl -d '{"userId": "asdf", "userPassword": "fddsa"}' -X POST http://localhost:8080/v1/demo/register

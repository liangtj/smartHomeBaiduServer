curl -d '{"userId": "asdf"}' -X GET http://localhost:8080/v1/demo/is-new-user

curl -d '{"userId": "asdf", "userPassword": "fddsa"}' -X POST http://localhost:8080/v1/demo/register


# `content-type` required because of the `c.Bind`/`c.ShouldBind` (auto-detect MIME-type)
curl -H 'content-type: application/json' -d '{"userId": "asdf"}' -X GET http://localhost:8080/v1/demo/is-new-user

curl -H 'content-type: application/json' -d '{"userId": "asdf", "userPassword": "fddsa"}' -X POST http://localhost:8080/v1/demo/register


curl -H 'content-type: application/json'  -d '{"code": "asdf"}' -X GET http://localhost:8080/v1/test-login

curl -H 'content-type: application/json'  -H 'Wx-Session-Token: XXX0' -X GET http://localhost:8080/v1/test-ru

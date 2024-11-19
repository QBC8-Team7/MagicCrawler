# Bootup

1. run `go mod tidy`  for install packages
2. install docker
3. install dockerized postgres
4. install dockerized redis
5. get bot id token from botfather and set in config
6. run project by `make run-app`


## Database
write your query inside `pkg/db/query` then run `make sqlc-gen` to generate repositories. you can import generated function to interact with db.


if you made any changes on db scheama please make sure that you create a migration files by `make new-migrate`.

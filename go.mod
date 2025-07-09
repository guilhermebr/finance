module finance

go 1.24.2

require (
	github.com/ardanlabs/conf/v3 v3.8.0
	github.com/go-chi/chi/v5 v5.2.1
	github.com/go-chi/cors v1.2.1
	github.com/go-chi/render v1.0.3
	github.com/gofrs/uuid/v5 v5.3.2
	github.com/gorilla/mux v1.8.1
	github.com/guilhermebr/gox/logger v0.0.0-20250531115130-f761d05ebb90
	github.com/guilhermebr/gox/monetary v0.0.0-20250709010243-aa7489a9f384
	github.com/guilhermebr/gox/postgres v0.0.0-20250531115130-f761d05ebb90
	github.com/jackc/pgx/v5 v5.7.5
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.10.0
	github.com/swaggo/http-swagger/v2 v2.0.2
	github.com/swaggo/swag v1.16.4
)

//replace github.com/guilhermebr/gox/postgres v0.0.0 => ../gox/postgres
//replace github.com/guilhermebr/gox/logger v0.0.0 => ../gox/logger

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/ajg/form v1.5.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-openapi/jsonpointer v0.21.1 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/spec v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/swaggo/files/v2 v2.0.2 // indirect
	golang.org/x/crypto v0.39.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/tools v0.34.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

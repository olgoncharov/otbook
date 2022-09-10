module github.com/olgoncharov/otbook

go 1.18

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/gojuno/minimock/v3 v3.0.10
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/gorilla/mux v1.8.0
	github.com/ilyakaznacheev/cleanenv v1.3.0
	github.com/json-iterator/go v1.1.12
	github.com/rs/zerolog v1.27.0
	github.com/stretchr/testify v1.8.0
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
)

require (
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rs/cors v1.8.2 // indirect
	golang.org/x/sys v0.0.0-20220731174439-a90be440212d // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

replace github.com/olgoncharov/otbook/internal => ./internal

export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

.PHONY: install-minimock
install-minimock:
	go install github.com/gojuno/minimock/v3/cmd/minimock

.PHONY: generate-mocks
generate-mocks:
	minimock -i github.com/olgoncharov/otbook/internal/usecase/access/command/login.* -o ./internal/usecase/access/command/login/mocks -s "_mock.go"
	minimock -i github.com/olgoncharov/otbook/internal/usecase/access/command/refresh_token.* -o ./internal/usecase/access/command/refresh_token/mocks -s "_mock.go"
	minimock -i github.com/olgoncharov/otbook/internal/usecase/profile/command/create.* -o ./internal/usecase/profile/command/create/mocks -s "_mock.go"
	minimock -i github.com/olgoncharov/otbook/internal/usecase/profile/command/update.* -o ./internal/usecase/profile/command/update/mocks -s "_mock.go"
	minimock -i github.com/olgoncharov/otbook/internal/usecase/profile/query/list.* -o ./internal/usecase/profile/query/list/mocks -s "_mock.go"
	minimock -i github.com/olgoncharov/otbook/internal/usecase/profile/query/user_profile.* -o ./internal/usecase/profile/query/user_profile/mocks -s "_mock.go"
	minimock -i github.com/olgoncharov/otbook/internal/usecase/friends/command/become_friends.* -o ./internal/usecase/friends/command/become_friends/mocks -s "_mock.go"
	minimock -i github.com/olgoncharov/otbook/internal/usecase/friends/query/list.* -o ./internal/usecase/friends/query/list/mocks -s "_mock.go"

.PHONY: install-goose
install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@v3.6.1

.PHONY: create-migration
create-migration:
	goose -dir ./migrations create ' ' sql

.PHONY: migrations-up
migrations-up:
	goose -dir ./migrations mysql "admin:admin@/otbook?parseTime=true" up

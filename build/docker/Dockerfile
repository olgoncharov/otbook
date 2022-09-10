FROM golang:1.18-alpine

WORKDIR /app
RUN mkdir bin
ENV GOBIN /app/bin

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.6.1

COPY build/docker/entrypoint.sh .
COPY cmd/ cmd/
COPY config/ config/
COPY migrations/ migrations/
COPY internal/ internal/

RUN go build -o /bin/otbook ./cmd/main.go

EXPOSE 8000

ENTRYPOINT [ "/app/entrypoint.sh" ]
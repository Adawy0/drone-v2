FROM golang:1.18 as development

ARG VERSION=dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go install github.com/cespare/reflex@latest

EXPOSE 4000

CMD reflex -g '*.go' go run main.go --start-service
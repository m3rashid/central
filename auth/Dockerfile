FROM golang:1.21.5-alpine

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

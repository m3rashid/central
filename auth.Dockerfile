FROM golang:1.21.5-alpine

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app
RUN mkdir auth internal

COPY internal/go.mod ./internal/go.sum ./internal/
RUN cd internal && go mod download && cd ..
COPY ./internal ./internal

COPY ./auth/go.mod ./auth/go.sum ./auth/
RUN cd auth && go mod download && cd ..
COPY ./auth ./auth

WORKDIR /app/auth

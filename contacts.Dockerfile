FROM golang:1.21.5-alpine

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app
RUN mkdir contacts internal

COPY internal/go.mod ./internal/go.sum ./internal/
RUN cd internal && go mod download && cd ..
COPY ./internal ./internal

COPY ./contacts/go.mod ./contacts/go.sum ./contacts/
RUN cd contacts && go mod download && cd ..
COPY ./contacts ./contacts

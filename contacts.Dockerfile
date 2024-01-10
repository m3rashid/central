FROM golang:1.21.5-alpine

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

WORKDIR /app

COPY . .

RUN cd internal && go mod download && cd ..
RUN cd contacts && go mod download && cd ..

WORKDIR /app/contacts

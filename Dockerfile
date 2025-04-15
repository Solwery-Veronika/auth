FROM golang:1.23.4-alpine

COPY . .

RUN go build -o main ./cmd/service/main.go

CMD ./main 
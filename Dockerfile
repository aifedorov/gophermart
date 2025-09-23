FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN (cd cmd/gophermart && go build -buildvcs=false -o gophermart)

EXPOSE 8080

CMD ["./cmd/gophermart/gophermart"]

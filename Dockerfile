FROM golang:1.23.2

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y libxml2-dev postgresql-client

COPY . .

RUN go build -o /bank-service ./src/main.go

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 8080

CMD ["/bank-service"]

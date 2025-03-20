FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod .
# COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o urlshortener .

EXPOSE 8080

CMD ["./urlshortener", "-storage=memory", "-port=8080"]

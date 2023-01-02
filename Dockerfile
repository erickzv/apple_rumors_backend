FROM golang:1.19-bullseye

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o server

CMD ["./server"]

FROM golang:bookworm as build

WORKDIR /go/src/app

COPY . .

RUN go mod download

RUN go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=buil /go/bin/app /

COPY src/ src/

CMD ["./app"]

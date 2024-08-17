FROM golang:latest AS build

WORKDIR /build

COPY . .

RUN go build

FROM golang:latest

WORKDIR /app

COPY --from=build /build/http_server .

RUN chmod +x /app/http_server

CMD ["go", "run", "http_server"]

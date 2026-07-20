FROM golang:1.26-alpine as build

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go build -o bookmark-service cmd/api/main.go

FROM alpine:3.24.1

WORKDIR /app

COPY --from=build /opt/app/bookmark-service /app/bookmark-service
COPY --from=build /opt/app/docs /app/docs

CMD ["/app/bookmark-service"]

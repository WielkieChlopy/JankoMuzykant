FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . . 
#RUN CGO_ENABLED=0 GOOS=linux go build -o cmd/server.go
RUN go build -o app cmd/server.go

FROM golang:1.21 as runner

COPY --from=builder /app/app /app/app

EXPOSE 8080

CMD ["/app/app"]
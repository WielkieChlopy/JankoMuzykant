FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN mkdir /app/bin && curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /app/bin/yt-dlp && chmod a+rx /app/bin/yt-dlp


COPY . . 
#RUN CGO_ENABLED=0 GOOS=linux go build -o cmd/server.go
RUN go build -o app cmd/server.go

FROM golang:1.23 as runner

COPY --from=builder /app/bin/yt-dlp /app/bin/yt-dlp
COPY --from=builder /app/app /app/app

EXPOSE 8080

CMD ["/app/app"]

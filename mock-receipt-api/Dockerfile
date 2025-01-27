FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM scratch

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 1234

CMD ["./main"]
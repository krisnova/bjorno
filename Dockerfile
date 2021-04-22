FROM golang:1-alpine
WORKDIR /go/src/github.com/kris-nova/bjorno
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app cmd/main.go

# --- Multistage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /app /app
CMD ["/app"]
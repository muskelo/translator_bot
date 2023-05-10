from golang:1.20.3 as build-stage
workdir /build
copy . .
run CGO_ENABLED=0 go build -o translator_bot ./cmd/bot

from alpine:latest
run apk --no-cache add ca-certificates
workdir /app
copy --from=build-stage /build/translator_bot .
cmd ["/app/translator_bot"]

from golang:1.20.3 as build-stage
workdir /build
copy . .
run CGO_ENABLED=0 go build -o bot ./cmd/bot
run CGO_ENABLED=0 go build -o migrate ./cmd/migrate

from alpine:latest
run apk --no-cache add ca-certificates
workdir /app
copy --from=build-stage /build/bot .
copy --from=build-stage /build/migrate .
cmd ["/app/bot"]

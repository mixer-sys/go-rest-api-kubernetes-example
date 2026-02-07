FROM alpine:3.20 as root-certs
RUN apk add --no-cache ca-certificates && \
    addgroup -g 1001 app && \
    adduser -u 1001 -D -G app -h /home/app app

FROM golang:1.25-alpine as builder
WORKDIR /youtube-api-files
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o youtube-stats ./cmd/app

FROM scratch as final
COPY --from=root-certs /etc/passwd /etc/group /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=builder /youtube-api-files/youtube-stats /youtube-stats
USER app
ENTRYPOINT ["/youtube-stats"]

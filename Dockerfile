FROM alpine:3.16 as root-certs
RUN apk add -u --no-cache ca-certificates
RUN  addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app /home/app

FROM golang:1.17 as builder
WORKDIR /youtube-api-files
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o ./youtube-stats ./cmd/app/./...

FROM scratch as final
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ca-certs.crt /etc/ssl
COPY --chown=1001:1001 --from=builder /youtube-api/files/youtube-stats /youtube-stats
USER app
ENTRYPOINT ["/youtube-stats"]
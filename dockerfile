FROM golang:1.24-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /out/api ./cmd/scanner-api

FROM alpine:3.20

RUN apk add --no-cache clamav clamav-daemon tzdata ca-certificates \
  && mkdir -p /run/clamav /var/lib/clamav /var/log/clamav /app \
  && chown -R clamav:clamav /run/clamav /var/lib/clamav /var/log/clamav

COPY docker/clamd.conf /etc/clamav/clamd.conf
COPY docker/freshclam.conf /etc/clamav/freshclam.conf
COPY docker/entrypoint.sh /entrypoint.sh
RUN sed -i 's/\r$//' /etc/clamav/clamd.conf /etc/clamav/freshclam.conf /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY --from=build /out/api /app/api

EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
ENV GO_ENV=production
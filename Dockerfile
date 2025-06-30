# syntax=docker/dockerfile:1.4
FROM golang:1.24-bullseye AS builder

WORKDIR /go/src/github.com/ory/kratos

COPY go.mod go.sum ./
COPY internal/client-go/go.* internal/client-go/
RUN go mod download

COPY . .

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

RUN go build \
  -ldflags="-X 'github.com/ory/kratos/driver/config.Version=${VERSION}' -X 'github.com/ory/kratos/driver/config.Date=${BUILD_DATE}' -X 'github.com/ory/kratos/driver/config.Commit=${COMMIT}'" \
  -o /usr/bin/kratos

# --- runner ---
FROM debian:bullseye

# Install ca-cert + gettext (for envsubst)
RUN apt-get update && \
    apt-get install -y ca-certificates gettext && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/bin/kratos /usr/bin/kratos

# Use envsubst from gettext
COPY config/kratos.yml /etc/config/kratos_template.yml
COPY stub /stub
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 4433 4434
ENTRYPOINT ["/entrypoint.sh"]
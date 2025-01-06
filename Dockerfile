FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server .

################################################################################
FROM alpine:latest AS final

# Install any runtime dependencies that are needed to run your application.
# Leverage a cache mount to /var/cache/apk/ to speed up subsequent builds.
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

COPY --from=builder /bin/server /bin/
COPY config.yaml /config.yaml

EXPOSE 8080

ENTRYPOINT [ "/bin/server" ]

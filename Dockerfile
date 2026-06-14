FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

ARG VERSION=dev
ARG COMMIT=unknown

RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT}" \
    -o hermes .

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/hermes /usr/local/bin/hermes

ENTRYPOINT ["hermes"]

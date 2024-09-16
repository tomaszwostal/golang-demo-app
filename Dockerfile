FROM golang:alpine AS builder

RUN apk update && apk add --no-cache gcc libc-dev make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .
# Run audit and test
RUN make audit

# Build the binary
RUN make build

FROM alpine

ENV DB_HOST= \
    DB_PORT= \
    DB_USER= \
    DB_PASSWORD= \
    DB_NAME= \
    DB_SSL=

WORKDIR /app
COPY --from=builder /tmp/bin/demo .

RUN apk --no-cache add ca-certificates tzdata

ENTRYPOINT ["/app/demo"]

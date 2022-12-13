# builder image
FROM golang:1.17-alpine3.14 AS builder
RUN apk add build-base
WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/main ./cmd/server/

# generate clean, final image for end users
FROM alpine:3.14
COPY --chown=65534:65534 --from=builder /go/bin/main /go/src/app/.env /go/src/app/aliasWords.txt ./

USER 65534

# executable
ENTRYPOINT [ "./main" ]

############################
# STEP 1 build executable binary
############################
FROM golang:1.18.1-alpine3.15 AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
# Fetch dependencies.
# Using go get.
RUN go mod download
RUN go mod verify
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/bot
############################
# STEP 2 build a small image
############################
FROM alpine
# To load timezone from TZ env
RUN apk update && apk add --no-cache tzdata
# Copy our static executable.
COPY --from=builder /go/bin/bot /go/bin/bot
# Run the bot binary.
ENTRYPOINT ["/go/bin/bot"]
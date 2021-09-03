# syntax=docker/dockerfile:1

##
## Build
##

FROM golang:1.17-alpine AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /ecommerceApi

##
## Deploy
##

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /ecommerceApi /ecommerceApi

EXPOSE 8080

# Run
CMD [ "./ecommerceApi" ]
FROM golang:alpine as builder


RUN apk update && apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files 
COPY go.mod go.sum ./
COPY Makefile ./

RUN go mod download 

# Copy the source from the current directory to the working Directory inside the container 
COPY . .

RUN apk update
RUN apk add make # install make
# Build the Go app
RUN make build BUILD_OUTPUT=/app/bin/server

# Start a new stage from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates


RUN apk --no-cache add curl

RUN apk add --update curl && \
    rm -rf /var/cache/apk/*

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/bin/server .
COPY --from=builder /app/.env .       

# Expose port 8080 to the outside world
# EXPOSE 8080

#Command to run the executable
CMD ["./server"]
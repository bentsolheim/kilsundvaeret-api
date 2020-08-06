FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    OUTPUT_DIR=_out

# Move to working directory
WORKDIR /work

# Install ca-certificates that can be copied into scratch later
RUN apk add -U --no-cache ca-certificates

RUN apk --no-cache add tzdata

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY ./ ./

# Build the application
RUN rm -rf $OUTPUT_DIR && mkdir $OUTPUT_DIR \
    && go build -o $OUTPUT_DIR ./...

# Build a small image
FROM scratch

ENV TZ=Europe/Oslo

ENV SERVER_PORT=9010 \
    GIN_MODE=release

EXPOSE 9010

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /work/_out/app /
COPY --from=builder /work/migrations /migrations

# Command to run
ENTRYPOINT ["/app"]

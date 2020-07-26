FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    OUTPUT_DIR=_out

# Move to working directory
WORKDIR /work

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

ENV SERVER_PORT=9010 \
    GIN_MODE=release

EXPOSE 9010

COPY --from=builder /work/_out/api /

# Command to run
ENTRYPOINT ["/api"]
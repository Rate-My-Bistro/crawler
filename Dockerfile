FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download necessacry go dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application and name it 'main'
# and strip debug symbols to reduce the size (38->33 MB)
RUN go build -ldflags '-s' -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to /dist folder
RUN cp /build/main .

# Build a small image
FROM scratch

# copy applicatoin binary 'main' to the container
COPY --from=builder /dist/main /

# Command to run
ENTRYPOINT ["/main"]

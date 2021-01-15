FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Change to working directory /build
WORKDIR /build

# Copy and download necessacry go dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application and name it 'app'
# strip debug symbols to reduce the size (reduces the size by ~13%)
RUN go build -ldflags '-s' -o app .

############################################################
# Build a small image
FROM scratch

# copy application binary 'app' to the container
COPY --from=builder /build/app /

# copy swagger doc
COPY restapi/docs/swagger.json /

# Command to run
ENTRYPOINT ["/app"]

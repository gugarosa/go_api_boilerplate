# Imports a GO alpine image as the builder
FROM golang:1.14-alpine AS builder

# Sets environment variables necessary for building
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Creates the source's directory
RUN mkdir -p /src

# Sets the work directory to source's folder
WORKDIR /src

# Copy files into source's folder
COPY ./src .

# Install the dependencies
RUN go mod download

# Builds the application
RUN go build -o api api.go

# Creates a smaller-sized image
FROM scratch

# Sets the work directory to distribution's folder
WORKDIR /dist

# Copies the binary from the builder
COPY --from=builder /src/api ./

# Runs the application
CMD ["./api"]
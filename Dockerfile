FROM golang:1.14-alpine

# Creates the application's directory
RUN mkdir -p /api

# Sets the work directory to application's folder
WORKDIR /api

# Copy files into application's folder
COPY . .

# Install the dependencies
RUN go mod download

# Builds the application
RUN CGO_ENABLED=0 GOARCH=amd64 go build -o app app.go

# Runs the application
# CMD ["./app"]
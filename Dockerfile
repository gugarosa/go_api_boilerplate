FROM golang:1.14-alpine

# Creates the application's directory
RUN mkdir -p /vivere/api

# Sets the work directory to application's folder
WORKDIR /vivere/api

# Copy files into application's folder
COPY . .

# Install the dependencies
RUN go mod download

#
RUN go get github.com/codegangsta/gin

#
EXPOSE 3001

# 
CMD ["gin", "--port", "8080", "run", "app.go"]
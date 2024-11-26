# Base image
FROM golang:latest

# Working directory
WORKDIR /app

# Copy the source code
COPY . ./
# Install necessary modules
RUN go mod download

# Copy the content of docker env file
RUN cat .env_docker > .env

# Install python and necessary packages
RUN apt-get update
RUN apt-get -y install python3
RUN apt-get -y install python3-setuptools
RUN apt-get -y install python3-pip
RUN pip install cryptography --break-system-packages

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/Safety-Tracker-Server

# Command to run when start the container
CMD ["/build/Safety-Tracker-Server"]
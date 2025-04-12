FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project directory into the container
COPY . .

CMD ["/bin/bash"]

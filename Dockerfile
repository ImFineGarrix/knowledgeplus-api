# Dockerfile.gin

FROM golang:latest

WORKDIR /app

COPY go.mod .
COPY go.sum .
ARG ENV

# Set environment
ENV ENV=${ENV}

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8081

CMD ["./main"]

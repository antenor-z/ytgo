FROM golang:latest

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ytgo

EXPOSE 5100

CMD ["./ytgo"]

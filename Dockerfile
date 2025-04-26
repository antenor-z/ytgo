FROM golang:latest

RUN apt-get update && apt-get install -y python3 python3-pip
RUN pip install yt-dlp --break-system-packages
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ytgo

EXPOSE 5100

CMD ["./ytgo"]

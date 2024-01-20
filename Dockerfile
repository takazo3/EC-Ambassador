FROM golang:1.20

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod tidy

COPY . .

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./
RUN go mod tidy

CMD ["air", "-c", ".air.toml"]
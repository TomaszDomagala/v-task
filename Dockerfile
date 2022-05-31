FROM golang:1.18-alpine

# TODO: image can be smaller by multi-stage build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o ./emailmanager ./cmd/emailmanager

EXPOSE 8080

CMD ["./emailmanager"]

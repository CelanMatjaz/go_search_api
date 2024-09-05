FROM golang:1.23-alpine

WORKDIR /backend

RUN go install github.com/air-verse/air@latest

RUN apk add vim neovim

COPY .air.toml ./

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]

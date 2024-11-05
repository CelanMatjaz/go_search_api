FROM golang:1.23-alpine

RUN apk add vim neovim make postgresql

RUN addgroup -S testing_group && adduser -S testing_user -G testing_group
USER testing_user
WORKDIR /testing_user

COPY go.mod go.sum ./
RUN go mod download

CMD ["make", "test"]

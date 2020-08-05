FROM golang:1.14.6-buster as base

ENV GO111MODULE=auto
WORKDIR /app

FROM base as modules

COPY go.mod go.sum ./
RUN go mod download

FROM modules as build

COPY cmd ./cmd
COPY pkg ./pkg
RUN go build -a -o bin/bork ./cmd/bork

FROM modules as dev

RUN go get golang.org/x/tools/gopls@latest
RUN go get github.com/cosmtrek/air

CMD ["./bin/bork"]

FROM build as test

CMD ["go", "test", "./pkg/..."]

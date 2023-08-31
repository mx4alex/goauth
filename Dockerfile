FROM golang:1.20

WORKDIR /usr/src/app

# install dependencies
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./

# build app
RUN go build -o goauth ./cmd/app/main.go

CMD ["./goauth"]
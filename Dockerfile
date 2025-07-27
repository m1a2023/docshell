FROM golang:1.24 AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./ 
RUN go mod download

COPY . . 
RUN CGO_ENABLED=0 go build -o docshell ./cmd

CMD [ "./docshell" ]
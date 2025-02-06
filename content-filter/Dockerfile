FROM golang:alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o content-filter ./cmd/main.go

CMD ["./content-filter"]
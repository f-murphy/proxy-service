FROM golang:alpine

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o anomaly-detection ./cmd/main.go

CMD ["./anomaly-detection"]
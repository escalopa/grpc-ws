FROM golang:alpine AS builder
WORKDIR /go/src/github.com/escalopa/grpc-ws
COPY . .
RUN go mod download
RUN go build -o /go/bin/grpc-ws ./cmd/.

FROM alpine:3.19 AS production
COPY --from=builder /go/bin/grpc-ws /go/bin/grpc-ws
COPY /go/bin/grpc-ws /go/bin/grpc-ws
COPY ./.cataystgo/ /go/bin/.catalystgo/
ENTRYPOINT ["/go/bin/grpc-ws"]

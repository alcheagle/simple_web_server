FROM golang:1.9.0 AS builder
WORKDIR /go/src/github.com/alcheagle/example_server
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build

FROM alpine
COPY --from=builder /go/src/github.com/alcheagle/example_server .
EXPOSE 8080
ENTRYPOINT ["/example_server"]

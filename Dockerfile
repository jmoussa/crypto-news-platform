FROM golang:1.18beta2
COPY . /go/src/
WORKDIR /go/src/

RUN go clean -modcache
RUN go mod tidy 
EXPOSE 3000
ENTRYPOINT go run main.go -api
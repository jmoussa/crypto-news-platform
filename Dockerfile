FROM golang
COPY . /go/src/
WORKDIR /go/src/
RUN go get .
ENTRYPOINT go run main.go -api
EXPOSE 3000
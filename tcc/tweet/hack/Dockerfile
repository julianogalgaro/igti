# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.8

RUN go get github.com/julienschmidt/httprouter
RUN go get gopkg.in/mgo.v2
RUN go get github.com/dghubble/go-twitter/twitter
RUN go get github.com/dghubble/oauth1
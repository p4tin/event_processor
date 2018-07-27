# create image from the official Go image
FROM golang:latest

# install dependancies
#RUN apk add --update tzdata bash wget curl git;

# Create binary directory, install fresh
RUN mkdir -p $$GOPATH/bin && \
    go get github.com/pilu/fresh

# define work directory
ADD . /go/src/event_processor/
WORKDIR /go/src/event_processor/cmd/

# serve the app
CMD fresh

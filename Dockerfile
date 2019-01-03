##################################
FROM golang

# Build Executable Binary

ADD . /go/src/github.com/CelesteComet/celeste-auth-service
WORKDIR /go/src/github.com/CelesteComet/celeste-auth-service

# Fetch Dependencies
RUN go get 

# Build Binary
RUN go build .

# Run server when container is run

CMD ./celeste-auth-service

# Expose port 6800 of container

EXPOSE 6800

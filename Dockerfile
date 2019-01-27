##################################
FROM golang

# Build Executable Binary

ADD . /go/src/github.com/CelesteComet/celeste-auth-service
WORKDIR /go/src/github.com/CelesteComet/celeste-auth-service

# Fetch Dependencies
RUN go get 

# Build Binary
RUN go install

# Run server when container is run

CMD /go/bin/celeste-auth-service

# Expose port 1337 of container

EXPOSE 1337

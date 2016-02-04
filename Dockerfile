# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/RichardKnop/recall

ENV GO15VENDOREXPERIMENT 1
WORKDIR /go/src/github.com/RichardKnop/recall

# Build the recall command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/RichardKnop/recall

# Copy the docker-entrypoint.sh script and use it as entrypoint
COPY ./docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]

# Document that the service listens on port 8080.
EXPOSE 8080

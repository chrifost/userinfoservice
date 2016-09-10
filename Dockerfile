# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go
# Use 'onbuild' variant - this automatically copies source, buids and configures
# for startup
FROM golang:onbuild

# Copy the local package files to the container's workspace
#ADD . /go/src/userinfoservice

# Build the outyet command inside the continer.
# (You may fetch or manage depdencies here,
# either manually or with a tool like "godep".)
#RUN go install userinfoservice

# Run the jsonparse command by default when the container starts.
#ENTRYPOINT /go/bin/userinfoservice

# Document that the service listens on port 8080
EXPOSE 8080

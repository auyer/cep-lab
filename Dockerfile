#Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/latitude-RESTsec-lab/api-gingonic

# # Build the outyet command inside the container.
# # (You may fetch or manage dependencies here,
# # either manually or with a tool like "godep".)
#RUN HTTPS_PROXY=https://10.30.0.10:3128 go get -u all
RUN HTTPS_PROXY=https://10.30.0.10:3128 go get github.com/lib/pq github.com/go-gorp/gorp github.com/gin-gonic/gin
RUN HTTPS_PROXY=https://10.30.0.10:3128 go install github.com/latitude-RESTsec-lab/api-gingonic

# Run the outyet command by default when the container starts.
#RUN cp /go/src/github.com/latitude-RESTsec-lab/api-gingonic/config.json /go/bin/config.json

ADD ./config.json .
ADD ./devssl ./devssl

EXPOSE 443
ENTRYPOINT /go/bin/api-gingonic
# Document that the service listens on port 8080


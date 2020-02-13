#http://blog.traintracks.io/monorepo-shared-code-and-isolated-fast-docker-builds/
#https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
###################################
# STEP 1 build executable binary #
##################################
FROM golang:alpine as builder
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

ENV SRC_DIR=/go/src/github.com/sharkyze/lbc
COPY . $SRC_DIR
WORKDIR $SRC_DIR/cmd/api

# enable go modules support
ARG GO111MODULE=on
ARG VERSION="v0.0.0"
# Fetch dependencies.
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/api .


##############################
# STEP 2 build a small image #
##############################
FROM scratch
# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwdw
# Copy our static executable.
COPY --from=builder /go/bin/api .
ENTRYPOINT ["./api"]

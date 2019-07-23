FROM golang:1.12 AS builder

RUN apt-get update && apt-get install protobuf-compiler -y

RUN go get -u \
	github.com/golang/dep/cmd/dep \
	github.com/golang/protobuf/protoc-gen-go \
	github.com/mwitkow/go-proto-validators/protoc-gen-govalidators \
	github.com/twitchtv/twirp/protoc-gen-twirp \
	&& cd /go/src/github.com/twitchtv/twirp && git checkout v6.0.0-prerelease.alpha1 && go install ./protoc-gen-twirp

ARG package=/go/src/github.com/datalinkE/yet-another-chat
ADD . ${package}
WORKDIR ${package}

RUN make codegen

RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /service ./cmd/service/main.go

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /service ./
ENTRYPOINT ["./service"]
EXPOSE 9000

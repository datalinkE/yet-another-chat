IMAGE?=yet-another-chat

clean:
	rm -rf rpc

codegen: clean
	mkdir -p ./rpc
	cd proto && \
		protoc \
		--proto_path=${GOPATH}/src \
		--proto_path=. \
		--go_out=${GOPATH}/src \
		--govalidators_out=${GOPATH}/src \
		--twirp_out=${GOPATH}/src \
		./*.proto
	cd ./rpc && ls ./*.pb.go | grep -v validator | xargs -I @ protoc-go-inject-tag -input=@

run-local:
	go run cmd/service/main.go

run: build
	docker run -ti --rm -p 9000:9000 ${IMAGE}:latest
build:
	docker build -t ${IMAGE} .

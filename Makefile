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

run:
	go run cmd/service/main.go
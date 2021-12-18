.PHONY: proto
proto:
	@protoc -I ./proto pancake.proto --go_out=api/gen --go-grpc_out=api/gen --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative
.PHONY: proto1
proto1:
	@protoc -Iproto image_uploader.proto --go_out=plugins=grpc:api --go_opt=paths=source_relative
.PHONY: prp
prp:
	@protoc -I .  call.proto --go_out=call --go-grpc_out=call --go-grpc_opt=paths=source_relative --go_opt=paths=source_relative --js_out=import_style=commonjs:client/src/proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:client/src/proto

.PHONY: dev
dev:
	docker-compose -f docker/docker-compose.yml up --build
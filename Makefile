.PHONY: api gw zrpc gorm-gen

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
PROTO_FILE="dsl/simple.proto"
API_FILE="dsl/simple.api"

# gw produces grpc-gateway .pb.gw.go file
gw:
	protoc -I dsl --grpc-gateway_out ./pb \
	 --grpc-gateway_opt logtostderr=true \
	 --grpc-gateway_opt paths=source_relative \
	 $(PROTO_FILE)

swagger:
	protoc -I dsl --openapiv2_out . \
	--openapiv2_opt logtostderr=true \
	$(PROTO_FILE)

zapi:
	goctl api go --api=$(API_FILE) --dir=./apitmp &&\
	if ! [ -e ./internal/api ]; then \
	  mkdir ./internal/api; \
	fi &&\
	cp -rn ./apitmp/internal/handler/. ./internal/api/handler 2>/dev/null; \
	cp -rn ./apitmp/internal/logic/. ./internal/api/logic 2>/dev/null; \
	cp ./apitmp/internal/handler/routes.go ./internal/api/handler/routes.go &&\
	cp -r ./apitmp/internal/types/. ./internal/api/types &&\
	find ./internal/api -type f -exec sed -i 's/apitmp\/internal/internal\/api/g' {} + &&\
	find ./internal/api -type f -exec sed -i  's/internal\/api\/svc/internal\/svc/g' {} + &&\
	rm -rf ./apitmp

zrpc:
	goctl rpc protoc $(PROTO_FILE) -I dsl --go_out=. --go-grpc_out=. --zrpc_out=. -m

#gorm-gen:
#	cd $(ROOT_DIR)/rpc/internal && \
#	gentool -dsn 'root:toor@tcp(localhost)/main?charset=utf8mb4&parseTime=true&loc=Local' -tables ...
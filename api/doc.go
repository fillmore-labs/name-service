package api

//go:generate sh -c "protoc --proto_path=../proto --go_out=. --go_opt=paths=import ../proto/fillmore_labs/name_service/*/*.proto"
//go:generate sh -c "protoc --proto_path=../proto --go-grpc_out=. --go-grpc_opt=paths=import ../proto/fillmore_labs/name_service/*/*.proto"

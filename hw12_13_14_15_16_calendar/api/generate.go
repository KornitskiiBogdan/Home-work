package api

//go:generate protoc -I . --go_out=./eventpb --go_opt=paths=source_relative --go-grpc_out=./eventpb --go-grpc_opt=paths=source_relative EventService.proto

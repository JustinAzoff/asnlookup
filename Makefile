all: hostlookup
pb/hostlookup.pb.go: protos/hostlookup.proto
	protoc -I protos hostlookup.proto --go_out=plugins=grpc:pb

hostlookup: pb/hostlookup.pb.go *.go */*.go
	go build

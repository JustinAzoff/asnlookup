asnlookup/asnlookup_pb2.py: protos/asnlookup.proto
	python -m grpc.tools.protoc -Iprotos --python_out=asnlookup --grpc_python_out=asnlookup protos/asnlookup.proto

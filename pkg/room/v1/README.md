# Room gRPC contract

`room.proto` defines the shared backend/room API. The generated Go files are
committed, so ordinary builds do not require protobuf tools.

To regenerate, install protoc (generated with v36.2) and these Go plugins:

```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.12
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.2
```

Put protoc and the Go plugins on `PATH`, then run `make proto` from the repository
root. Commit the `.proto` and regenerated `.pb.go` files together.

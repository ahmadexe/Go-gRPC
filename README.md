# gRPC
gRPC (gRPC Remote Procedure Calls) is a high-performance, open-source framework developed by Google for building distributed systems and microservices. It uses HTTP/2 for transport, Protocol Buffers (protobufs) as the interface description language, and provides features such as authentication, load balancing, and bidirectional streaming. gRPC allows clients and servers to communicate transparently, making it easier to create scalable and efficient applications across various platforms and languages.

## Generating gRPC files
1. Install protoc
2. Go to root of the project
3. Run
  ```
   protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    grpc/user.proto
   ```
4. To run the project go to users service, run ```go run main.go```
5. Go to orders service, run ```go run main.go```
   

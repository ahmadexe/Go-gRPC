# gRPC
It is an open-source remote procedure call framework, developed by Google. So as you might have guessed gRPC stands for google remote procedure call.

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
   

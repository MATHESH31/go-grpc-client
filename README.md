# Go gRPC Client

This repository contains a small Go gRPC client for an `Employee` service. The client connects to a gRPC server running on `localhost:9090`, sends an employee ID, and prints the employee details returned by the server.

## Project Structure

- `main.go`: The application entry point. It creates the gRPC connection, builds the client stub, sends the request, and logs the response.
- `proto/Employee.proto`: The service contract. It defines the `Employee` service, the `getEmployee` RPC, and the request/response message shapes.
- `generated/Employee.pb.go`: Generated protobuf message types for `EmployeeRequest` and `EmployeeResponse`.
- `generated/Employee_grpc.pb.go`: Generated gRPC client and server interfaces and the RPC method metadata.
- `go.mod`: The Go module definition and dependency versions.

## Service Contract

The client is driven by the protobuf definition in `proto/Employee.proto`.

```proto
syntax = "proto3";
package employee.v1;
option go_package = "./generated;employeev1";

service Employee {
  rpc getEmployee (EmployeeRequest) returns (EmployeeResponse) {}
}

message EmployeeRequest {
  int32 id = 1;
}

message EmployeeResponse {
  int32 id = 1;
  string name = 2;
  int32 age = 3;
  double salary = 4;
}
```

Important points:

- `package employee.v1` becomes the gRPC service namespace.
- `option go_package = "./generated;employeev1"` tells the Go generators to place generated code in the `generated` folder and use `employeev1` as the Go package name.
- `EmployeeRequest` carries the employee ID the client wants to look up.
- `EmployeeResponse` contains the server response payload.

## How the Generated Code Works

Two generated files are created from the proto definition:

- `generated/Employee.pb.go` contains the Go structs for the protobuf messages.
- `generated/Employee_grpc.pb.go` contains the gRPC client interface and RPC invocation code.

The generated gRPC client interface looks conceptually like this:

```go
type EmployeeClient interface {
    GetEmployee(ctx context.Context, in *EmployeeRequest, opts ...grpc.CallOption) (*EmployeeResponse, error)
}
```

That generated `GetEmployee` method is what `main.go` calls. Even though the proto method name is written as `getEmployee`, the generated Go method follows Go naming conventions and becomes `GetEmployee`.

## Client Implementation

The client flow in `main.go` is straightforward:

1. Create a gRPC connection to `localhost:9090`.
2. Use insecure transport credentials for local development.
3. Build the generated `EmployeeClient`.
4. Create a request-scoped context with a 5-second timeout.
5. Send an `EmployeeRequest` with `Id: 1`.
6. Receive the `EmployeeResponse`.
7. Log the response fields.

Here is the implementation pattern used by the repo:

```go
conn, err := grpc.NewClient(
    "localhost:9090",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
if err != nil {
    log.Fatalf("Failed to connect: %v", err)
}
defer conn.Close()

client := pb.NewEmployeeClient(conn)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

request := &pb.EmployeeRequest{Id: 1}

response, err := client.GetEmployee(ctx, request)
if err != nil {
    log.Fatalf("Failed to make RPC: %v", err)
}
```

## Runtime Request Flow

When the client runs, the request flows like this:

1. `main.go` creates a protobuf request object.
2. `client.GetEmployee(...)` calls into generated code in `generated/Employee_grpc.pb.go`.
3. The generated stub serializes the request and invokes the gRPC method `/employee.v1.Employee/getEmployee`.
4. The remote server handles the request and sends back a protobuf response.
5. The generated code deserializes the response into `EmployeeResponse`.
6. `main.go` logs the fields returned by the server.

## Generate the gRPC Files

If you update `proto/Employee.proto`, regenerate the Go files from the project root with:

```bash
PATH="$HOME/go/bin:$PATH" protoc \
  --proto_path=proto \
  --go_out=. \
  --go-grpc_out=. \
  proto/Employee.proto
```

This command produces:

- `generated/Employee.pb.go`
- `generated/Employee_grpc.pb.go`

If the plugins are not installed, add them with:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Run the Client

Make sure a compatible gRPC server is already running on `localhost:9090`, then start the client:

```bash
go run .
```

Expected behavior:

- The client connects to the server.
- It requests employee data for `Id: 1`.
- It logs `Id`, `Name`, `Age`, and `Salary`.

## Dependencies

The main runtime dependencies are:

- `google.golang.org/grpc`: gRPC client support.
- `google.golang.org/protobuf`: protobuf runtime and code generation support.

These are declared in `go.mod` and locked in `go.sum`.

## Current Scope

This repository contains only the client and generated contract files. It does not include the gRPC server implementation, so the client depends on an external server that implements the same `Employee` service contract.

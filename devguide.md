## Tools:
- Kreya: To test the gRPC services and generate client code.

anything else is typical grpc stuff.

Use this command in the root of the grpc project to generate the gRPC code from the proto files:
```bash
protoc --proto_path=. --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative proto/*/*.proto
```

## Proto Files:
All proto files are located in the `proto` directory. Each service has its own subdirectory, and the proto file for the service is named `service.proto`. For example, the proto file for the DocumentService is located at `proto/document/service.proto`.

## Services:
All services are defined in the `services` directory. Each service has its own subdirectory, and the implementation of the service is in a file named `service.go`. For example, the DocumentService is implemented in `services/document/document.go`.

## Local development:
There are still several things to setup for local development, such as mtls certificates and certificate management as well as a guide on how to set it up. It will most likely end in me scripting this whole step to make it easier for everyone to get started with local development. For now, tough it out.

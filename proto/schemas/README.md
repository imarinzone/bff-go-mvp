# Proto Schemas

This directory contains all protobuf schema definitions that can be shared with downstream services.

## Structure

```
schemas/
├── common/
│   └── context.proto      # Common context messages (beckn protocol)
└── discover/
    └── discover.proto     # Discover service definitions
```

## Usage for Downstream Services

To use these proto files in your downstream service:

1. **Copy the `schemas` directory** to your downstream service repository, or reference it as a git submodule/subtree.

2. **Generate protobuf code** using protoc:

```bash
# Set proto_path to the schemas directory
protoc --proto_path=schemas \
  --go_out=./gen \
  --go_opt=paths=source_relative \
  --go-grpc_out=./gen \
  --go-grpc_opt=paths=source_relative \
  schemas/common/context.proto \
  schemas/discover/discover.proto
```

3. **Import in your proto files**:

```protobuf
import "common/context.proto";

message YourRequest {
  common.Context context = 1;
  // ... your fields
}
```

## Proto Path

When using these schemas, always set `--proto_path=schemas` (or the path where you've placed the schemas directory) so that imports resolve correctly.

## Example: Using in a Go Service

```go
import (
    "bff-go-mvp/proto/common/gen" // or your import path
    "bff-go-mvp/proto/discover/gen"
)

// Use the generated types
req := &discover.DiscoverRequest{
    Context: &common.Context{
        Version: "1.0.0",
        Action:  "discover",
        // ...
    },
    // ...
}
```


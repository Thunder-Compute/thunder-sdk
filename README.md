# Thunder SDK

Thunder SDK is a small Go client library for Thunder Compute integrations. It
wraps the Thunder Central organization API-token surface used by automation that
creates enrollment tokens, manages nodes and clients, and works with zones.

## Install

```sh
go get github.com/Thunder-Compute/thunder-sdk
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	thunder "github.com/Thunder-Compute/thunder-sdk"
)

func main() {
	ctx := context.Background()
	client := thunder.NewClient("", os.Getenv("THUNDER_API_TOKEN"))

	zones, err := client.ListZones(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("zones: %+v\n", zones)
}
```

Pass an empty base URL to use the default Thunder Central API endpoint:

```go
client := thunder.NewClient("", apiToken)
```

Use a custom endpoint, HTTP client, user agent, or installer URL when needed:

```go
client := thunder.NewClient(
	"https://api.thundercompute.com:2096",
	apiToken,
	thunder.WithUserAgent("my-integration/1.0"),
)
```

## Common Operations

Create a client enrollment token:

```go
token, err := client.CreateClientEnrollment(ctx, thunder.CreateClientEnrollmentRequest{
	ZoneID:   "zone-1",
	GPUType:  "nvidia-l4",
	GPUCount: 1,
})
```

Create a node enrollment token:

```go
token, err := client.CreateNodeEnrollment(ctx, thunder.CreateNodeEnrollmentRequest{
	ZoneID: "zone-1",
})
```

List registered nodes and clients:

```go
nodes, err := client.ListNodes(ctx, "zone-1")
clients, err := client.ListClients(ctx, "zone-1")
```

Create and delete zones:

```go
zone, err := client.CreateZone(ctx, thunder.CreateZoneRequest{
	DisplayName: "production",
})

err = client.DeleteZone(ctx, zone.ZoneID)
```

Generate install commands for enrollment:

```go
cmd := client.NodeEnrollmentCommand(thunder.NodeEnrollmentCommandRequest{
	EnrollmentToken: token.EnrollmentToken,
	NodeName:        "worker-1",
})
fmt.Println(cmd)
```

## Errors

Non-2xx API responses are returned as `*thunder.APIError`. Helper functions are
available for common status checks:

```go
if thunder.IsForbidden(err) {
	// The API token is valid but lacks the required capability.
}
```

## Development

Run the test suite with:

```sh
go test ./...
```

This module has no third-party Go dependencies.

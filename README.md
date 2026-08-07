# Thunder SDK

Thunder SDK is a small Go client library for Thunder Compute integrations. It
wraps the Thunder Central organization API-token surface used by automation that
creates enrollment tokens, manages servers and clients, and works with zones.

## Audience

This SDK is intended for **Thunder Enterprise** users building automation against
Thunder Compute. It is not intended for Thunder
Cloud users; cloud users should use the [Thunder CLI](https://github.com/Thunder-Compute/thunder-cli)
and [Thunder Compute Documentation](https://www.thundercompute.com/docs). Those resources
detail how to interface with instances and snapshots for our cloud offering.

## Install

Requires Go 1.22 or newer.

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

The client sends the token as a bearer token in the `Authorization` header.

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

Create a server enrollment token:

```go
token, err := client.CreateServerEnrollment(ctx, thunder.CreateServerEnrollmentRequest{
	ZoneID: "zone-1",
})
```

List registered servers and clients:

```go
servers, err := client.ListServers(ctx, "zone-1")
clients, err := client.ListClients(ctx, "zone-1")
```

Create and delete zones:

```go
zone, err := client.CreateZone(ctx, thunder.CreateZoneRequest{
	DisplayName: "production",
})

err = client.DeleteZone(ctx, zone.ZoneID)
```

Read and replace zone GPU oversubscription targets:

```go
targets, err := client.ListZoneOversubscriptionTargets(ctx, "zone-1")

updated, err := client.SetZoneOversubscriptionTargets(ctx, "zone-1", []thunder.ZoneOversubscriptionTarget{
	{GPUType: "nvidia-l4", OversubscriptionTarget: 2.5},
})
```

Generate install commands for enrollment:

```go
cmd := client.ServerEnrollmentCommand(thunder.ServerEnrollmentCommandRequest{
	EnrollmentToken: token.EnrollmentToken,
	ServerName:      "worker-1",
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

## API Token Capabilities

Thunder Central API tokens need the capabilities required by the operation being
called. The SDK exports capability constants such as
`CapabilityCreateClientEnrollmentToken`, `CapabilityReadZones`, and
`CapabilityWriteZones` so integrations can keep their setup code and
documentation aligned with the API surface they use.

## Development

Run the test suite with:

```sh
go test ./...
```

If you are working inside Thunder's Bazel workspace, the equivalent target is:

```sh
bazel test //thunder-sdk:thunder_test
```

This module has no third-party Go dependencies.

## Versioning

Public releases should be tagged with semantic versions, for example `v0.1.0`.
Consumers can then pin a release with:

```sh
go get github.com/Thunder-Compute/thunder-sdk@v0.1.0
```

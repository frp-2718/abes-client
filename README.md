# abes-client

Go client for the [ABES / SUDOC web services](https://abes.fr/api-et-web-services/).

[![Latest tag](https://img.shields.io/github/v/tag/frp-2718/abes-client)](https://github.com/frp-2718/abes-client/tags)

## Wrapped services

| Service | Endpoint | Methods |
|---|---|---|
| isbn2ppn | `https://www.sudoc.fr/services/isbn2ppn/` | `Isbn2ppn.GetPPNs`, `Isbn2ppn.GetMultiPPNs` |
| multiwhere | `https://www.sudoc.fr/services/multiwhere/` | `Multiwhere.GetLocations`, `Multiwhere.GetMultiLocations`, `Multiwhere.GetMultiLocationsWithErrors` |
| unimarc2marcxml | `https://www.sudoc.fr/<ppn>.xml` | `UnimarcXML.GetRecord` (a MARCXML parser is also provided in `marc.go`) |

Official ABES documentation: [webservices manual](https://documentation.abes.fr/sudoc/manuels/administration/aidewebservices/index.html).

## Installation

```bash
go get github.com/frp-2718/abes-client/abes
```

Requires Go 1.22 or later.

## Quick start

```go
package main

import (
    "fmt"
    "log"

    "github.com/frp-2718/abes-client/abes"
)

func main() {
    // nil = http.DefaultClient. Pass your own *http.Client (with Timeout,
    // custom transport, proxy...) to control the HTTP layer.
    client := abes.NewAbesClient(nil)

    // ISBN -> PPNs (with and without holdings)
    res, err := client.Isbn2ppn.GetPPNs("978-0134190440")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(res)

    // PPN -> holding libraries
    libs, err := client.Multiwhere.GetLocations("154923206")
    if err != nil {
        log.Fatal(err)
    }
    for _, l := range libs {
        fmt.Println(l)
    }
}
```

## Error handling

All services return typed errors, inspectable with `errors.As`:

| Error type | Meaning | Fields |
|---|---|---|
| `*abes.NotFoundError` | the resource does not exist | `Op`, `URL`, `Message` (server's `<error>` text) |
| `*abes.APIError` | any other non-2xx HTTP status | `Op`, `URL`, `StatusCode`, `Body` |
| `*abes.NetworkError` | transport failure (DNS, TLS, ...) | `Op`, `URL`, `Err` (via `Unwrap`) |
| `*abes.SystemError` | unreadable body / invalid XML | `Op`, `URL`, `Err` (via `Unwrap`) |

ABES services report "not found" inconsistently: `isbn2ppn` and `<ppn>.xml`
answer HTTP 404, while `multiwhere` answers HTTP 500 with an `<error>`
element. The client normalizes this: use `abes.NotFound(err)` for the common
case.

```go
libs, err := client.Multiwhere.GetLocations(ppn)
switch {
case err == nil:
    // use libs
case abes.NotFound(err):
    // unknown PPN: ignore it or mark the record as missing
default:
    var apiErr *abes.APIError
    if errors.As(err, &apiErr) {
        log.Printf("SUDOC answered HTTP %d on %s", apiErr.StatusCode, apiErr.URL)
    }
    log.Fatal(err)
}
```

Note on batch calls (`GetMultiPPNs`, `GetMultiLocations`): identifiers not
found by the server are simply absent from the returned map (the server
omits them from its 200 response). Use `GetMultiLocationsWithErrors` to also
get the list of invalid PPNs. A single unknown identifier returns a
`*NotFoundError`.

## Examples

Runnable programs live in [examples/](examples/):

```bash
go run ./examples/isbn
go run ./examples/locations
go run ./examples/multiLocations
go run ./examples/checkMSC
```

<!-- ## API documentation -->

<!-- Full reference on [pkg.go.dev](https://pkg.go.dev/github.com/frp-2718/abes-client/abes), -->
<!-- generated from the doc comments; the runnable examples from -->
<!-- `abes/example_test.go` are displayed there too. -->

## Development

```bash
go test ./...
```

Tests run against local `httptest` servers: no network access needed.

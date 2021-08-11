# StreamingFast Data Integrity Proof Protocol middleware in Go
[![reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/streamingfast/dipp)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Simple protocol to provide your users with irrefutable proofs of the
output of your service, as pioneered by dfuse.io

### Using the protocol

Have your users pass the:

```
X-Data-Integrity-Proof: true
```

header in requests to your service.  The protocol then returns an `X-Data-Integrity-Proof` hash back in the response headers.

If any data integrity issues arise, bring the returned payload along with the proof to claim your bounty. This will provide us with irrefutable proof that we did serve the payload.

Note: this will not work on websocket streams.


### Integrate the middleware

In a Go service, use:

```go
router := http.NewServeMux()
// Add routes to router

secret := "this is a randomly generated secret that is at least 32 bytes long"

_ = http.ListenAndServe(":8080", dipp.NewProofMiddleware(secret, router))
```


### Validation

When you receive report of faulty data, users can provide the proof
you gave them, with the payload they received.

Use the provided `dipp-checker` tool to do validation, which you can install with:

```
go get github.com/streamingfast/dipp/cmd/dipp-checker
```

and run as:

```
dipp-checker \
    "this is a randomly generated secret that is at least 32 bytes long" \
    "the-proof-provided-by-the-user-in-hexadecimal" \
    ./path/to/payload.file
```

to validate integrity.


## Initial work

The dfuse.io platform started this initiative alongside its Data Integrity Bounty Program.

See https://hackerone.com/dfuse for details.


## Contributing

**Issues and PR in this repo related strictly to the dipp library.**

Report any protocol-specific issues in their
[respective repositories](https://github.com/streamingfast/streamingfast#protocols)

**Please first refer to the general
[StreamingFast contribution guide](https://github.com/streamingfast/streamingfast/blob/master/CONTRIBUTING.md)**,
if you wish to contribute to this code base.

This codebase uses unit tests extensively, please write and run tests.

## License

[Apache 2.0](LICENSE)

Data Integrity Proof Protocol middleware in Go
==============================================

Simple protocol to provide your users with irrefutable proofs of the
output of your service.


### Integrate middleware

Example:

```go
router := http.NewServeMux()
// Add routes to router

secret := "this is a randomly generated secret that is at least 32 bytes long"

_ = http.ListenAndServe(":8080", dipp.NewProofMiddleware(secret, router))
```

Then ask your users to add this header to their request:

```
X-Data-Integrity-Proof: true
```

and they will receive an `X-Data-Integrity-Proof` hash back in the
response headers.

Note: this will not work on websocket streams.

### Validation

When you receive report of faulty data, users can provide the proof
you gave them, with the payload they received.

Use the provided `dipp-checker` tool to do validation, which you can install with:

```
go get github.com/eoscanada/dipp/cmd/dipp-checker
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



## LICENSE

MIT

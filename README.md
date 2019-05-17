Data Integrity Proof Protocol middleware in Go
==============================================

Simple protocol to give proofs to your users of the data you provide them.

Integrate the middleware:

```go
    dipp.ProofMiddleware("this is a randomly generated secret that is at least 32 bytes long", http.Handler(next))
```

Then ask your users to add this header to their request:

```
X-Data-Integrity-Proof: true
```

and they will receive an `X-Data-Integrity-Proof` hash back in the
response headers.

Note: this will not work on websocket streams.


Initiator
---------

The dfuse.io platform started this initiative alongside its Data Integrity Bounty Program.

See https://hackerone.com/dfuse-io for details.



LICENSE
-------

MIT

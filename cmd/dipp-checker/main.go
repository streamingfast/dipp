package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dfuse-io/dipp"
)

func main() {
	if len(os.Args) != 4 {
		errorCheck("Usage: dipp-checker [secret-as-string] [proof-hash-in-hex] [filename with contents]\n", fmt.Errorf("  -> error: found %d arguments", len(os.Args)))
	}

	secret := os.Args[1]
	proofHash := os.Args[2]
	filename := os.Args[3]

	cnt, err := ioutil.ReadFile(filename)
	errorCheck("opening file", err)

	computedHash := dipp.HashMac(secret, cnt)

	if proofHash == computedHash {
		fmt.Printf("Valid. Hash %q matches.\n", computedHash)
	} else {
		fmt.Printf("INVALID. Provided proof %q does not match computed hash %q.\n", proofHash, computedHash)
		os.Exit(1)
	}
}

func errorCheck(prefix string, err error) {
	if err != nil {
		fmt.Println(prefix, err)
		os.Exit(1)
	}
}

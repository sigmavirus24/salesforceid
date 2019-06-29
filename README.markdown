# salesforceid

[![GoDoc](https://godoc.org/github.com/sigmavirus24/salesforceid?status.svg)](https://godoc.org/github.com/sigmavirus24/salesforceid)
[![CircleCI](https://circleci.com/gh/sigmavirus24/salesforceid.svg?style=svg)](https://circleci.com/gh/sigmavirus24/salesforceid)
[![GolangCI](https://golangci.com/badges/github.com/sigmavirus24/salesforceid.svg)](https://golangci.com)

This is a library that enables handling and manipulation of Salesforce 
Identifiers. Specifically it enables converting identifiers between 15 
character (case-sensitive) and 18 character (case-insensitive) versions.

This library lazily parses a Salesforce identifier for the user and provides 
other convenience methods for handling portions of the identifier.

## Examples

It's possible to have different versions of a Salesforce identifier. 15 
character identifiers are case-sensitive while 18 character identifers are 
not. As a result, this library always checks and adjusts the casing of 
identifiers and always produces 18 character identifiers. This allows users to 
confidently use which ever form they prefer.

```go
package main

import (
	"fmt"
	"os"

	"github.com/sigmavirus24/salesforceid"
)

func main() {
	// Note that we're using an 18-character, entirely lower-cased Salesforce identifier here
	sfid, err := salesforceid.New("00d000000000062eaa")
	if err != nil {
		fmt.Printf("encountered unexpected error: %q", err)
		os.Exit(1)
	}
	fmt.Println(sfid)
}
```

Furthermore, one can use this library to perform arithmetic on the 
identifiers. For an example of where this might be useful see the 
[example](./example_test.go) in this project.

# salesforceid

[![GoDoc Reference](https://pkg.go.dev/badge/github.com/sigmavirus24/salesforceid?utm_source=godoc)](https://pkg.go.dev/github.com/sigmavirus24/salesforceid)

This is a library that enables handling and manipulation of Salesforce
Identifiers. Specifically it enables converting identifiers between 15
character (case-sensitive) and 18 character (case-insensitive) versions.

This library parses a Salesforce identifier for the user and provides
some convenience methods for handling interactions with the identifier.

For explanations of how the identifier works see:

* ["Obscure Salesforce Object Key Prefixes"](https://www.fishofprey.com/2011/09/obscure-salesforce-object-key-prefixes.html) [Archive.org](https://web.archive.org/web/20251219015022/https://www.fishofprey.com/2011/09/obscure-salesforce-object-key-prefixes.html)
* ["Salesforce IDs Explained"](https://codebycody.com/salesforce-ids-explained/) [Archive.org](https://web.archive.org/web/20251111164629/https://codebycody.com/salesforce-ids-explained/)
* ["Salesforce Object ID Is Refined to Use Three Characters for Server IDs
  (Release
  Update)"](https://help.salesforce.com/s/articleView?language=en_US&id=release-notes.rn_hyperforce_object_id.htm&release=246&type=5)

## Changelog

### v1.0.0 - 2026-02-13

* Fix parsing for `SalesforceID` in `New` to correct number of bytes
  `Reserved` and number of bytes used in the `NumericIdentifier` fields.

  Previously, we allocated one of the reserved bytes to the numeric identifier
  incorrectly.

* Add `SalesfoceIDV2` and `NewV2` for Salesforce Object Identifiers. These
  reflect the changes to Salesforce Object IDs to use 3 bytes (using the 6th
  byte that was previously reserved) to represent Instances (Servers).

  This also renames `PodIdentifier` from `SalesforceID` to
  `InstanceIdentifier` on `SalesforceIDV2` reflect the new
  language used.

## Examples

It's possible to have different versions of a Salesforce identifier. 15
character identifiers are case-sensitive while 18 character identifiers are
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

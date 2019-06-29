package salesforceid_test

import (
	"fmt"

	sfid "github.com/sigmavirus24/salesforceid"
)

func GenerateChunks(starting *sfid.SalesforceID, chunkSize uint64) func() (*sfid.SalesforceID, *sfid.SalesforceID) {
	next := starting
	return func() (*sfid.SalesforceID, *sfid.SalesforceID) {
		current := next
		end, _ := current.Add(chunkSize)
		next, _ = end.Add(1)
		return current, end
	}
}

func Example() {
	s, _ := sfid.New("001000000000000")
	chunker := GenerateChunks(s, 250000)

	for i := 0; i < 15; i++ {
		start, end := chunker()
		fmt.Printf("SELECT Id FROM Account WHERE Id >= %s AND Id <= %s\n", start, end)
	}
	// Output:
	// SELECT Id FROM Account WHERE Id >= 001000000000000AAA AND Id <= 00100000000132GAAQ
	// SELECT Id FROM Account WHERE Id >= 00100000000132HAAQ AND Id <= 00100000000264XAAQ
	// SELECT Id FROM Account WHERE Id >= 00100000000264YAAQ AND Id <= 00100000000396oAAA
	// SELECT Id FROM Account WHERE Id >= 00100000000396pAAA AND Id <= 001000000004C95AAE
	// SELECT Id FROM Account WHERE Id >= 001000000004C96AAE AND Id <= 001000000005FBMAA2
	// SELECT Id FROM Account WHERE Id >= 001000000005FBNAA2 AND Id <= 001000000006IDdAAM
	// SELECT Id FROM Account WHERE Id >= 001000000006IDeAAM AND Id <= 001000000007LFuAAM
	// SELECT Id FROM Account WHERE Id >= 001000000007LFvAAM AND Id <= 001000000008OIBAA2
	// SELECT Id FROM Account WHERE Id >= 001000000008OICAA2 AND Id <= 001000000009RKSAA2
	// SELECT Id FROM Account WHERE Id >= 001000000009RKTAA2 AND Id <= 00100000000AUMjAAO
	// SELECT Id FROM Account WHERE Id >= 00100000000AUMkAAO AND Id <= 00100000000BXP0AAO
	// SELECT Id FROM Account WHERE Id >= 00100000000BXP1AAO AND Id <= 00100000000CaRHAA0
	// SELECT Id FROM Account WHERE Id >= 00100000000CaRIAA0 AND Id <= 00100000000DdTYAA0
	// SELECT Id FROM Account WHERE Id >= 00100000000DdTZAA0 AND Id <= 00100000000EgVpAAK
	// SELECT Id FROM Account WHERE Id >= 00100000000EgVqAAK AND Id <= 00100000000FjY6AAK
}

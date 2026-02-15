package salesforceid

import (
	"errors"
	"fmt"
)

var ErrInvalidEdition = errors.New("invalid edition provided")

// ErrInvalidSFID is returned when an 18 character identifier is provided
// whose check bytes do not make sense with the 15 character identifier
var ErrInvalidSFID = errors.New("check bytes do not match identifier")

// ErrInvalidLengthSFID indicates when an SFID is not a known good length
var ErrInvalidLengthSFID = errors.New("sfids should be 15 or 18 characters")

// ErrValueTooLarge is returned when attempting to encode a Base62 value
// larger than will fit in a NumericIdentifier
var ErrValueTooLarge = fmt.Errorf("value is larger than %d", MaxIdentifierValue)

// ErrInvalidNumericIdentifier is returned when attempting to decode a numeric
// identifier more than 9 bytes long or has invalid bytes
var ErrInvalidNumericIdentifier = errors.New("numeric identifier being decoded is invalid")

// ErrInvalidAddition is returned when the NumericIdentifier is already the
// largest allowed in a 9 digit base62 encoded number
var ErrInvalidAddition = fmt.Errorf("addition would overflow maximum value: %d", MaxIdentifierValue)

// ErrInvalidSubtraction is returned when the amount to subtract from the
// identifier is greater than the decoded value of the NumericIdentifier
var ErrInvalidSubtraction = errors.New("subtraction would result in a negative identifier")

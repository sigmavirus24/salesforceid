package salesforceid

import (
	"bytes"
	"errors"
)

// SalesforceID stores and manages the Salesforce Identifier and its
// components. The full identifier is not an accessible attribute. To retrieve
// it, use the String method.
type SalesforceID struct {
	id                []byte
	KeyPrefix         []byte // KeyPrefix consists of the first 3 bytes of an id. It is used to identify the object
	PodIdentifier     []byte // PodIdentifier consists of the fourth and fifth bytes. It maps to the pod on which the record was created
	Reserved          []byte // Reserved is a single byte reserved for future use. It should always be '0'
	NumericIdentifier []byte // NumericIdentifier is a Base 62 encoded number that auto-increments for each record. You can decode it with Decode
	Suffix            []byte // Suffix are the three bytes at the end of an 18 character identifier. These help determine the casing of a 15 character identifier
}

// While Salesforce IDs use A-Z, a-z, and 0-9 we use this slice to find a
// value from 0-31.
var sfidSeq = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ012345")

// ErrInvalidSFID is returned when an 18 character identifier is provided
// whose check bytes do not make sense with the 15 character identifier
var ErrInvalidSFID = errors.New("check bytes do not match identifier")

// ErrInvalidLengthSFID indicates when an SFID is not a known good length
var ErrInvalidLengthSFID = errors.New("sfids should be 15 or 18 characters")

// ErrValueTooLarge is returned when attempting to encode a Base62 value
// larger than will fit in a NumericIdentifier
var ErrValueTooLarge = errors.New("value is larger than 62^9")

// ErrInvalidNumericIdentifier is returned when attempting to decode a numeric
// identifier more than 9 bytes long or has invalid bytes
var ErrInvalidNumericIdentifier = errors.New("numeric identifier being decoded is invalid")

// ErrInvalidAddition is returned when the NumericIdentifier is already the
// largest allowed in a 9 digit base62 encoded number
var ErrInvalidAddition = errors.New("addition would overflow maximum value: 13537086546263552")

// ErrInvalidSubtraction is returned when the amount to subtract from the
// identifier is greater than the decoded value of the NumericIdentifier
var ErrInvalidSubtraction = errors.New("subtraction would result in a negative identifier")

// New generates a SalesforceID for usage
func New(id string) (*SalesforceID, error) {
	var err error
	if len(id) != 15 && len(id) != 18 {
		return nil, ErrInvalidLengthSFID
	}
	idBytes := []byte(id)
	if len(idBytes) == 18 {
		idBytes, err = normalize(idBytes)
		if err != nil {
			return nil, err
		}
	}
	if len(id) == 15 {
		idBytes = computeEighteen(idBytes)
	}
	return &SalesforceID{
		idBytes,
		idBytes[0:3],
		idBytes[3:5],
		idBytes[5:6],
		idBytes[6:15],
		idBytes[15:18],
	}, nil
}

func computeEighteen(id []byte) []byte {
	var chunkSum uint
	newSFID := make([]byte, 15)
	copy(newSFID, id[0:15])

	for i, b := range id {
		if i%5 == 0 {
			if i != 0 {
				newSFID = append(newSFID, sfidSeq[chunkSum])
			}
			chunkSum = 0
		}
		if 'A' <= b && b <= 'Z' {
			chunkSum += 1 << uint(i%5)
		}
	}
	return append(newSFID, sfidSeq[chunkSum])
}

func normalize(id []byte) ([]byte, error) {
	check := bytes.ToUpper([]byte(id[15:]))
	copy(id[15:18], check)
	for i, b := range id[:15] {
		checkByte := check[i/5]
		pow := 1 << uint(i%5)
		checkVal := bytes.IndexByte(sfidSeq, checkByte)
		// In each chunk, our left most byte is the least significant
		// digit in the calculation that creates the last three
		// digits. Those three digits are based on the case of the
		// first 15 bytes. If the chunk looks like 0b11011 that means
		// the 1st, 2nd, 4th, and 5th bytes in this chunk should be
		// upper case. So we calculate 1 << j which translates to one
		// of:
		// - 0b00001
		// - 0b00010
		// - 0b00100
		// - 0b01000
		// - 0b10000
		// And use bitwise and to determine this byte's case
		if checkVal&pow == pow {
			if '0' <= b && b <= '9' {
				return nil, ErrInvalidSFID
			}
			if 'a' <= b && b <= 'z' {
				// If this is lower case but should be upper, handle that
				id[i] = b - 32
			}
		} else if checkVal&pow == 0 && 'A' <= b && b <= 'Z' {
			// If this is upper case but should be lower, handle that
			id[i] = b + 32
		}
	}
	return id, nil
}

func (s *SalesforceID) String() string {
	return string(s.id)
}

// Add a value to the numeric identifier and ensure the resulting identifier
// is valid.
func (s *SalesforceID) Add(i uint64) (*SalesforceID, error) {
	newID := make([]byte, 15)
	copy(newID, s.id)
	decoded, err := Decode(s.NumericIdentifier)
	if err != nil {
		return nil, err
	}
	if decoded == maxID-1 {
		return nil, ErrInvalidAddition
	}
	encoded, err := Encode(decoded + i)
	if err != nil {
		return nil, err
	}
	copy(newID[6:15], encoded)
	return New(string(newID))
}

// Subtract a value from the numeric identifier and ensure the resulting
// identifier is valid.
func (s *SalesforceID) Subtract(i uint64) (*SalesforceID, error) {
	decoded, err := Decode(s.NumericIdentifier)
	if err != nil {
		return nil, err
	}
	if decoded < i {
		return nil, ErrInvalidSubtraction
	}
	newID := make([]byte, 15)
	copy(newID, s.id)
	encoded, err := Encode(decoded - i)
	if err != nil {
		return nil, err
	}
	copy(newID[6:15], encoded)
	return New(string(newID))
}

var table = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

const (
	zeroString = "000000000"
	maxID      = 13537086546263552 // == 62 ^ 9
	base       = 62
)

// Encode converts an unsigned integer to Base62 for use as a
// NumericIdentifier. This can only handle 62^9.
func Encode(u uint64) (string, error) {
	if u == 0 {
		return zeroString, nil
	}
	if u >= maxID {
		return "", ErrValueTooLarge
	}
	ebytes := make([]byte, 0, 9)
	for p := uint64(maxID / base); p >= 1; p /= base {
		m := u / p
		ebytes = append(ebytes, table[m])
		u -= m * p
	}
	return string(ebytes), nil
}

// Decode converts bytes to an unsigned integer
func Decode(src []byte) (uint64, error) {
	var v uint64
	var err error

	if len(src) != 9 {
		return v, ErrInvalidNumericIdentifier
	}

	m := uint64(1)
	for i := 8; i >= 0; i-- {
		c := src[i]
		if c >= '0' && c <= '9' {
			v += uint64((c - '0')) * m
		} else if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			v += uint64(bytes.IndexByte(table, c)) * m
		} else {
			v = 0
			err = ErrInvalidNumericIdentifier
			break
		}
		m *= 62
	}

	return v, err
}

package salesforceid

import (
	"bytes"
	"errors"
)

// SFID provides certain methods that
type SFID interface {
	To15() (SFID, error)
	To18() (SFID, error)
	String() string
	KeyPrefix() string
	PodIdentifier() string
	Reserved() string
	NumericIdentifier() string
	Suffix() (string, error)
}

// SalesforceID stores and manages the Salesforce Identifier
type SalesforceID struct {
	id string
}

// While Salesforce IDs use A-Z, a-z, and 0-9 we use this slice to find a
// value from 0-31.
var sfidSeq = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ012345")

// ErrInvalidSFID indicates when an SFID does not meet the necessary
// requirements of a valid SFID.
var ErrInvalidSFID = errors.New("sfid provided is invalid")

// ErrInvalidLengthSFID indicates when an SFID is not a known good length
var ErrInvalidLengthSFID = errors.New("sfids should be 15 or 18 characters")

// ErrNoSuffix is returned when requesting a suffix for an SFID that is fewer
// than 18 characters.
var ErrNoSuffix = errors.New("to retrieve a suffix, the sfid must be 18 characters")

// New generates a SalesforceID for usage
func New(id string) SalesforceID {
	return SalesforceID{id}
}

func (s SalesforceID) String() string {
	return s.id
}

// To15 converts the given SalesforceID to a 15-character, case-sensitive,
// Salesforce ID.
func (s SalesforceID) To15() (SFID, error) {
	oldSfid := s.id
	if len(oldSfid) != 15 && len(oldSfid) != 18 {
		return nil, ErrInvalidLengthSFID
	}
	if len(oldSfid) == 15 {
		return SalesforceID{oldSfid}, nil
	}
	fifteen, check := []byte(oldSfid[:15]), bytes.ToUpper([]byte(oldSfid[15:]))
	chunks := [][]byte{
		fifteen[:5],
		fifteen[5:10],
		fifteen[10:],
	}
	for i, chunk := range chunks {
		// NOTE(sigmavirus24) We can probably make this more efficient than
		// having a nested loop. If we iterate over `fifteen` we can calculate
		// the index of the check byte (or precalculate the indices) and then
		// flatten this logic.
		checkVal := bytes.IndexByte(sfidSeq, check[i])
		for j, b := range chunk {
			pow := 1 << uint(j)
			// In our chunk, our left most byte is the least significant
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
			if checkVal&pow == pow && 'a' <= b && b <= 'z' {
				// If this is lower case but should be upper, handle that
				fifteen[(i*5)+j] = b - 32
			} else if checkVal&pow == 0 && 'A' <= b && b <= 'Z' {
				// If this is upper case but should be lower, handle that
				fifteen[(i*5)+j] = b + 32
			}
		}
	}
	return SalesforceID{string(fifteen)}, nil
}

// To18 converts a 15 character Salesforce ID to an 18 character Salesforce
// ID. 15 character SFIDs are case-sensitive but 18 character SFIDs are
// case-insensitive.
func (s SalesforceID) To18() (SFID, error) {
	sfid := s.id
	if len(sfid) != 15 && len(sfid) != 18 {
		return nil, ErrInvalidLengthSFID
	}
	if len(sfid) == 18 {
		return SalesforceID{sfid}, nil
	}
	sfidBytes := []byte(sfid)
	chunks := [][]byte{
		sfidBytes[:5],
		sfidBytes[5:10],
		sfidBytes[10:],
	}
	newSFID := []byte(sfid)
	for _, chunk := range chunks {
		chunkSum := sumUppercase(chunk)
		newSFID = append(newSFID, sfidSeq[chunkSum])
	}
	return SalesforceID{string(newSFID)}, nil
}

// Normalize will create an 18-character Salesforce ID that can be treated as
// case-insensitive but uses the proper casing so that a user can handle the
// string version with confidence.
func (s SalesforceID) Normalize() (SFID, error) {
	if len(s.id) == 15 {
		// We can only assume that this wa sa proper 15 character,
		// case-sensitive SFID so converting it to 18 will do the right thing.
		return s.To18()
	}
	if len(s.id) != 15 && len(s.id) != 18 {
		return nil, ErrInvalidLengthSFID
	}
	fifteen, err := s.To15()
	if err != nil {
		return nil, err
	}
	chunks := [][]byte{
		[]byte(fifteen.String()),
		bytes.ToUpper([]byte(s.id[15:])),
	}
	newSFID := bytes.Join(chunks, []byte{})

	return SalesforceID{string(newSFID)}, nil
}

// KeyPrefix finds and returns the key prefix section of this SFID. This
// prefix tends to map to the object but may not be unique within
// organizations or across organizations.
// *Note*: The KeyPrefix _is_ case-sensitive but this library does not check
// prior to returning the prefix.
func (s SalesforceID) KeyPrefix() string {
	return s.id[:3]
}

// PodIdentifier finds and returns the identifier of the pod the record was
// created in.
func (s SalesforceID) PodIdentifier() string {
	return s.id[3:5]
}

// Reserved returns the reserved byte that may be used for future purposes.
func (s SalesforceID) Reserved() string {
	return s.id[5:6]
}

// NumericIdentifier is the numeric record identifier section for the
// particular record in Salesforce.
func (s SalesforceID) NumericIdentifier() string {
	return s.id[6:15]
}

// Suffix may return the three character suffix if this SalesforceID is an
// 18-character SFID. If it's only a 15 character SFID, then it returns an
// error.
func (s SalesforceID) Suffix() (string, error) {
	if len(s.id) == 18 {
		return s.id[15:], nil
	}
	return "", ErrNoSuffix
}

func sumUppercase(bs []byte) int {
	chunkSum := 0

	for i, b := range bs {
		if 'A' <= b && b <= 'Z' {
			chunkSum += 1 << uint(i)
		}
	}
	return chunkSum
}

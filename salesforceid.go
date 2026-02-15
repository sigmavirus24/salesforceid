package salesforceid

import "fmt"

// Format is used to select the identifier format for representing a
// SalesforceID.
type Format uint8

// IdentifierEdition is used to indicate the format for the identifier.
// Prior to the Summer '23 release, identifiers used two characters to identify
// the Salesforce Pod (a.k.a., Instance, Server). Starting with Summer '23 and
// enforced with Winter '24, one of the two reserved characters was used to
// expand the identifier to three characters.
type IdentifierEdition uint8

const (
	FifteenCharacterFormat Format = iota + 1
	EighteenCharacterFormat

	PreSummer23IdentifierEdition IdentifierEdition = iota + 1
	PostSummer23IdentifierEdition
)

// SalesforceID stores and manages the Salesforce Identifier and its
// components. The full identifier is not an accessible attribute. To retrieve
// it, use the [SalesforceID.String] method or [SalesforceID.Format] method.
// See also:
// * https://codebycody.com/salesforce-ids-explained/
// * https://help.salesforce.com/s/articleView?id=release-notes.rn_hyperforce_object_id.htm&release=246&type=5
type SalesforceID struct {
	id                []byte
	KeyPrefix         []byte // KeyPrefix consists of the first 3 bytes of an id. It is used to identify the object
	PodIdentifier     []byte // PodIdentifier consists of the fourth and fifth bytes prior to Summer '23, or fourth, fifth, and sixth bytes starting with Summer '23. It maps to the pod on which the record was created
	Reserved          []byte // Reserved consists of the sixth and seventh bytes reserved for future use prior to Summer '23 or just the seventh byte starting with Summer '23. It should always be either `[]byte{'0', '0'}` or `[]byte{'0'}`
	NumericIdentifier []byte // NumericIdentifier is a Base 62 encoded number that auto-increments for each record. You can decode it with Decode. It may not be negative or lager than [MaxIdentifierValue].
	Suffix            []byte // Suffix are the three bytes at the end of an 18 character identifier. These help determine the casing of a 15 character identifier
	Edition           IdentifierEdition
}

// New generates a SalesforceID for usage. It assumes the edition to be
// [PreSummer23IdentifierEdition] for backwards compatibility.
func New(id string) (*SalesforceID, error) {
	return Parse(id, PreSummer23IdentifierEdition)
}

func Parse(id string, edition IdentifierEdition) (*SalesforceID, error) {
	idBytes, err := prepareID(id)
	if err != nil {
		return nil, err
	}
	switch edition {
	case PreSummer23IdentifierEdition:
		return &SalesforceID{
			id:                idBytes,
			KeyPrefix:         idBytes[0:3],
			PodIdentifier:     idBytes[3:5],
			Reserved:          idBytes[5:7],
			NumericIdentifier: idBytes[7:15],
			Suffix:            idBytes[15:18],
			Edition:           edition,
		}, nil
	case PostSummer23IdentifierEdition:
		return &SalesforceID{
			id:                idBytes,
			KeyPrefix:         idBytes[0:3],
			PodIdentifier:     idBytes[3:6],
			Reserved:          idBytes[6:7],
			NumericIdentifier: idBytes[7:15],
			Suffix:            idBytes[15:18],
			Edition:           edition,
		}, nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrInvalidEdition, edition)
	}
}

func (s *SalesforceID) String() string {
	return s.Format(EighteenCharacterFormat)
}

func (s *SalesforceID) Format(f Format) string {
	switch f {
	case FifteenCharacterFormat:
		return string(s.id[:15])
	default:
		return string(s.id)
	}
}

// Add a value to the numeric identifier and ensure the resulting identifier
// is valid.
func (s *SalesforceID) Add(i uint64) (*SalesforceID, error) {
	encoded, err := addToID(s.NumericIdentifier, i)
	if err != nil {
		return nil, err
	}
	newID := make([]byte, 15)
	copy(newID, s.id[:7])
	copy(newID[7:15], encoded)
	return New(string(newID))
}

// Subtract a value from the numeric identifier and ensure the resulting
// identifier is valid.
func (s *SalesforceID) Subtract(i uint64) (*SalesforceID, error) {
	encoded, err := subtractFromID(s.NumericIdentifier, i)
	if err != nil {
		return nil, err
	}
	newID := make([]byte, 15)
	copy(newID, s.id[:7])
	copy(newID[7:15], encoded)
	return New(string(newID))
}

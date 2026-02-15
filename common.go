package salesforceid

import "bytes"

var (
	table = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

	// While Salesforce IDs use A-Z, a-z, and 0-9 we use this slice to find a
	// value from 0-31.
	checkSeq = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ012345")
)

const (
	zeroString = "00000000"
	maxString  = "zzzzzzzz"
	// oldMaxID  = 13_537_086_546_263_552 // == 62 ^ 9
	MaxIdentifierValue = 218_340_105_584_896 // == 62 ^ 8
	base               = 62
)

// Encode converts an unsigned integer to Base62 for use as a
// NumericIdentifier. This can only handle 62^8. This returns
// [ErrValueTooLarge] if u is larger than [MaxIdentifierValue].
func Encode(u uint64) (string, error) {
	switch {
	case u == 0:
		return zeroString, nil
	case u == MaxIdentifierValue:
		return maxString, nil
	case u > MaxIdentifierValue:
		return "", ErrValueTooLarge
	}
	ebytes := make([]byte, 0, 8)
	for p := uint64(MaxIdentifierValue / base); p >= 1; p /= base {
		m := u / p
		ebytes = append(ebytes, table[m])
		u -= m * p
	}
	return string(ebytes), nil
}

// Decode converts bytes to an unsigned integer. This returns
// [ErrInvalidNumericIdentifier] if the length of [src] is not 8 or if one of
// the bytes is not a valid Base62 identifier.
func Decode(src []byte) (uint64, error) {
	var v uint64
	var err error

	if len(src) != 8 {
		return v, ErrInvalidNumericIdentifier
	}

	m := uint64(1)
_loop:
	for i := 7; i >= 0; i-- {
		c := src[i]
		switch {
		case c >= '0' && c <= '9':
			v += uint64((c - '0')) * m
		case (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z'):
			v += uint64(bytes.IndexByte(table, c)) * m
		default:
			v = 0
			err = ErrInvalidNumericIdentifier
			break _loop
		}
		m *= 62
	}

	return v, err
}

func computeEighteen(id []byte) []byte {
	var chunkSum uint
	newSFID := make([]byte, 15)
	copy(newSFID, id[0:15])

	for i, b := range id {
		if i%5 == 0 {
			if i != 0 {
				newSFID = append(newSFID, checkSeq[chunkSum])
			}
			chunkSum = 0
		}
		if 'A' <= b && b <= 'Z' {
			chunkSum += 1 << uint(i%5)
		}
	}
	return append(newSFID, checkSeq[chunkSum])
}

func normalize(id []byte) ([]byte, error) {
	check := bytes.ToUpper(id[15:])
	copy(id[15:18], check)
	for i, b := range id[:15] {
		checkByte := check[i/5]
		pow := 1 << uint(i%5)
		checkVal := bytes.IndexByte(checkSeq, checkByte)
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

func prepareID(id string) ([]byte, error) {
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
	return idBytes, nil
}

func addToID(numeric []byte, i uint64) (string, error) {
	decoded, err := Decode(numeric)
	if err != nil {
		return "", err
	}
	newID := decoded + i
	if newID >= MaxIdentifierValue {
		return "", ErrInvalidAddition
	}
	encoded, err := Encode(newID)
	if err != nil {
		return "", err
	}
	return encoded, nil
}

func subtractFromID(numeric []byte, i uint64) (string, error) {
	decoded, err := Decode(numeric)
	if err != nil {
		return "", err
	}
	if decoded < i {
		return "", ErrInvalidSubtraction
	}
	encoded, err := Encode(decoded - i)
	if err != nil {
		return "", err
	}
	return encoded, nil
}

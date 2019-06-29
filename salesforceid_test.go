package salesforceid_test

import (
	"fmt"
	"testing"

	"github.com/sigmavirus24/salesforceid"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		sfid        string
		expected    string
		shouldErr   bool
		expectedErr error
	}{
		{"00D000000000062EAA", "00D000000000062EAA", false, nil},
		{"00D000000000062", "00D000000000062EAA", false, nil},
		{"00d000000000062", "00d000000000062AAA", false, nil},
		{"00d000000000062eaa", "00D000000000062EAA", false, nil},
		{"003D0000001aH2A", "003D0000001aH2AIAU", false, nil},
		{"0a3D0000001aH2A", "0a3D0000001aH2AIAU", false, nil},
		{"0A3D0000001aH2A", "0A3D0000001aH2AKAU", false, nil},
		{"0a3d0000001ah2a", "0a3d0000001ah2aAAA", false, nil},
		{"000000000000000", "000000000000000AAA", false, nil},
		{"999999999999999", "999999999999999AAA", false, nil},
		{"aaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaAAA", false, nil},
		{"zzzzzzzzzzzzzzz", "zzzzzzzzzzzzzzzAAA", false, nil},
		{"AAAAAAAAAAAAAAA", "AAAAAAAAAAAAAAA555", false, nil},
		{"ZZZZZZZZZZZZZZZ", "ZZZZZZZZZZZZZZZ555", false, nil},
		{"ZZZZZZZZZZZZZZ", "", true, salesforceid.ErrInvalidLengthSFID},   // 14 char sfid
		{"ZZZZZZZZZZZZZZZZ", "", true, salesforceid.ErrInvalidLengthSFID}, // 16 char sfid
		{"001000000000062EAA", "", true, salesforceid.ErrInvalidSFID},     // 001 should be 00<cap>
		{"aaaaaaaaaaaaaaa555", "AAAAAAAAAAAAAAA555", false, nil},
		{"ZZZZZZZZZZZZZZZ555", "ZZZZZZZZZZZZZZZ555", false, nil},
		{"zzzzzzzzzzzzzzz555", "ZZZZZZZZZZZZZZZ555", false, nil},
		{"ZzZzZzZzZZzZZzzAAA", "zzzzzzzzzzzzzzzAAA", false, nil},
		{"ZZZZZZZZZZZZZZZAAA", "zzzzzzzzzzzzzzzAAA", false, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.sfid, func(t *testing.T) {
			t.Parallel()
			id, err := salesforceid.New(tc.sfid)
			if tc.shouldErr && err == nil {
				t.Errorf("expected error but didn't get an error")
			}
			if tc.shouldErr && err != tc.expectedErr {
				t.Errorf("expected err %q, got err %q", tc.expectedErr, err)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("expected no error but got %q", err)
			}
			if id != nil {
				s := id.String()
				if tc.expected != s {
					t.Errorf("expected %s, got %s", tc.expected, s)
				}
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	b.Run("15 char sfid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = salesforceid.New("0a3d0000001ah2a")
		}
	})
	b.Run("18 char sfid", func(b *testing.B) {
		b.Run("no changes necessary", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = salesforceid.New("0A3D0000001aH2AKAU")
			}
		})
		b.Run("everything needs correction", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = salesforceid.New("zzzzzzzzzzzzzzz555")
			}
		})
	})
}

func TestAdd(t *testing.T) {
	testCases := []struct {
		sfid        string
		add         uint64
		expected    string
		shouldErr   bool
		expectedErr error
	}{
		{"001000000000000AAA", 1, "001000000000001AAA", false, nil},
		{"001000zzzzzzzzzAAA", 1, "", true, salesforceid.ErrInvalidAddition},
		{"00100000000000zAAA", 1, "001000000000010AAA", false, nil},
		{"00100000000000zAAA", 238328, "00100000000100zAAA", false, nil},
		{"001000000000000AAA", 10, "00100000000000AAAQ", false, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.sfid, func(t *testing.T) {
			want := tc.expected
			sfid, _ := salesforceid.New(tc.sfid)
			got, err := sfid.Add(tc.add)
			if tc.shouldErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if tc.shouldErr && err != tc.expectedErr {
				t.Errorf("expected err %q but got %q", tc.expectedErr, err)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("didn't expect an error but got %q", err)
			}
			if got != nil && want != got.String() {
				t.Errorf("wanted %s, got %s", want, got)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	testCases := []struct {
		sfid        string
		sub         uint64
		expected    string
		shouldErr   bool
		expectedErr error
	}{
		{"001000000000000AAA", 1, "", true, salesforceid.ErrInvalidSubtraction},
		{"001000zzzzzzzzzAAA", 1, "001000zzzzzzzzyAAA", false, nil},
		{"001000000000010AAA", 1, "00100000000000zAAA", false, nil},
		{"00100000000000aAAA", 1, "00100000000000ZAAQ", false, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.sfid, func(t *testing.T) {
			want := tc.expected
			sfid, _ := salesforceid.New(tc.sfid)
			got, err := sfid.Subtract(tc.sub)
			if tc.shouldErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if tc.shouldErr && err != tc.expectedErr {
				t.Errorf("expected err %q but got %q", tc.expectedErr, err)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("didn't expect an error but got %q", err)
			}
			if got != nil && want != got.String() {
				t.Errorf("wanted %s, got %s", want, got)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	testCases := []struct {
		i           uint64
		o           string
		shouldErr   bool
		expectedErr error
	}{
		// Boundaries
		{0, "000000000", false, nil},
		{13537086546263553, "", true, salesforceid.ErrValueTooLarge},
		// Within boundaries
		{5, "000000005", false, nil},
		{15, "00000000F", false, nil},
		{30, "00000000U", false, nil},
		{45, "00000000j", false, nil},
		{61, "00000000z", false, nil},
		{1024, "0000000GW", false, nil},
		{10241024, "00000gy9o", false, nil},
		{6768543273131776, "V00000000", false, nil},
		{13537086546263552 - 1, "zzzzzzzzz", false, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.o, func(t *testing.T) {
			want := tc.o
			got, err := salesforceid.Encode(tc.i)
			if tc.shouldErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if tc.shouldErr && err != tc.expectedErr {
				t.Errorf("expected err %q but got %q", tc.expectedErr, err)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("didn't expect an error but got %q", err)
			}
			if want != got {
				t.Errorf("wanted \"%s\", got \"%s\"", want, got)
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	benchmarks := []struct {
		name string
		v    uint64
	}{
		{"max value", 13537086546263552 - 1},
		{"min value", 0},
		{"middle value", 6768543273131776},
	}
	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = salesforceid.Encode(benchmark.v)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	testCases := []struct {
		in          string
		out         uint64
		shouldErr   bool
		expectedErr error
	}{
		{"00000gy9o", 10241024, false, nil},
		{"0000000GW", 1024, false, nil},
		{"000000GW", 0, true, salesforceid.ErrInvalidNumericIdentifier},
		{"", 0, true, salesforceid.ErrInvalidNumericIdentifier},
		{"0000000000", 0, true, salesforceid.ErrInvalidNumericIdentifier},
		{"000000000", 0, false, nil},
		{"zzzzzzzzz", 13537086546263552 - 1, false, nil}, // 13537086546263552 = 62 ^ 9
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.in, func(t *testing.T) {
			want := tc.out
			got, err := salesforceid.Decode([]byte(tc.in))
			if tc.shouldErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if tc.shouldErr && err != tc.expectedErr {
				t.Errorf("expected err %q but got %q", tc.expectedErr, err)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("didn't expect an error but got %q", err)
			}
			if want != got {
				t.Errorf("wanted %d, got %d", want, got)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	benchmarks := []struct {
		name string
		v    string
	}{
		{"max value", "zzzzzzzzz"},
		{"min value", "000000000"},
		{"middle value", "V00000000"},
	}
	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = salesforceid.Encode(13537086546263552 - 1)
			}
		})
	}
}

// Examples

func ExampleNew() {
	id, err := salesforceid.New("00D000000000062")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("00D000000000062 => %s", id)
	// Output: 00D000000000062 => 00D000000000062EAA
}

func ExampleNew_second() {
	id, err := salesforceid.New("00d000000000062eaa")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("00d000000000062eaa => %s", id)
	// Output: 00d000000000062eaa => 00D000000000062EAA
}

func ExampleDecode() {
	id, _ := salesforceid.New("00D000000000062")
	val, _ := salesforceid.Decode(id.NumericIdentifier)
	fmt.Printf("%s == %d", string(id.NumericIdentifier), val)
	// Output: 000000062 == 374
}

func ExampleDecode_second() {
	id, _ := salesforceid.New("00d000000000062eaa")
	val, _ := salesforceid.Decode(id.NumericIdentifier)
	fmt.Printf("%s == %d", string(id.NumericIdentifier), val)
	// Output: 000000062 == 374
}

func ExampleEncode() {
	id, _ := salesforceid.New("00d000000000062eaa")
	val, _ := salesforceid.Decode(id.NumericIdentifier)
	encoded, _ := salesforceid.Encode(val + 238328) // 238328 == 62 * 62 * 62
	fmt.Printf("%s + 238328 == %s", string(id.NumericIdentifier), encoded)
	// Output: 000000062 + 238328 == 000001062
}

func ExampleEncode_second() {
	id, _ := salesforceid.New("00d000000000062eaa")
	val, _ := salesforceid.Decode(id.NumericIdentifier)
	encoded, _ := salesforceid.Encode(val)
	fmt.Printf("%s == %s", string(id.NumericIdentifier), encoded)
	// Output: 000000062 == 000000062
}

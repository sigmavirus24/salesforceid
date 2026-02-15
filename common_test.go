package salesforceid_test

import (
	"testing"

	"github.com/sigmavirus24/salesforceid"
)

func TestEncode(t *testing.T) {
	testCases := []struct {
		i           uint64
		o           string
		shouldErr   bool
		expectedErr error
	}{
		// Boundaries
		{0, "00000000", false, nil},
		{13537086546263553, "", true, salesforceid.ErrValueTooLarge},
		// Within boundaries
		{5, "00000005", false, nil},
		{15, "0000000F", false, nil},
		{30, "0000000U", false, nil},
		{45, "0000000j", false, nil},
		{61, "0000000z", false, nil},
		{1024, "000000GW", false, nil},
		{10241024, "0000gy9o", false, nil},
		{109_170_052_792_448, "V0000000", false, nil},
		{218_340_105_584_896 - 1, "zzzzzzzz", false, nil},
	}

	for _, tc := range testCases {
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
		{"0000gy9o", 10241024, false, nil},
		{"000000GW", 1024, false, nil},
		{"0000000GW", 0, true, salesforceid.ErrInvalidNumericIdentifier},
		{"", 0, true, salesforceid.ErrInvalidNumericIdentifier},
		{"ßßßß", 0, true, salesforceid.ErrInvalidNumericIdentifier},
		{"0000000000", 0, true, salesforceid.ErrInvalidNumericIdentifier},
		{"00000000", 0, false, nil},
		{"zzzzzzzz", 218_340_105_584_896 - 1, false, nil}, // 218_340_105_584_896 = 62 ^ 8
	}

	for _, tc := range testCases {
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
		v    []byte
	}{
		{"max value", []byte("zzzzzzzzz")},
		{"min value", []byte("000000000")},
		{"middle value", []byte("V00000000")},
	}
	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = salesforceid.Decode(benchmark.v)
			}
		})
	}
}

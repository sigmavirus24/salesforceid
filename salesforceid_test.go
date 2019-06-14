package salesforceid

import (
	"fmt"
	"testing"
)

func TestTo18(t *testing.T) {
	testCases := []struct {
		sfid        string
		expected    string
		shouldErr   bool
		expectedErr error
	}{
		{"00D000000000062EAA", "00D000000000062EAA", false, nil},
		{"00D000000000062", "00D000000000062EAA", false, nil},
		{"00d000000000062", "00d000000000062AAA", false, nil},
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
		{"ZZZZZZZZZZZZZZ", "", true, ErrInvalidLengthSFID},
		{"ZZZZZZZZZZZZZZZZ", "", true, ErrInvalidLengthSFID},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.sfid, func(t *testing.T) {
			t.Parallel()
			eighteen, err := New(tc.sfid).To18()
			if tc.shouldErr && err == nil {
				t.Errorf("expected error but didn't get an error")
			}
			if tc.shouldErr && err != tc.expectedErr {
				t.Errorf("expected err %q, got err %q", tc.expectedErr, err)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("expected no error but got %q", err)
			}
			if eighteen != nil {
				s := eighteen.String()
				if tc.expected != s {
					t.Errorf("expected %s, got %s", tc.expected, s)
				}
			}
		})
	}
}

func BenchmarkTo18(b *testing.B) {
	sfid := New("0a3d0000001ah2a")
	for i := 0; i < b.N; i++ {
		sfid.To18()
	}
}

func TestTo15(t *testing.T) {
	testCases := []struct {
		sfid        string
		expected    string
		shouldErr   bool
		expectedErr error
	}{
		{"ZZZZZZZZZZZZZZ", "", true, ErrInvalidLengthSFID},
		{"ZZZZZZZZZZZZZZZZ", "", true, ErrInvalidLengthSFID},
		{"AAAAAAAAAAAAAAA555", "AAAAAAAAAAAAAAA", false, nil},
		{"aaaaaaaaaaaaaaa555", "AAAAAAAAAAAAAAA", false, nil},
		{"ZZZZZZZZZZZZZZZ555", "ZZZZZZZZZZZZZZZ", false, nil},
		{"zzzzzzzzzzzzzzz555", "ZZZZZZZZZZZZZZZ", false, nil},
		{"ZzZzZzZzZZzZZzzAAA", "zzzzzzzzzzzzzzz", false, nil},
		{"ZZZZZZZZZZZZZZZAAA", "zzzzzzzzzzzzzzz", false, nil},
		{"ZZZZZZZZZZZZZZZ", "ZZZZZZZZZZZZZZZ", false, nil},
		{"00D000000000062EAA", "00D000000000062", false, nil},
		{"00d000000000062eaa", "00D000000000062", false, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.sfid, func(t *testing.T) {
			t.Parallel()
			sfid := New(tc.sfid)
			fifteen, err := sfid.To15()
			if tc.shouldErr && err == nil {
				t.Errorf("expected error but didn't get an error")
			}
			if tc.shouldErr && err != tc.expectedErr {
				t.Errorf("expected err %q, got err %q", tc.expectedErr, err)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("didn't expect an error but got %q", err)
			}
			if fifteen != nil {
				s := fifteen.String()
				if tc.expected != s {
					t.Errorf("expected %s, got %s", tc.expected, s)
				}
			}
		})
	}
}

func BenchmarkTo15(b *testing.B) {
	sfid := New("00d000000000062eaa")
	for i := 0; i < b.N; i++ {
		sfid.To15()
	}
}

func TestNormalize(t *testing.T) {
	testCases := []struct {
		sfid        string
		expected    string
		shouldErr   bool
		expectedErr error
	}{
		{"00d000000000062eaa", "00D000000000062EAA", false, nil},
		{"00D000000000062EAA", "00D000000000062EAA", false, nil},
		{"00D000000000062eaa", "00D000000000062EAA", false, nil},
		{"00d000000000062EAA", "00D000000000062EAA", false, nil},
		{"00D000000000062", "00D000000000062EAA", false, nil},
		{"00D000000000062", "00D000000000062EAA", false, nil},
		{"003D0000001aH2A", "003D0000001aH2AIAU", false, nil},
		{"0a3D0000001aH2A", "0a3D0000001aH2AIAU", false, nil},
		{"0A3D0000001aH2A", "0A3D0000001aH2AKAU", false, nil},
		{"0a3d0000001ah2a", "0a3d0000001ah2aAAA", false, nil},
		{"0a3d0000001ah2aa", "", true, ErrInvalidLengthSFID},
		{"0a3d0000001ah2aaaaa", "", true, ErrInvalidLengthSFID},
		{"0a3d000000", "", true, ErrInvalidLengthSFID},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.sfid, func(t *testing.T) {
			t.Parallel()
			sfid := New(tc.sfid)
			normalized, err := sfid.Normalize()
			if tc.shouldErr && err == nil {
				t.Errorf("expected error but didn't get an error")
			}
			if tc.shouldErr && err != tc.expectedErr {
				t.Errorf("expected err %q, got err %q", tc.expectedErr, err)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("expected no error but got %q", err)
			}
			if normalized != nil {
				s := normalized.String()
				if tc.expected != s {
					t.Errorf("expected %s, got %s", tc.expected, s)
				}
			}
		})
	}
}

func BenchmarkNormalize(b *testing.B) {
	sfid := New("00d000000000062eaa")
	for i := 0; i < b.N; i++ {
		sfid.Normalize()
	}
}

func TestKeyPrefix(t *testing.T) {
	testCases := []struct {
		sfid     string
		expected string
	}{
		{"003D0000001aH2A", "003"},
		{"0a3D0000001aH2A", "0a3"},
		{"0A3D0000001aH2A", "0A3"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.sfid, func(t *testing.T) {
			t.Parallel()
			s := New(tc.sfid)
			p := s.KeyPrefix()
			if p != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, p)
			}
		})
	}
}

func BenchmarkKeyPrefix(b *testing.B) {
	sfid := New("0a3D0000001aH2A")
	for i := 0; i < b.N; i++ {
		sfid.KeyPrefix()
	}
}

func TestPodIdentifier(t *testing.T) {
	testCases := []struct {
		sfid     string
		expected string
	}{
		{"003D0000001aH2A", "003"},
		{"0a3D0000001aH2A", "0a3"},
		{"0A3D0000001aH2A", "0A3"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.sfid, func(t *testing.T) {
			t.Parallel()
			s := New(tc.sfid)
			p := s.KeyPrefix()
			if p != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, p)
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		sfid string
	}{
		{"003D0000001aH2A"},
		{"0a3D0000001aH2A"},
		{"0A3D0000001aH2A"},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.sfid, func(t *testing.T) {
			t.Parallel()
			s := New(tc.sfid).String()
			if s != tc.sfid {
				t.Errorf("expected %s, got %s", tc.sfid, s)
			}
		})
	}
}

func TestSuffix(t *testing.T) {
	testCases := []struct {
		sfid        string
		expected    string
		shouldErr   bool
		expectedErr error
	}{
		{"00D000000000062", "", true, ErrNoSuffix},
		{"00D000000000062EAA", "EAA", false, nil},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.sfid, func(t *testing.T) {
			t.Parallel()
			sfid := New(tc.sfid)
			suffix, err := sfid.Suffix()
			if tc.shouldErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if tc.shouldErr && err != tc.expectedErr {
				t.Errorf("expected err %q but got %q", tc.expectedErr, err)
			}
			if !tc.shouldErr && err != nil {
				t.Errorf("didn't expect an error but got %q", err)
			}
			if suffix != "" && suffix != tc.expected {
				t.Errorf("expected suffix %s, got %s", tc.expected, suffix)
			}
		})
	}
}

// Examples

func ExampleSalesforceID_To18() {
	fifteen := New("00D000000000062")
	eighteen, err := fifteen.To18()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s => %s", fifteen, eighteen)
	// Output: 00D000000000062 => 00D000000000062EAA
}

func ExampleSalesforceID_Normalize() {
	fifteen := New("00D000000000062")
	normalized, err := fifteen.Normalize()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s => %s", fifteen, normalized)
	// Output: 00D000000000062 => 00D000000000062EAA
}

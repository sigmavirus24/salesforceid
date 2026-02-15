package salesforceid_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

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

func TestSalesforceID_Add(t *testing.T) {
	testCases := []struct {
		sfid        string
		add         uint64
		expected    string
		shouldErr   bool
		expectedErr error
	}{
		{"001000000000000AAA", 1, "001000000000001AAA", false, nil},
		{"0010000zzzzzzzzAAA", 1, "", true, salesforceid.ErrInvalidAddition},
		{"00100000000000zAAA", 1, "001000000000010AAA", false, nil},
		{"00100000000000zAAA", 238328, "00100000000100zAAA", false, nil},
		{"001000000000000AAA", 10, "00100000000000AAAQ", false, nil},
	}

	for _, tc := range testCases {
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

func TestSalesforceID_Subtract(t *testing.T) {
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

// func TestParse(t *testing.T) {
// 	testCases := []struct {
// 		sfid        string
// 		expected    string
// 		shouldErr   bool
// 		expectedErr error
// 	}{
// 		{"00D000000000062EAA", "00D000000000062EAA", false, nil},
// 		{"00D000000000062", "00D000000000062EAA", false, nil},
// 		{"00d000000000062", "00d000000000062AAA", false, nil},
// 		{"00d000000000062eaa", "00D000000000062EAA", false, nil},
// 		{"003D0000001aH2A", "003D0000001aH2AIAU", false, nil},
// 		{"0a3D0000001aH2A", "0a3D0000001aH2AIAU", false, nil},
// 		{"0A3D0000001aH2A", "0A3D0000001aH2AKAU", false, nil},
// 		{"0a3d0000001ah2a", "0a3d0000001ah2aAAA", false, nil},
// 		{"000000000000000", "000000000000000AAA", false, nil},
// 		{"999999999999999", "999999999999999AAA", false, nil},
// 		{"aaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaAAA", false, nil},
// 		{"zzzzzzzzzzzzzzz", "zzzzzzzzzzzzzzzAAA", false, nil},
// 		{"AAAAAAAAAAAAAAA", "AAAAAAAAAAAAAAA555", false, nil},
// 		{"ZZZZZZZZZZZZZZZ", "ZZZZZZZZZZZZZZZ555", false, nil},
// 		{"ZZZZZZZZZZZZZZ", "", true, salesforceid.ErrInvalidLengthSFID},   // 14 char sfid
// 		{"ZZZZZZZZZZZZZZZZ", "", true, salesforceid.ErrInvalidLengthSFID}, // 16 char sfid
// 		{"001000000000062EAA", "", true, salesforceid.ErrInvalidSFID},     // 001 should be 00<cap>
// 		{"aaaaaaaaaaaaaaa555", "AAAAAAAAAAAAAAA555", false, nil},
// 		{"ZZZZZZZZZZZZZZZ555", "ZZZZZZZZZZZZZZZ555", false, nil},
// 		{"zzzzzzzzzzzzzzz555", "ZZZZZZZZZZZZZZZ555", false, nil},
// 		{"ZzZzZzZzZZzZZzzAAA", "zzzzzzzzzzzzzzzAAA", false, nil},
// 		{"ZZZZZZZZZZZZZZZAAA", "zzzzzzzzzzzzzzzAAA", false, nil},
// 	}
//
// 	for _, tc := range testCases {
// 		tc := tc
// 		t.Run(tc.sfid, func(t *testing.T) {
// 			t.Parallel()
// 			id, err := salesforceid.Parse(tc.sfid, salesforceid.PostSummer23IdentifierEdition)
// 			if tc.shouldErr && err == nil {
// 				t.Errorf("expected error but didn't get an error")
// 			}
// 			if tc.shouldErr && err != tc.expectedErr {
// 				t.Errorf("expected err %q, got err %q", tc.expectedErr, err)
// 			}
// 			if !tc.shouldErr && err != nil {
// 				t.Errorf("expected no error but got %q", err)
// 			}
// 			if id != nil {
// 				s := id.String()
// 				if tc.expected != s {
// 					t.Errorf("expected %s, got %s", tc.expected, s)
// 				}
// 			}
// 		})
// 	}
// }

func TestParse(t *testing.T) {
	type args struct {
		id      string
		edition salesforceid.IdentifierEdition
	}
	tests := []struct {
		name    string
		args    args
		want    *salesforceid.SalesforceID
		wantErr bool
	}{
		{
			name: "successful parsing no modification necessary (pre Summer '23)",
			args: args{
				id:      "00D000000000062EAA",
				edition: salesforceid.PreSummer23IdentifierEdition,
			},
			want: &salesforceid.SalesforceID{
				KeyPrefix:         []byte("00D"),
				PodIdentifier:     []byte("00"),
				Reserved:          []byte("00"),
				NumericIdentifier: []byte("00000062"),
				Suffix:            []byte("EAA"),
				Edition:           salesforceid.PreSummer23IdentifierEdition,
			},
			wantErr: false,
		},
		{
			name: "successful parsing no modification necessary (post Summer '23)",
			args: args{
				id:      "00D000000000062EAA",
				edition: salesforceid.PostSummer23IdentifierEdition,
			},
			want: &salesforceid.SalesforceID{
				KeyPrefix:         []byte("00D"),
				PodIdentifier:     []byte("000"),
				Reserved:          []byte("0"),
				NumericIdentifier: []byte("00000062"),
				Suffix:            []byte("EAA"),
				Edition:           salesforceid.PostSummer23IdentifierEdition,
			},
			wantErr: false,
		},
		{
			name: "successful parsing calcuating check suffix (pre Summer '23)",
			args: args{
				id:      "00D000000000062",
				edition: salesforceid.PreSummer23IdentifierEdition,
			},
			want: &salesforceid.SalesforceID{
				KeyPrefix:         []byte("00D"),
				PodIdentifier:     []byte("00"),
				Reserved:          []byte("00"),
				NumericIdentifier: []byte("00000062"),
				Suffix:            []byte("EAA"),
				Edition:           salesforceid.PreSummer23IdentifierEdition,
			},
			wantErr: false,
		},
		{
			name: "successful parsing calcuating check suffix (post Summer '23)",
			args: args{
				id:      "00D000000000062",
				edition: salesforceid.PostSummer23IdentifierEdition,
			},
			want: &salesforceid.SalesforceID{
				KeyPrefix:         []byte("00D"),
				PodIdentifier:     []byte("000"),
				Reserved:          []byte("0"),
				NumericIdentifier: []byte("00000062"),
				Suffix:            []byte("EAA"),
				Edition:           salesforceid.PostSummer23IdentifierEdition,
			},
			wantErr: false,
		},
		{
			name: "successful parsing calcuating check suffix (pre Summer '23)",
			args: args{
				id:      "00d000000000062",
				edition: salesforceid.PreSummer23IdentifierEdition,
			},
			want: &salesforceid.SalesforceID{
				KeyPrefix:         []byte("00d"),
				PodIdentifier:     []byte("00"),
				Reserved:          []byte("00"),
				NumericIdentifier: []byte("00000062"),
				Suffix:            []byte("AAA"),
				Edition:           salesforceid.PreSummer23IdentifierEdition,
			},
			wantErr: false,
		},
		{
			name: "successful parsing calcuating check suffix (post Summer '23)",
			args: args{
				id:      "00d000000000062",
				edition: salesforceid.PostSummer23IdentifierEdition,
			},
			want: &salesforceid.SalesforceID{
				KeyPrefix:         []byte("00d"),
				PodIdentifier:     []byte("000"),
				Reserved:          []byte("0"),
				NumericIdentifier: []byte("00000062"),
				Suffix:            []byte("AAA"),
				Edition:           salesforceid.PostSummer23IdentifierEdition,
			},
			wantErr: false,
		},
		{
			name: "successful parsing fixing casing with check suffix (pre Summer '23)",
			args: args{
				id:      "00d000000000062eaa",
				edition: salesforceid.PreSummer23IdentifierEdition,
			},
			want: &salesforceid.SalesforceID{
				KeyPrefix:         []byte("00D"),
				PodIdentifier:     []byte("00"),
				Reserved:          []byte("00"),
				NumericIdentifier: []byte("00000062"),
				Suffix:            []byte("EAA"),
				Edition:           salesforceid.PreSummer23IdentifierEdition,
			},
			wantErr: false,
		},
		{
			name: "successful parsing fixing casing with check suffix (post Summer '23)",
			args: args{
				id:      "00d000000000062eaa",
				edition: salesforceid.PostSummer23IdentifierEdition,
			},
			want: &salesforceid.SalesforceID{
				KeyPrefix:         []byte("00D"),
				PodIdentifier:     []byte("000"),
				Reserved:          []byte("0"),
				NumericIdentifier: []byte("00000062"),
				Suffix:            []byte("EAA"),
				Edition:           salesforceid.PostSummer23IdentifierEdition,
			},
			wantErr: false,
		},
		{
			name:    "generates suffix (object prefixes can have identical check-sums 003)",
			args:    args{id: "003D0000001aH2A", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("003"), PodIdentifier: []byte("D00"), Reserved: []byte("0"), NumericIdentifier: []byte("0001aH2A"), Suffix: []byte("IAU"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "generates suffix (object prefixes can have identical check-sums 0a3)",
			args:    args{id: "0a3D0000001aH2A", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("0a3"), PodIdentifier: []byte("D00"), Reserved: []byte("0"), NumericIdentifier: []byte("0001aH2A"), Suffix: []byte("IAU"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "generates suffix (object prefixes do influence check-sums 0A3)",
			args:    args{id: "0A3D0000001aH2A", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("0A3"), PodIdentifier: []byte("D00"), Reserved: []byte("0"), NumericIdentifier: []byte("0001aH2A"), Suffix: []byte("KAU"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "generates suffix (casing matters)",
			args:    args{id: "0a3d0000001ah2a", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("0a3"), PodIdentifier: []byte("d00"), Reserved: []byte("0"), NumericIdentifier: []byte("0001ah2a"), Suffix: []byte("AAA"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "generates suffix (all zeros)",
			args:    args{id: "000000000000000", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("000"), PodIdentifier: []byte("000"), Reserved: []byte("0"), NumericIdentifier: []byte("00000000"), Suffix: []byte("AAA"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "generates suffix (all nines)",
			args:    args{id: "999999999999999", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("999"), PodIdentifier: []byte("999"), Reserved: []byte("9"), NumericIdentifier: []byte("99999999"), Suffix: []byte("AAA"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "generates suffix (all a's)",
			args:    args{id: "aaaaaaaaaaaaaaa", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("aaa"), PodIdentifier: []byte("aaa"), Reserved: []byte("a"), NumericIdentifier: []byte("aaaaaaaa"), Suffix: []byte("AAA"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "generates suffix (all z's)",
			args:    args{id: "zzzzzzzzzzzzzzz", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("zzz"), PodIdentifier: []byte("zzz"), Reserved: []byte("z"), NumericIdentifier: []byte("zzzzzzzz"), Suffix: []byte("AAA"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "generates suffix (all A's)",
			args:    args{id: "AAAAAAAAAAAAAAA", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("AAA"), PodIdentifier: []byte("AAA"), Reserved: []byte("A"), NumericIdentifier: []byte("AAAAAAAA"), Suffix: []byte("555"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "generates suffix (all Z's)",
			args:    args{id: "ZZZZZZZZZZZZZZZ", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("ZZZ"), PodIdentifier: []byte("ZZZ"), Reserved: []byte("Z"), NumericIdentifier: []byte("ZZZZZZZZ"), Suffix: []byte("555"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "corrects casing with check suffix (all 15 characters are wrong - lower a instead of upper A)",
			args:    args{id: "aaaaaaaaaaaaaaa555", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("AAA"), PodIdentifier: []byte("AAA"), Reserved: []byte("A"), NumericIdentifier: []byte("AAAAAAAA"), Suffix: []byte("555"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "does nothing",
			args:    args{id: "ZZZZZZZZZZZZZZZ555", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("ZZZ"), PodIdentifier: []byte("ZZZ"), Reserved: []byte("Z"), NumericIdentifier: []byte("ZZZZZZZZ"), Suffix: []byte("555"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "corrects casing with check suffix (all 15 characters are wrong - lower z instead of upper Z)",
			args:    args{id: "zzzzzzzzzzzzzzz555", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("ZZZ"), PodIdentifier: []byte("ZZZ"), Reserved: []byte("Z"), NumericIdentifier: []byte("ZZZZZZZZ"), Suffix: []byte("555"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "corrects casing with check suffix (all 15 characters are wrong - upper Z instead of lower z)",
			args:    args{id: "ZZZZZZZZZZZZZZZAAA", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("zzz"), PodIdentifier: []byte("zzz"), Reserved: []byte("z"), NumericIdentifier: []byte("zzzzzzzz"), Suffix: []byte("AAA"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "corrects casing with check suffix (intermittently wrong casing - upper Z instead of lower z)",
			args:    args{id: "ZzZzZzZzZZzZZzzAAA", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    &salesforceid.SalesforceID{KeyPrefix: []byte("zzz"), PodIdentifier: []byte("zzz"), Reserved: []byte("z"), NumericIdentifier: []byte("zzzzzzzz"), Suffix: []byte("AAA"), Edition: salesforceid.PostSummer23IdentifierEdition},
			wantErr: false,
		},
		{
			name:    "errors when identifier is 14 characters",
			args:    args{id: "ZZZZZZZZZZZZZZ", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "errors when identifier is 16 characters",
			args:    args{id: "ZZZZZZZZZZZZZZZZ", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "object prefix is corrupted - check suffix indicate it should be a capital letter instead of a 1",
			args:    args{id: "001000000000062EAA", edition: salesforceid.PostSummer23IdentifierEdition},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := salesforceid.Parse(tt.args.id, tt.args.edition)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Parse() error = %+v, wantErr %+v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if diff := cmp.Diff(got, tt.want, cmpopts.IgnoreUnexported(salesforceid.SalesforceID{})); diff != "" {
				t.Errorf("Parse() = %+v, want %+v, diff = %s", got, tt.want, diff)
			}
		})
	}
}

func BenchmarkNewV2(b *testing.B) {
	b.Run("15 char sfid", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = salesforceid.Parse("0a3d0000001ah2a", salesforceid.PostSummer23IdentifierEdition)
		}
	})
	b.Run("18 char sfid", func(b *testing.B) {
		b.Run("no changes necessary", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = salesforceid.Parse("0A3D0000001aH2AKAU", salesforceid.PostSummer23IdentifierEdition)
			}
		})
		b.Run("everything needs correction", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = salesforceid.Parse("zzzzzzzzzzzzzzz555", salesforceid.PostSummer23IdentifierEdition)
			}
		})
	})
}

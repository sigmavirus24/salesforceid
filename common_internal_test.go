package salesforceid

import (
	"testing"
)

func Test_addToID(t *testing.T) {
	type args struct {
		numeric []byte
		i       uint64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "invalid identifier",
			args:    args{numeric: []byte("ßßßß"), i: 1},
			want:    "",
			wantErr: true,
		},
		{
			name:    "would overflow",
			args:    args{numeric: []byte("zzzzzzzz"), i: 10},
			want:    "",
			wantErr: true,
		},
		{
			name:    "adding nothing to the max is safe",
			args:    args{numeric: []byte("zzzzzzzz"), i: 0},
			want:    "zzzzzzzz",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addToID(tt.args.numeric, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Fatalf("addToID(%v, %v) error = %v, wantErr %v", tt.args.numeric, tt.args.i, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("addToID(%v, %v) = %v, want %v", tt.args.numeric, tt.args.i, got, tt.want)
			}
		})
	}
}

func Test_subtractFromID(t *testing.T) {
	type args struct {
		numeric []byte
		i       uint64
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "invalid identifier",
			args:    args{numeric: []byte("ßßßß"), i: 1},
			want:    "",
			wantErr: true,
		},
		{
			name:    "would overflow",
			args:    args{numeric: []byte("00000000"), i: 10},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := subtractFromID(tt.args.numeric, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Fatalf("subtractFromID(%v, %v) error = %v, wantErr %v", tt.args.numeric, tt.args.i, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("subtractFromID(%v, %v) = %v, want %v", tt.args.numeric, tt.args.i, got, tt.want)
			}
		})
	}
}

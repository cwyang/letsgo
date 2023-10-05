package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// slice of anonymous test
	tests := []struct {
		name string
		tm time.Time
		want string
	}{
		{
			name: "UTC",
			tm: time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want: "17 Dec 2020 at 10:00",
		},
		{
			name: "Empty",
			tm: time.Time{},
			want: "",
		},
		{
			name: "KST",
			tm: time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("KST",-9*60*60)),
			want: "17 Dec 2020 at 19:00",
		},
	}

	for _, tt := range tests {
		// sub-test run
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			if hd != tt.want {
				t.Errorf("want %q, got %q", tt.want, hd)
			}
		})
	}
}

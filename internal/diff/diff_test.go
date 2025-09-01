package diff

import (
	"reflect"
	"testing"
)

func TestPlaylistDiff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		apple      []string
		spotify    []string
		wantAdd    []string
		wantDelete []string
	}{
		{
			name:       "identical",
			apple:      []string{"a", "b"},
			spotify:    []string{"a", "b"},
			wantAdd:    nil,
			wantDelete: nil,
		},
		{
			name:       "add only",
			apple:      []string{"a", "b", "c"},
			spotify:    []string{"a", "b"},
			wantAdd:    []string{"c"},
			wantDelete: nil,
		},
		{
			name:       "delete only",
			apple:      []string{"a"},
			spotify:    []string{"a", "b"},
			wantAdd:    nil,
			wantDelete: []string{"b"},
		},
		{
			name:       "add and delete",
			apple:      []string{"a", "c"},
			spotify:    []string{"a", "b"},
			wantAdd:    []string{"c"},
			wantDelete: []string{"b"},
		},
		{
			name:       "empty inputs",
			apple:      nil,
			spotify:    nil,
			wantAdd:    nil,
			wantDelete: nil,
		},
		{
			name:       "duplicates to add preserved",
			apple:      []string{"x", "y", "y"},
			spotify:    []string{"x"},
			wantAdd:    []string{"y", "y"},
			wantDelete: nil,
		},
		{
			name:       "duplicates to delete preserved",
			apple:      []string{"x"},
			spotify:    []string{"x", "y", "y"},
			wantAdd:    nil,
			wantDelete: []string{"y", "y"},
		},
		{
			name:       "order preserved",
			apple:      []string{"c", "a", "b"},
			spotify:    []string{"b", "c"},
			wantAdd:    []string{"a"},
			wantDelete: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotAdd, gotDelete := PlaylistDiff(tt.apple, tt.spotify)

			if !reflect.DeepEqual(gotAdd, tt.wantAdd) {
				t.Fatalf("toAdd mismatch:\n  got:  %#v\n  want: %#v", gotAdd, tt.wantAdd)
			}
			if !reflect.DeepEqual(gotDelete, tt.wantDelete) {
				t.Fatalf("toDelete mismatch:\n  got:  %#v\n  want: %#v", gotDelete, tt.wantDelete)
			}
		})
	}
}

package store

import (
	"testing"

	"github.com/cunyat/feeder/pkg/utils"
)

func TestStore_Insert(t *testing.T) {
	tests := []struct {
		name     string
		skus     []string
		expDups  int
		expCount int
	}{
		{
			name:     "unque skus",
			skus:     []string{"one", "two", "three", "four"},
			expDups:  0,
			expCount: 4,
		}, {
			name:     "all duplicateds",
			skus:     []string{"one", "one", "one", "one"},
			expDups:  3,
			expCount: 1,
		}, {
			name:     "mix - 2-2",
			skus:     []string{"one", "one", "two", "two"},
			expDups:  2,
			expCount: 2,
		}, {
			name:     "empty",
			skus:     []string{},
			expDups:  0,
			expCount: 0,
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := New()
			for _, sku := range tt.skus {
				db.Insert(sku)
			}
			dups := db.DuplicatedCount()
			if dups != tt.expDups {
				t.Errorf("expected duplicates: %d, got: %d", tt.expDups, dups)
			}

			count := db.SKUCount()
			if count != tt.expCount {
				t.Errorf("expected count: %d, got: %d", tt.expCount, count)
			}
		})
	}
}

func BenchmarkStore(b *testing.B) {
	db := New()
	var skus = utils.GenerateSKUs(100000)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		for _, sku := range skus {
			db.Insert(sku)
		}
	}
}

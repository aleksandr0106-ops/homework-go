package storage

import (
	"homework-go/homework-7/internal/item"
	"path/filepath"
	"testing"
)

func BenchmarkStorage_Save(b *testing.B) {
	items := make([]item.Item, 100)
	for i := 0; i < 100; i++ {
		items[i] = item.NewWeapon(i, "Меч", 5.0, 10, 5)
	}

	b.Run("MemoryStorage", func(b *testing.B) {
		store := NewMemoryStorage()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = store.Save(items)
		}
	})

	b.Run("FileStorage", func(b *testing.B) {
		tmpDir := b.TempDir()
		filePath := filepath.Join(tmpDir, "bench.txt")
		store := NewFileStorage(filePath)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = store.Save(items)
		}
	})
}

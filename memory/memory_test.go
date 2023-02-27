package memory_test

import (
	"fmt"
	"github.com/aacfactory/systems/memory"
	"testing"
)

func TestStats(t *testing.T) {
	stats, err := memory.Stats()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(stats, stats.Heap)
}

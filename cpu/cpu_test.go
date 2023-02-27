package cpu_test

import (
	"fmt"
	"github.com/aacfactory/systems/cpu"
	"testing"
)

func TestOccupancy(t *testing.T) {
	r, err := cpu.Occupancy()
	fmt.Println(r, err)
	fmt.Println(r.Min(), r.Max(), r.AVG())
}

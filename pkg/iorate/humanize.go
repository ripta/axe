package iorate

import (
	"fmt"
	"math"
)

// GenerateHumanize generates custom humanization routines.
func GenerateHumanize(sz, base float64, units []string) string {
	magnitude := math.Floor(math.Log(sz) / math.Log(base))
	value := sz / math.Pow(base, magnitude)

	if int(magnitude) >= len(units) {
		return fmt.Sprintf("!VALUE TOO LARGE (%d)", sz)
	}

	var sigdigits int
	if value < 10 {
		sigdigits = 2
	} else if value < 100 {
		sigdigits = 1
	}

	suffix := units[int(magnitude)]
	return fmt.Sprintf(fmt.Sprintf("%%.%df%%s", sigdigits), value, suffix)
}

// HumanizeBinary formats sz to a human-readable IEC bytes (2^10).
func HumanizeBinary(sz float64) string {
	return GenerateHumanize(sz, 1024, []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"})
}

// HumanizeBytes formats sz to a human-readable SI bytes (10^3).
func HumanizeBytes(sz float64) string {
	return GenerateHumanize(sz, 1000, []string{"B", "kB", "MB", "GB", "TB", "PB"})
}

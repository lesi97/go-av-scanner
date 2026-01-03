package utils

import "fmt"

func FormatBytes(bytes int64) string {
	const (
		kib = 1 << 10
		mib = 1 << 20
		gib = 1 << 30

		kb = 1_000
		mb = 1_000_000
		gb = 1_000_000_000
	)

	switch {
	case bytes >= gib:
		return fmt.Sprintf(
			"%.2f GiB (%.2f GB)",
			float64(bytes)/float64(gib),
			float64(bytes)/float64(gb),
		)
	case bytes >= mib:
		return fmt.Sprintf(
			"%.2f MiB (%.2f MB)",
			float64(bytes)/float64(mib),
			float64(bytes)/float64(mb),
		)
	case bytes >= kib:
		return fmt.Sprintf(
			"%.2f KiB (%.2f KB)",
			float64(bytes)/float64(kib),
			float64(bytes)/float64(kb),
		)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

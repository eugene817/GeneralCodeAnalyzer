package services

import (
	"fmt"
	"time"
)

func timeTrack(start time.Time, name string) string {
	elapsed := time.Since(start)
	return fmt.Sprintf("%s took %s", name, elapsed)
}

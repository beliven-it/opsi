package helpers

import (
	"fmt"
	"time"
)

func Log(message string) {
	now := time.Now().Format(time.RFC3339)

	fmt.Printf("[%s] %s \n", now, message)
}

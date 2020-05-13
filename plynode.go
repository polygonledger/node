package main

import (
	"fmt"
	"os"
)

func main() {
	GitCommit := os.Getenv("GIT_COMMIT")
	fmt.Printf("--- run polygon ---\ngit commit: %s ----\n", GitCommit)

	runNodeWithConfig()
}

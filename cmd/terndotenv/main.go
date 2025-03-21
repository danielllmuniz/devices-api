package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	args := []string{
		"migrate",
	}

	// Check for a custom steps argument like steps=1
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "steps=") {
			// You can customize this destination flag as needed
			args = []string{
				"migrate",
				"--destination", "-1",
			}
			break
		}
	}

	// Add common flags
	args = append(args,
		"--migrations", "./internal/store/pgstore/migrations",
		"--config", "./internal/store/pgstore/migrations/tern.conf",
	)

	cmd := exec.Command("tern", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Command execution failed: ", err)
		fmt.Println("Output: ", string(output))
		panic(err)
	}

	fmt.Println("Command executed successfully:", string(output))
}

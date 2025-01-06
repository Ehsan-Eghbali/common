package utils

import (
	"fmt"
	"os"
)

func GetHomeDir() string {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("could not get home directory")
	}

	return currentDir
}

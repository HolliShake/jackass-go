package main

import "os"

func ReadFile(filePath string) string {
	data, err := os.ReadFile(filePath)

	if err != nil {
		panic(err)
	}

	return string(data)
}

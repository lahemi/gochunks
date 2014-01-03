// Get the name of the current directory
package main

import (
	"os"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for i := len(pwd) - 1; i >= 0; i-- {
		if pwd[i] == '/' {
			os.Stdout.Write([]byte(pwd[i+1:]))
			break
		}
	}
}

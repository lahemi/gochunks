package main

import (
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) > 1 {
		arg, err := strconv.ParseInt(os.Args[1], 10, 64)
		if err != nil {
            return
		}
		time.Sleep(time.Duration(arg) * time.Second)
	}
}

// init things for a user
package main

import (
	"fmt"
	"os"
)

func CreateLink() {
	//use user provied source
	// target is fixed
	err := os.Link("test", "logs")
	if err != nil {
		fmt.Println("create link error", err)
		return
	}
	os.Chdir("test")
	if err != nil {
		fmt.Println("change directory error", err)
		return
	}
}

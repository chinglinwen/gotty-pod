// init things for a user
package main

import (
	"os"
)

func CreateLink() {
	//use user provied source
	// target is fixed
	os.Symlink("test", "logs")
}

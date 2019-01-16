// init things for a user
package main

import (
	"os"
)

func CreateLink(src, dst string) {
	//use user provied source
	// target is fixed
	os.Symlink(src, dst)
}

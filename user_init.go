// init things for a user
package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateLink(src, dst string) (err error) {
	// check src exist
	_, err = os.Stat(src)
	if err != nil {
		return fmt.Errorf("%v does not exist", src) //dst is git
	}

	os.MkdirAll(filepath.Dir(dst), 0755)
	//use user provied source
	// target is fixed, only exist error, ignore it
	os.Symlink(src, dst)

	return
}

func RemoveLink(dst string) {
	os.RemoveAll(filepath.Dir(dst))
}

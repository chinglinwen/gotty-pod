package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestA(t *testing.T) {
	base := "/data/fluentd"
	err := filepath.Walk(base,
		func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !f.IsDir() {
				return nil
			}
			fmt.Println(path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

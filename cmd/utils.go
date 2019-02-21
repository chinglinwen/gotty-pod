package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tidwall/gjson"
)

func Walk(base string) (list []string, err error) {
	// base := "/data/fluentd"
	err = filepath.Walk(base,
		func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !f.IsDir() {
				return nil
			}
			git := strings.Split(path, "/")
			if len(git) != 5 {
				return nil
			}
			list = append(list, fmt.Sprintf("%v/%v", git[3], git[4]))
			return nil
		})
	if err != nil {
		return
	}
	return
}

func Filter(dirlist, gitlist []string) []string {
	loglist := []string{}
	for _, v1 := range dirlist {
		for _, v2 := range gitlist {
			git := strings.Replace(v2, "_", "-", -1)
			if v1 == git {
				loglist = append(loglist, v2) // we need gitlab name
			}
		}
	}
	return loglist
}

func Compress(data string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(data)); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

func UnCompress(str string) (s string, err error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return
	}
	rdata := bytes.NewReader(data)
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}
	s = string(b)
	return
}

func ParseUserInfo(str string) (name string, id int, err error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return
	}
	// fmt.Println(string(data))
	id = int(gjson.GetBytes(data, "id").Int())
	name = gjson.GetBytes(data, "username").String()
	return
}

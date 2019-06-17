package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"./k8s"
)

func printprogress() func() {
	ctx, cancel := context.WithCancel(context.Background())

	chars := []string{"/", "-", "\\", "|"}

	go func() {
		for i := 0; ; i++ {
			fmt.Printf("searching projects... %v\r", chars[i%4])
			time.Sleep(100 * time.Millisecond)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
	return cancel
}

func listpods() (list map[string]k8s.Pod, err error) {
	pods, err := k8s.PodItems()
	if err != nil {
		return
	}
	list = make(map[string]k8s.Pod)
	for _, v := range pods {
		// fmt.Printf("got %v,%v\n", v.Name, v.Namespace)
		// list = append(list, v.Namespace+"/"+v.Name)
		list[v.Namespace+"/"+v.PodName] = v
	}
	return
}

// func Walk(base string) (list []string, err error) {
// 	// base := "/data/fluentd"
// 	err = filepath.Walk(base,
// 		func(path string, f os.FileInfo, err error) error {
// 			if err != nil {
// 				return err
// 			}
// 			if !f.IsDir() {
// 				return nil
// 			}
// 			git := strings.Split(path, "/")
// 			if len(git) != 5 {
// 				return nil
// 			}
// 			list = append(list, fmt.Sprintf("%v/%v", git[3], git[4]))
// 			return nil
// 		})
// 	if err != nil {
// 		return
// 	}
// 	return
// }

func Filter(dirlist, gitlist []string) []string {
	loglist := []string{}
	for _, v1 := range dirlist {
		for _, v2 := range gitlist {
			git := strings.Replace(v2, "_", "-", -1)
			// log.Println("compare", v1, v2)
			if strings.HasPrefix(v1, git) { // support pod suffixs
				// log.Println("got", v1)
				loglist = append(loglist, v1) // we need pod name
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

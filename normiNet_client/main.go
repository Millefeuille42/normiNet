package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func getFileList(files []string, prefix string) ([]string, error) {
	fileList := make([]string, 0)

	for _, file := range files {
		file = strings.TrimSuffix(file, "/")
		info, err := os.Stat(prefix + file)
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			fFile, err := os.Open(prefix + file)
			if err != nil {
				return nil, err
			}
			dirFileList, err := fFile.Readdirnames(-1)
			if err != nil {
				return nil, err
			}
			nFileList, err := getFileList(dirFileList, prefix+file+"/")
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, nFileList...)
			continue
		}
		if ext := filepath.Ext(prefix + file); ext == ".c" || ext == ".h" {
			fileList = append(fileList, prefix+file)
		}
	}
	return fileList, nil
}

func main() {
	userName := "anon"

	if currUser, err := user.Current(); err == nil {
		userName = currUser.Username + "_-_" + currUser.Name
	}

	fileList, err := getFileList(os.Args[1:], "")
	if err != nil {
		panic(err)
	}

	resp, err := http.Get("http://localhost:8080/norm")
	if err != nil {
		panic(err)
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(respData))
	for _, file := range fileList {
		info, err := os.Stat(file)
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		resp, err := http.Post("http://localhost:8080/norm?username="+userName+"&filename="+info.Name(),
			"text/x-c", bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Print(string(respData))
	}
}

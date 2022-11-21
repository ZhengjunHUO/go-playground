package main

import (
	"fmt"
	"net/http"

	eg "golang.org/x/sync/errgroup"
)

var (
	urls = []string{
		"http://www.golang.org/",
		//"http://www.huo.com/",
		"http://www.huozj.com/",
	}
)

func main() {
	seq_curl()
}

func seq_curl() {
	grp := new(eg.Group)

	for _, url := range urls {
		url := url
		grp.Go(func() error{
			resp, err := http.Get(url)
			if err == nil {
				resp.Body.Close()
			}
			return err
		})
	}

	if err := grp.Wait(); err == nil {
		fmt.Println("Ok")
	}else{
		fmt.Println("Get some errors when fetching the urls: ", err)
	}
}

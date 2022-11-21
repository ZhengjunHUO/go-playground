package main

import (
	"os"
	"fmt"
	"sync"
	"context"
	"net/http"

	eg "golang.org/x/sync/errgroup"
)

type concatErr struct {
	errMsgs		string
	lock		sync.RWMutex
}

func (ce *concatErr) Merge(err error) {
	ce.lock.Lock()
	defer ce.lock.Unlock()
	ce.errMsgs += "\n  " + err.Error()
}

var (
	urls = []string{
		"http://www.golang.org/",
		"http://www.huozj.com/",
		"http://www.huo.com/",
	}

        repo = pseudoQueryWithSuccess("repo")
	//docs = pseudoQueryWithSuccess("docs")
	//refs = pseudoQueryWithSuccess("refs")
	docs = pseudoQueryWithFailure("docs")
	refs = pseudoQueryWithFailure("refs")
)

type Rslt string
type Query func(ctx context.Context, enq string) (Rslt, error)

func pseudoQueryWithSuccess(name string) Query {
	return func(_ context.Context, enq string) (Rslt, error) {
		return Rslt(fmt.Sprintf("Got content %s of %s", enq, name)), nil
	}
}

func pseudoQueryWithFailure(name string) Query {
	return func(_ context.Context, enq string) (Rslt, error) {
		return Rslt(""), fmt.Errorf("Failed to get %s content of %s", enq, name)
	}
}

func main() {
	//seq_curl()
        paral_curl()
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

func paral_curl() {
	ce := &concatErr{
		errMsgs: "Error Trace: ",
	}

	queries := func(ctx context.Context, qry string) ([]Rslt, error) {
		grp, ctx := eg.WithContext(ctx)

		//qs := []Query{repo, docs, refs}
		qs := []Query{repo, refs, docs}
		rslts := make([]Rslt, len(qs))

		for i, q := range qs {
			i, q := i, q
			grp.Go(func() error {
				rslt, err := q(ctx, qry)
				if err == nil {
					rslts[i] = rslt
					return nil
				}

				ce.Merge(err)
				//fmt.Printf("[inside goroutine]: %s\n", ce.errMsgs)
				return fmt.Errorf("%s", ce.errMsgs)
			})
		}

		// grp.Wait() returns the first non-nil error (if any) from them, use ce.errMsgs instead.
		if err := grp.Wait(); err != nil {
			return nil, fmt.Errorf("%s", ce.errMsgs)
		}
		return rslts, nil
        }

	rslts, err := queries(context.Background(), "rust")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, rslt := range rslts {
		fmt.Println(rslt)
	}
}

package main

import (
	"coursecodes/internal/flow096"
	"coursecodes/internal/store"
	"flag"
	"fmt"
	"os"
)

func main() {
	path := flag.String("db", "redeem.db", "database path")
	flag.Parse()
	s, e := store.Open(*path)
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		os.Exit(1)
	}
	defer s.Close()
	f := flow096.New(s)
	if flag.NArg() > 0 && flag.Arg(0) == "demo" {
		r, e := f.Register("demo-1", "COURSE-001", "Go", 10)
		fmt.Println(r.ID, e)
		return
	}
	rows, e := s.ListRecords()
	fmt.Printf("records=%d err=%v\n", len(rows), e)
}

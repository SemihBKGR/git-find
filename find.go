package main

import (
	"regexp"
	"strings"
	"sync"
)

type findOptions struct {
	ignoreCase bool
	removed    bool
	regex      bool
}

type findResult struct {
	searchTerm string
	file       *diffFile
	lines      []*diffLine
}

func newFindResult() findResult {
	return findResult{
		lines: make([]*diffLine, 0, 10),
	}
}

func find(diff *diff, sts []string, fo findOptions) <-chan findResult {
	c := make(chan findResult)
	wg := sync.WaitGroup{}
	wg.Add(len(sts))
	for _, st := range sts {
		go func(st string) {
			findBySearchTerm(diff, st, fo, c)
			wg.Done()
		}(st)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	return c
}

func findBySearchTerm(diff *diff, st string, fo findOptions, c chan<- findResult) {
	if !fo.regex {
		findByKeyword(diff, st, fo, c)
	} else {
		regexp, err := regexp.Compile(st)
		if err != nil {
			panic(err)
		}
		findByRegex(diff, regexp, fo, c)
	}
}

func findByKeyword(diff *diff, keyword string, fo findOptions, c chan<- findResult) {
	wg := sync.WaitGroup{}
	wg.Add(len(diff.files))
	if fo.ignoreCase {
		keyword = strings.ToLower(keyword)
	}
	for _, df := range diff.files {
		go func(df *diffFile) {
			r := newFindResult()
			r.searchTerm = keyword
			r.file = df
			for _, dl := range df.lines {
				lineContent := dl.content
				if fo.ignoreCase {
					lineContent = strings.ToLower(lineContent)
				}
				if strings.Contains(lineContent, keyword) {
					if !fo.removed && !dl.added {
						continue
					}
					r.lines = append(r.lines, dl)
				}
			}
			c <- r
			wg.Done()
		}(df)
	}
	wg.Wait()
}

func findByRegex(diff *diff, regexp *regexp.Regexp, fo findOptions, c chan<- findResult) {
	wg := sync.WaitGroup{}
	wg.Add(len(diff.files))
	for _, df := range diff.files {
		go func(df *diffFile) {
			r := newFindResult()
			r.searchTerm = regexp.String()
			r.file = df
			for _, dl := range df.lines {
				if regexp.MatchString(dl.content) {
					if !fo.removed && !dl.added {
						continue
					}
					r.lines = append(r.lines, dl)
				}
			}
			c <- r
			wg.Done()
		}(df)
	}
	wg.Wait()
}

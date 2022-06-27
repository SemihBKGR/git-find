package main

import (
	"strings"
	"sync"
)

type findOptions struct {
	ignoreCase bool
	removed    bool
}

type findResult struct {
	file        *diffFile
	occurrences []*lineOccurrence
}

func newFindResult() findResult {
	return findResult{
		occurrences: make([]*lineOccurrence, 0, 10),
	}
}

type lineOccurrence struct {
	line        *diffLine
	occurrences map[string][]uint
}

func newLineOccurrences() *lineOccurrence {
	return &lineOccurrence{
		occurrences: make(map[string][]uint),
	}
}

func find(diff *diff, sts []string, fo findOptions) <-chan findResult {
	c := make(chan findResult)
	wg := sync.WaitGroup{}
	wg.Add(len(diff.files))
	for _, df := range diff.files {
		go func(df *diffFile) {
			r := findByKeyword(df, sts, fo)
			c <- r
			wg.Done()
		}(df)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	return c
}

func findByKeyword(df *diffFile, keywords []string, fo findOptions) findResult {
	r := newFindResult()
	r.file = df
	if fo.ignoreCase {
		for i, keyword := range keywords {
			keywords[i] = strings.ToLower(keyword)
		}
	}
	for _, dl := range df.lines {
		lineContent := dl.content
		if fo.ignoreCase {
			lineContent = strings.ToLower(lineContent)
		}
		ko := findKeywordOccurrences(lineContent, keywords)
		if len(ko) > 0 {
			lo := newLineOccurrences()
			lo.line = dl
			lo.occurrences = ko
			r.occurrences = append(r.occurrences, lo)
		}
	}
	return r
}

func findKeywordOccurrences(s string, ks []string) map[string][]uint {
	m := make(map[string][]uint)
	for _, k := range ks {
		kLen := len(k)
		from := 0
		for i := strings.Index(s, k); i != -1; i = strings.Index(s[from:], k) {
			if indexes, ok := m[k]; ok {
				m[k] = append(indexes, uint(i+from))
			} else {
				indexes := make([]uint, 0, 10)
				m[k] = append(indexes, uint(i+from))
			}
			from += i + kLen
		}
	}
	return m
}

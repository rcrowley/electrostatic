package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type Inputs struct {
	pathnames []string
}

func (in *Inputs) Add(pathname string) {
	ext := filepath.Ext(pathname)
	if ext != ".html" && ext != ".md" {
		return
	}

	// If the Markdown variant of this pathname is in the list already,
	// we're done.
	mdPathname := fmt.Sprint(strings.TrimSuffix(pathname, ext), ".md")
	//log.Printf("Markdown %v", mdPathname)
	i := sort.SearchStrings(in.pathnames, mdPathname)
	if i < len(in.pathnames) && in.pathnames[i] == mdPathname {
		//log.Print("already got one")
		return
	}

	// If the HTML variant of this pathname is in the list already, convert
	// its extension to this extension.
	htmlPathname := fmt.Sprint(strings.TrimSuffix(pathname, ext), ".html")
	//log.Printf("HTML %v", htmlPathname)
	i = sort.SearchStrings(in.pathnames, htmlPathname)
	if i < len(in.pathnames) && in.pathnames[i] == htmlPathname {
		//log.Printf("replacing %v with %v", htmlPathname, pathname)
		in.pathnames[i] = pathname
		return
	}

	i = sort.SearchStrings(in.pathnames, pathname)
	if i == len(in.pathnames) || in.pathnames[i] != pathname {
		in.pathnames = append(in.pathnames, "")
		copy(in.pathnames[i+1:], in.pathnames[i:])
		in.pathnames[i] = pathname
	}
}

func (in *Inputs) Pathnames() []string {
	return in.pathnames
}

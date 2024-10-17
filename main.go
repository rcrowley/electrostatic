package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rcrowley/mergician/files"
	"github.com/rcrowley/mergician/html"
)

func init() {
	log.SetFlags(0)
}

func main() {
	input := flag.String("i", "", "directory containing input HTML and Markdown documents")
	layout := flag.String("l", "", "")
	output := flag.String("o", "", "document root directory where merged HTML documents will be placed")
	pretend := flag.Bool("p", false, "pretend to process all the inputs but don't write any outputs; implies -v")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: electrostatic -i <input> -l <layout> -o <output> [-p] [-v]
  -i <input>   directory containing input HTML and Markdown documents
  -l <layout>  site layout HTML document
  -o <output>  document root directory where merged HTML documents will be placed
  -p           pretend to process all the inputs but don't write any outputs; implies -v
  -v           verbose mode
`)
	}
	flag.Parse()
	if *input == "" || *layout == "" || *output == "" || flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}
	if *pretend {
		*verbose = true
	}

	l := &files.List{}
	must(fs.WalkDir(
		os.DirFS(*input),
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.Type().IsRegular() {
				l.Add(path)
			}
			return nil
		},
	))

	in0 := must2(html.ParseFile(*layout))

	var rules []html.Rule // TODO -r option
	if len(rules) == 0 {
		rules = html.DefaultRules()
	}

	for _, pathname := range l.Pathnames() {
		in := filepath.Join(*input, pathname)
		out := filepath.Join(*output, fmt.Sprint(strings.TrimSuffix(pathname, filepath.Ext(pathname)), ".html"))
		if *verbose {
			fmt.Printf(
				// "mergician -o %q %q %q\n",
				"mergician -o %s %s %s\n",
				out, *layout, in)
		}
		if !*pretend {
			must(os.MkdirAll(filepath.Dir(out), 0777))
			must(html.RenderFile(out, must2(html.Merge([]*html.Node{
				in0,
				must2(files.Parse(in)),
			}, rules))))
		}
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func must2[T any](v T, err error) T {
	must(err)
	return v
}

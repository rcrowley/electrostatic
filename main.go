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
	rules := new(html.Rules)
	flag.Var(rules, "r", "use a custom rule for merging inputs (overrides all defaults; may be repeated)")
	verbose := flag.Bool("v", false, "verbose mode")
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: electrostatic -i <input> -l <layout> -o <output> [-p] [-r <rule>[...]] [-v]
  -i <input>   directory containing input HTML and Markdown documents
  -l <layout>  site layout HTML document
  -o <output>  document root directory where merged HTML documents will be placed
  -p           pretend to process all the inputs but don't write any outputs; implies -v
  -r <rule>    use a custom rule for merging inputs (overrides all defaults;
               may be repeated)
               each rule is a destination HTML tag with optional attributes,
               "=" or "+=", and a source HTML tag with optional attributes
               default rules: <article class="body"> = <body>
                              <div class="body"> = <body>
                              <section class="body"> = <body>
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

	if len(*rules) == 0 {
		*rules = html.DefaultRules()
	}

	for _, pathname := range l.Pathnames() {
		inPathname := filepath.Join(*input, pathname)
		outPathname := filepath.Join(*output, fmt.Sprint(strings.TrimSuffix(pathname, filepath.Ext(pathname)), ".html"))

		in1 := must2(files.Parse(inPathname))
		in := must2(html.Merge([]*html.Node{in0, in1}, *rules))

		if *verbose {
			fmt.Printf(
				"mergician -o %s %s %s\n", // "mergician -o %q %q %q\n",
				outPathname, *layout, inPathname,
			)
		}
		if *pretend {
			continue
		}

		must(os.MkdirAll(filepath.Dir(outPathname), 0777))
		must(html.RenderFile(outPathname, in))
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

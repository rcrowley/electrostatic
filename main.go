package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rcrowley/mergician/files"
	"github.com/rcrowley/mergician/html"
	"golang.org/x/net/html/atom"
)

func Main(args []string, stdin io.Reader, stdout io.Writer) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	author := flags.String("a", "", "author's name (used to build an Atom feed)")
	input := flags.String("i", "", "directory containing input HTML and Markdown documents")
	layout := flags.String("l", "", "")
	output := flags.String("o", "", "document root directory where merged HTML documents will be placed")
	pretend := flags.Bool("p", false, "pretend to process all the inputs but don't write any outputs; implies -v")
	rules := new(html.Rules)
	flags.Var(rules, "r", "use a custom rule for merging inputs (overrides all defaults; may be repeated)")
	url := flags.String("u", "", "site URL with scheme and domain (used to build an Atom feed)")
	verbose := flags.Bool("v", false, "verbose mode")
	flags.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: electrostatic [-a <author>] -i <input> -l <layout> -o <output> [-p] [-r <rule>[...]] [-u <url>] [-v]
  -a <author>  author's name (used to build an Atom feed)
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
  -u <url>     site URL with scheme and domain (used to build an Atom feed)
  -v           verbose mode
`)
	}
	flags.Parse(args[1:])
	if *input == "" || *layout == "" || *output == "" || flags.NArg() > 0 {
		flags.Usage()
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

	feed := &Feed{
		Author: *author,
		Path:   "index.atom.xml", // maybe parameterize
		URL:    *url,
	}
	if title := html.Find(in0, html.IsAtom(atom.Title)); title != nil {
		feed.Title = html.Text(title).String()
	}

	if len(*rules) == 0 {
		*rules = html.DefaultRules()
	}

	var wg sync.WaitGroup
	for _, pathname := range l.Pathnames() {
		inPathname := filepath.Join(*input, pathname)
		outPathname := filepath.Join(*output, fmt.Sprint(strings.TrimSuffix(pathname, filepath.Ext(pathname)), ".html"))
		if *verbose {
			fmt.Printf(
				"mergician -o %s %s %s\n", // "mergician -o %q %q %q\n",
				outPathname, *layout, inPathname,
			)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			in1 := must2(files.Parse(inPathname))
			in := must2(html.Merge([]*html.Node{in0, in1}, *rules))

			if a := html.Find(in, html.IsAtom(atom.Article)); a != nil {
				if t := html.Find(a, html.IsAtom(atom.Time)); t != nil {
					feed.Add(html.Attr(t, "datetime"), outPathname, in)
				}
			}

			if *pretend {
				return
			}

			must(os.MkdirAll(filepath.Dir(outPathname), 0777))
			must(html.RenderFile(outPathname, in))
		}()
	}
	wg.Wait()

	if *author != "" && *url != "" {
		for _, entry := range feed.Entries {
			if *verbose {
				fmt.Printf(
					"frag %s %s # %s\n", // "frag %q %q # %s\n",
					"<h1>", entry.Path, entry.Date,
				)
				fmt.Printf(
					"frag %s %s # %s\n", // "frag %q %q # %s\n",
					"<article>", entry.Path, entry.Date,
				)
			}
		}
		if !*pretend {
			pathname := filepath.Join(*output, feed.Path)
			if *verbose {
				fmt.Printf("# wrote Atom feed to %s\n", pathname)
			}
			must(feed.RenderFile(pathname))
		}
	}

}

func init() {
	log.SetFlags(0)
}

func main() {
	Main(os.Args, os.Stdin, os.Stdout)
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

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/rcrowley/mergician/files"
	"github.com/rcrowley/mergician/html"
)

func Main(args []string, stdin io.Reader, stdout io.Writer) {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	layout := flags.String("l", "", "site layout HTML document")
	output := flags.String("o", ".", "document root directory where merged HTML documents will be placed")
	pretend := flags.Bool("p", false, "pretend to process all the inputs but don't write any outputs; implies -v")
	rules := new(html.Rules)
	flags.Var(rules, "r", "use a custom rule for merging inputs (overrides all defaults; may be repeated)")
	verbose := flags.Bool("v", false, "verbose mode")
	exclude := files.NewStringSliceFlag(flags, "x", "subdirectory of <input> to exclude (may be repeated)")
	flags.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: electrostatic -l <layout> [-o <output>] [-p] [-r <rule>[...]] [-v] [-x <exclude>[...]] <input>[...]
  -l <layout>   site layout HTML document
  -o <output>   document root directory where merged HTML documents will be placed (defaults to the current working directory)
  -p            pretend to process all the inputs but don't write any outputs; implies -v
  -r <rule>     use a custom rule for merging inputs (overrides all defaults;
                may be repeated)
                each rule is a destination HTML tag with optional attributes,
                "=" or "+=", and a source HTML tag with optional attributes
                default rules: <article class="body"> = <body>
                               <div class="body"> = <body>
                               <section class="body"> = <body>
  -v            verbose mode
  -x <exclude>  subdirectory of <input> to exclude (may be repeated)
  <input>       directory containing input HTML and Markdown documents (may be repeated)

Synopsis: electrostatic uses mergician to apply a consistent layout to a whole site.
`)
	}
	flags.Parse(args[1:])
	if *layout == "" || flags.NArg() == 0 {
		flags.Usage()
		os.Exit(1)
	}
	if *pretend {
		*verbose = true
	}

	lists := must2(files.AllInputs(flags.Args(), *exclude))

	in0 := must2(html.ParseFile(*layout))

	if len(*rules) == 0 {
		*rules = html.DefaultRules()
	}

	var wg sync.WaitGroup
	for _, list := range lists {
		for _, path := range list.RelativePaths() {
			inPath := filepath.Join(list.Root(), path)
			outPath := filepath.Join(*output, fmt.Sprint(strings.TrimSuffix(path, filepath.Ext(path)), ".html"))
			if *verbose {
				fmt.Printf(
					"mergician -o %s %s %s %s\n",
					outPath, rules, *layout, inPath,
				)
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				in1 := must2(files.Parse(inPath))
				in := must2(html.Merge([]*html.Node{in0, in1}, *rules))

				if *pretend {
					return
				}

				must(os.MkdirAll(filepath.Dir(outPath), 0777))
				must(html.RenderFile(outPath, in))
			}()
		}
	}
	wg.Wait()
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

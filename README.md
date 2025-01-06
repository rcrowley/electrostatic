Electrostatic
=============

HTML is the language of the Web but the Web has moved beyond the bare, unstyled Web of the '90s (and university professors). Visitors to your site expect coherent layouts, visual consistency, navigation. This is the point at which most folks reach for a templating language. [Mergician](https://github.com/rcrowley/mergician) reimagines the often frustrating juxtaposition of HTML and templating language in pure HTML. Instead of rendering a template, merge a content HTML document into a layout HTML document.

Electostatic extends Mergician just that little bit further into a very basic CMS. It processes a whole document root directory, rendering every Markdown document to HTML and wrapping every HTML document in a consistent layout. Someday it'll learn to generate reverse-chronological index pages. Maybe it'll even learn how to manage email subscribers, too.

[Feed](https://github.com/rcrowley/feed) can complete your site with an Atom feed. You might also want to test your inputs with [Deadlinks](https://github.com/rcrowley/deadlinks) before running Electrostatic or use [Sitesearch](https://github.com/rcrowley/sitesearch) to offer search on Electrostatic's output.

The easiest way to use Electrostatic is with a `make`(1) target about like this:

```make
all:
	electrostatic -l design-system/index.html -v raw
```

When run in the document root directory, this command turns <https://rcrowley.org/design-system/> and <https://rcrowley.org/raw/homelab/index.md> (for example) into <https://rcrowley.org/homelab/>, or <https://rcrowley.org/design-system/> and <https://rcrowley.org/raw/src-bin/> (for another) into <https://rcrowley.org/src-bin/>.

Installation
------------

```sh
go install github.com/rcrowley/electrostatic@latest
```

Usage
-----

```sh
electrostatic -l <layout> -o <output> [-p] [-r <rule>[...]] [-v] [-x <exclude>[...]] <input>[...]
```

* `-l <layout>`: site layout HTML document
* `-o <output>`: document root directory where merged HTML documents will be placed
* `-p`: pretend to process all the inputs but don't write any outputs; implies `-v`
* `-r <rule>`: use a custom rule for merging inputs (overrides all defaults; may be repeated); each rule is a destination HTML tag with optional attributes, "=" or "+=", and a source HTML tag with optional attributes default rules:
    * `<article class="body"> = <body>`
    * `<div class="body"> = <body>`
    * `<section class="body"> = <body>`
* `-v`: verbose mode
* `-x <exclude>`: subdirectory of `<input>` to exclude (may be repeated)
* `<input>`: directory containing input HTML and Markdown documents (may be repeated)

See also
--------

Electrostatic is part of the [Mergician](https://github.com/rcrowley/mergician) suite of tools that manipulate HTML documents:

* [Deadlinks](https://github.com/rcrowley/deadlinks): Scan a document root directory for dead links
* [Feed](https://github.com/rcrowley/feed): Scan a document root directory to construct an Atom feed
* [Frag](https://github.com/rcrowley/frag): Extract fragments of HTML documents
* [Sitesearch](https://github.com/rcrowley/sitesearch): Index a document root directory and serve queries to it in AWS Lambda

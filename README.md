Electrostatic
=============

HTML is the language of the Web but the Web has moved beyond the bare, unstyled Web of the '90s (and university professors). Visitors to your site expect coherent layouts, visual consistency, navigation. This is the point at which most folks reach for a templating language. [Mergician](https://github.com/rcrowley/mergician) reimagines the often frustrating juxtaposition of HTML and templating language in pure HTML. Instead of rendering a template, merge a content HTML document into a layout HTML document.

Electostatic extends Mergician just that little bit further into a very basic CMS. It processes a whole document root directory, wrapping every HTML document in a consistent layout. Soon, it'll learn to generate reverse-chronological index pages and Atom feeds. Maybe it'll even learn how to manage email subscribers. You might want to test your inputs with [Deadlinks](https://github.com/rcrowley/deadlinks) before running Electrostatic or use [Sitesearch](https://github.com/rcrowley/sitesearch) to offer...site...search on Electrostatic's output.

The easiest way to use Electrostatic is with a `make`(1) target about like this:

```make
all:
	electrostatic -i raw -l design-system/index.html -o . -v
```

Installation
------------

```sh
go get github.com/rcrowley/electrostatic
```

Usage
-----

```sh
electrostatic -i <input> -l <layout> -o <output> [-p] [-v]
```

* `-i <input>`: directory containing input HTML and Markdown documents
* `-l <layout>`: site layout HTML document
* `-o <output>`: document root directory where merged HTML documents will be placed
* `-p`: pretend to process all the inputs but don't write any outputs; implies -v
* `-v`: verbose mode

See also
--------

Electrostatic is part of the [Mergician](https://github.com/rcrowley/mergician) suite ot tools that manipulate HTML documents:

* [Deadlinks](https://github.com/rcrowley/deadlinks): Scan a document root directory for dead links
* [Frag](https://github.com/rcrowley/frag): Extract fragments of HTML documents
* [Sitesearch](https://github.com/rcrowley/sitesearch): Index a document root directory and serve queries to it in AWS Lambda

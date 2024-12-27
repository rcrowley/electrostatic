Electrostatic
=============

HTML is the language of the Web but the Web has moved beyond the bare, unstyled Web of the '90s (and university professors). Visitors to your site expect coherent layouts, visual consistency, navigation. This is the point at which most folks reach for a templating language. [Mergician](https://github.com/rcrowley/mergician) reimagines the often frustrating juxtaposition of HTML and templating language in pure HTML. Instead of rendering a template, merge a content HTML document into a layout HTML document.

Electostatic extends Mergician just that little bit further into a very basic CMS. It processes a whole document root directory, wrapping every HTML document in a consistent layout. It generates an Atom feed from the most recent `<article>`s that contain a `<time>`. Soon, it'll learn to generate reverse-chronological index pages, too. Maybe it'll even learn how to manage email subscribers. You might want to test your inputs with [Deadlinks](https://github.com/rcrowley/deadlinks) before running Electrostatic or use [Sitesearch](https://github.com/rcrowley/sitesearch) to offer ... site ... search on Electrostatic's output.

The easiest way to use Electrostatic is with a `make`(1) target about like this:

```make
all:
	electrostatic -i raw -l design-system/index.html -o . -v
```

Installation
------------

```sh
go install github.com/rcrowley/electrostatic@latest
```

Usage
-----

```sh
electrostatic [-a <author>] -i <input> -l <layout> -o <output> [-p] [-r <rule>[...]] [-u <url>] [-v]
```

* `-a <author>`: author's name (used to build an Atom feed)
* `-i <input>`: directory containing input HTML and Markdown documents
* `-l <layout>`: site layout HTML document
* `-o <output>`: document root directory where merged HTML documents will be placed
* `-p`: pretend to process all the inputs but don't write any outputs; implies `-v`
* `-u <url>`: site URL with scheme and domain (used to build an Atom feed)
* `-v`: verbose mode

See also
--------

Electrostatic is part of the [Mergician](https://github.com/rcrowley/mergician) suite of tools that manipulate HTML documents:

* [Deadlinks](https://github.com/rcrowley/deadlinks): Scan a document root directory for dead links
* [Feed](https://github.com/rcrowley/feed): Scan a document root directory to construct an Atom feed
* [Frag](https://github.com/rcrowley/frag): Extract fragments of HTML documents
* [Sitesearch](https://github.com/rcrowley/sitesearch): Index a document root directory and serve queries to it in AWS Lambda

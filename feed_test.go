package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/rcrowley/mergician/html"
)

func TestFeed(t *testing.T) {
	f := &Feed{
		Items: []Item{
			{
				Date: "2024-12-03 22:28:00",
				Path: "newest.html",

				Node: must2(html.ParseString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/>
<title>Newest article title — Site name</title>
</head>
<body>
<header><h1>Site name</h1></header>
<article class="body">
<time datetime="2024-12-03 22:28:00">2024-12-03 22:28:00</time>
<h1>Newest article title</h1>
<p>Neweset article body.</p>
</article>
</body>
</html>
`)),
			},
			{
				Date: "1970-01-01 00:00:00",
				Path: "oldest.html",
				Node: must2(html.ParseString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8"/>
<title>Oldest article title — Site name</title>
</head>
<body>
<header><h1>Site name</h1></header>
<article class="body">
<time datetime="1970-01-01 00:00:00">1970-01-01 00:00:00</time>
<h1>Oldest article title</h1>
<p>Oldest article body.</p>
</article>
</body>
</html>
`)),
			},
		},
		URL: "http://example.com",
		t:   time.Now(),
	}
	stdout := &bytes.Buffer{}
	if err := f.Render(stdout); err != nil {
		t.Fatal(err)
	}
	actual := stdout.String()
	expected := fmt.Sprintf(`<rss version="2.0">
	<channel>
		<title>Site name</title>
		<description></description>
		<link>http://example.com</link>
		<pubDate>%s</pubDate>
		<item>
			<title>Newest article title</title>
			<description>2024-12-03 22:28:00 Newest article title Neweset article body.</description>
			<link>http://example.com/newest.html</link>
			<pubDate>Tue, 03 Dec 2024 22:28:00 UTC</pubDate>
		</item>
		<item>
			<title>Oldest article title</title>
			<description>1970-01-01 00:00:00 Oldest article title Oldest article body.</description>
			<link>http://example.com/oldest.html</link>
			<pubDate>Thu, 01 Jan 1970 00:00:00 UTC</pubDate>
		</item>
	</channel>
</rss>
`, f.t.Format(time.RFC1123))
	if actual != expected {
		t.Fatalf("actual: %s != expected: %s", actual, expected)
	}
}

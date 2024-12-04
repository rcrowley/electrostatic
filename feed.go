package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/rcrowley/mergician/html"
	"golang.org/x/net/html/atom"
)

type Feed struct {
	Items []Item
	URL   string
	mu    sync.Mutex
	t     time.Time
}

func (f *Feed) Add(date, path string, n *html.Node) {
	f.mu.Lock()
	defer f.mu.Unlock()
	i := sort.Search(len(f.Items), func(i int) bool { return f.Items[i].Date < date })
	if i == len(f.Items) || f.Items[i].Path != path {
		f.Items = append(f.Items, Item{})
		copy(f.Items[i+1:], f.Items[i:])
		f.Items[i].Date = date
		f.Items[i].Path = path
		f.Items[i].Node = n
	}
}

func (f *Feed) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	t := time.Now()
	if !f.t.IsZero() {
		t = f.t
	}
	u, err := url.Parse(f.URL)
	if err != nil {
		return err
	}
	n, err := html.ParseFile("index.html") // TODO parameterize "index.html" or at least make it aware of what directory it should be in
	if err != nil {
		return err
	}
	title := html.Find(n, html.IsAtom(atom.Title))
	if title == nil {
		return fmt.Errorf("no <title> in %s", "index.html")
	}

	// TODO <?xml version="1.0" encoding="UTF-8"?>
	e.EncodeToken(xml.StartElement{xml.Name{Local: "rss"}, []xml.Attr{{xml.Name{Local: "version"}, "2.0"}}})
	e.EncodeToken(xml.StartElement{xml.Name{Local: "channel"}, nil})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "title"}, nil})
	e.EncodeToken(xml.CharData(html.Text(title).String()))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "title"}})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "description"}, nil})
	e.EncodeToken(xml.CharData(""))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "description"}})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "link"}, nil})
	e.EncodeToken(xml.CharData(u.String()))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "link"}})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "pubDate"}, nil})
	e.EncodeToken(xml.CharData(t.Format(time.RFC1123)))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "pubDate"}})

	for _, item := range f.Items {
		article := html.Find(item.Node, html.IsAtom(atom.Article))
		if article == nil {
			return fmt.Errorf("no <article> in %s", item.Path)
		}
		h1 := html.Find(article, html.IsAtom(atom.H1))
		if h1 == nil {
			return fmt.Errorf("no <h1> in %s", item.Path)
		}

		e.EncodeToken(xml.StartElement{xml.Name{Local: "item"}, nil})

		e.EncodeToken(xml.StartElement{xml.Name{Local: "title"}, nil})
		e.EncodeToken(xml.CharData(html.Text(h1).String()))
		e.EncodeToken(xml.EndElement{xml.Name{Local: "title"}})

		e.EncodeToken(xml.StartElement{xml.Name{Local: "description"}, nil})
		e.EncodeToken(xml.CharData(html.Text(article).String()))
		e.EncodeToken(xml.EndElement{xml.Name{Local: "description"}})

		u.Path = item.Path
		e.EncodeToken(xml.StartElement{xml.Name{Local: "link"}, nil})
		e.EncodeToken(xml.CharData(u.String()))
		e.EncodeToken(xml.EndElement{xml.Name{Local: "link"}})

		e.EncodeToken(xml.StartElement{xml.Name{Local: "pubDate"}, nil})
		if t, err := time.Parse(time.DateTime, item.Date); err == nil {
			e.EncodeToken(xml.CharData(t.Format(time.RFC1123)))
		} else if t, err := time.Parse(time.DateOnly, item.Date); err == nil {
			e.EncodeToken(xml.CharData(t.Format(time.RFC1123)))
		} else {
			log.Printf("error parsing date %q: %v", item.Date, err)
			e.EncodeToken(xml.CharData(item.Date))
		}
		e.EncodeToken(xml.EndElement{xml.Name{Local: "pubDate"}})

		e.EncodeToken(xml.EndElement{xml.Name{Local: "item"}})
	}
	e.EncodeToken(xml.EndElement{xml.Name{Local: "channel"}})
	e.EncodeToken(xml.EndElement{xml.Name{Local: "rss"}})
	e.EncodeToken(xml.CharData("\n"))

	e.Flush()
	return nil
}

func (f *Feed) Print() {
	must(f.Render(os.Stdout))
}

func (f *Feed) Render(w io.Writer) error {
	if _, err := w.Write([]byte(xml.Header)); err != nil {
		return err
	}
	e := xml.NewEncoder(w)
	e.Indent("", "\t")
	return e.Encode(f)
}

type Item struct {
	Date, Path string
	Node       *html.Node
}

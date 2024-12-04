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
	Author string // author name
	Path   string // feed path within site URL, like "index.atom.xml"
	Title  string // site title
	URL    string // site URL, like "http://example.com"

	Entries []Entry

	mu sync.Mutex
	t  time.Time
}

func (f *Feed) Add(date, path string, n *html.Node) {
	f.mu.Lock()
	defer f.mu.Unlock()
	i := sort.Search(len(f.Entries), func(i int) bool { return f.Entries[i].Date < date })
	if i == len(f.Entries) || f.Entries[i].Path != path {
		f.Entries = append(f.Entries, Entry{})
		copy(f.Entries[i+1:], f.Entries[i:])
		f.Entries[i].Date = date
		f.Entries[i].Path = path
		f.Entries[i].Node = n
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

	e.EncodeToken(xml.Header)
	e.EncodeToken(xml.StartElement{xml.Name{Local: "feed"}, []xml.Attr{{xml.Name{Local: "xmlns"}, "http://www.w3.org/2005/Atom"}}})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "author"}, nil})
	e.EncodeToken(xml.StartElement{xml.Name{Local: "name"}, nil})
	e.EncodeToken(xml.CharData(f.Author))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "name"}})
	e.EncodeToken(xml.EndElement{xml.Name{Local: "author"}})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "id"}, nil})
	e.EncodeToken(xml.CharData(u.String()))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "id"}})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "link"}, []xml.Attr{
		{xml.Name{Local: "href"}, u.String()},
		{xml.Name{Local: "rel"}, "alternate"},
	}})
	e.EncodeToken(xml.EndElement{xml.Name{Local: "link"}}) // encoding/xml doesn't support self-closing tags
	u.Path = f.Path
	e.EncodeToken(xml.StartElement{xml.Name{Local: "link"}, []xml.Attr{
		{xml.Name{Local: "href"}, u.String()},
		{xml.Name{Local: "rel"}, "self"},
	}})
	e.EncodeToken(xml.EndElement{xml.Name{Local: "link"}}) // encoding/xml doesn't support self-closing tags

	e.EncodeToken(xml.StartElement{xml.Name{Local: "title"}, nil})
	e.EncodeToken(xml.CharData(f.Title))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "title"}})

	e.EncodeToken(xml.StartElement{xml.Name{Local: "updated"}, nil})
	e.EncodeToken(xml.CharData(t.Format(time.RFC3339)))
	e.EncodeToken(xml.EndElement{xml.Name{Local: "updated"}})

	for _, entry := range f.Entries {
		article := html.Find(entry.Node, html.IsAtom(atom.Article))
		if article == nil {
			return fmt.Errorf("no <article> in %s", entry.Path)
		}
		h1 := html.Find(article, html.IsAtom(atom.H1))
		if h1 == nil {
			return fmt.Errorf("no <h1> in %s", entry.Path)
		}

		e.EncodeToken(xml.StartElement{xml.Name{Local: "entry"}, nil})

		u.Path = entry.Path
		e.EncodeToken(xml.StartElement{xml.Name{Local: "id"}, nil})
		e.EncodeToken(xml.CharData(u.String()))
		e.EncodeToken(xml.EndElement{xml.Name{Local: "id"}})
		e.EncodeToken(xml.StartElement{xml.Name{Local: "link"}, []xml.Attr{
			{xml.Name{Local: "href"}, u.String()},
			{xml.Name{Local: "rel"}, "alternate"},
		}})
		e.EncodeToken(xml.EndElement{xml.Name{Local: "link"}}) // encoding/xml doesn't support self-closing tags

		e.EncodeToken(xml.StartElement{xml.Name{Local: "title"}, nil})
		e.EncodeToken(xml.CharData(html.Text(h1).String()))
		e.EncodeToken(xml.EndElement{xml.Name{Local: "title"}})

		e.EncodeToken(xml.StartElement{xml.Name{Local: "updated"}, nil})
		if t, err := time.Parse(time.DateTime, entry.Date); err == nil {
			e.EncodeToken(xml.CharData(t.Format(time.RFC3339)))
		} else if t, err := time.Parse(time.DateOnly, entry.Date); err == nil {
			e.EncodeToken(xml.CharData(t.Format(time.RFC3339)))
		} else {
			log.Printf("error parsing date %q: %v", entry.Date, err)
			e.EncodeToken(xml.CharData(entry.Date))
		}
		e.EncodeToken(xml.EndElement{xml.Name{Local: "updated"}})

		e.EncodeToken(xml.StartElement{xml.Name{Local: "content"}, []xml.Attr{{xml.Name{Local: "type"}, "html"}}})
		e.EncodeToken(xml.CharData(html.String(article)))
		e.EncodeToken(xml.EndElement{xml.Name{Local: "content"}})

		e.EncodeToken(xml.EndElement{xml.Name{Local: "entry"}})
	}
	e.EncodeToken(xml.EndElement{xml.Name{Local: "feed"}})
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

type Entry struct {
	Date, Path string
	Node       *html.Node
}

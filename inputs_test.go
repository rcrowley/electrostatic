package main

import "testing"

func TestInputsDuplicate(t *testing.T) {
	inputs := &Inputs{}
	inputs.Add("test.html")
	inputs.Add("test.html")
	pathnames := inputs.Pathnames()
	if len(pathnames) != 1 || pathnames[0] != "test.html" {
		t.Fatal(pathnames)
	}
}

func TestInputsHTMLThenMarkdown(t *testing.T) {
	inputs := &Inputs{}
	inputs.Add("test.html")
	inputs.Add("test.md")
	pathnames := inputs.Pathnames()
	if len(pathnames) != 1 || pathnames[0] != "test.md" {
		t.Fatal(pathnames)
	}
}

func TestInputsMarkdownThenHTML(t *testing.T) {
	inputs := &Inputs{}
	inputs.Add("test.md")
	inputs.Add("test.html")
	pathnames := inputs.Pathnames()
	if len(pathnames) != 1 || pathnames[0] != "test.md" {
		t.Fatal(pathnames)
	}
}

func TestInputsNotMarkdownOrHTML(t *testing.T) {
	inputs := &Inputs{}
	inputs.Add("test.go")
	pathnames := inputs.Pathnames()
	if len(pathnames) != 0 {
		t.Fatal(pathnames)
	}
}

func TestInputsSorted(t *testing.T) {
	inputs := &Inputs{}
	inputs.Add("b.html")
	inputs.Add("a.html")
	inputs.Add("d.html")
	inputs.Add("c.html")
	pathnames := inputs.Pathnames()
	if len(pathnames) != 4 || pathnames[0] != "a.html" || pathnames[1] != "b.html" || pathnames[2] != "c.html" || pathnames[3] != "d.html" {
		t.Fatal(pathnames)
	}
}

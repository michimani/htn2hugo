package main

import (
	"fmt"
)

type Atom struct {
	Title    string  `xml:"title"`
	SubTitle string  `xml:"subtitle"`
	Links    []Link  `xml:"link"`
	Entries  []Entry `xml:"entry"`
	NextPage string
	HasNext  bool
	Host     string
}

func (a *Atom) prepare() error {
	for _, link := range a.Links {
		if link.Rel == "next" {
			a.NextPage = link.Href
		}
		if link.Rel == "alternate" {
			a.Host = link.Href
		}
	}

	if a.Host == "" {
		return fmt.Errorf("Cannot get host URL.")
	}

	a.HasNext = (a.NextPage != "")
	return nil
}

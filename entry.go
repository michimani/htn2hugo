package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Entry struct {
	Links          []Link     `xml:"link"`
	Author         string     `xml:"author>name"`
	Title          string     `xml:"title"`
	Published      string     `xml:"updated"`
	Content        string     `xml:"content"`
	Summary        string     `xml:"summary"`
	Draft          string     `xml:"control>draft"`
	Categories     []Category `xml:"category"`
	CategoriesStr  string
	PublishedYear  string
	PublishedMonth string
	Url            string
	Permalink      string
	FileName       string
}

type Category struct {
	Term string `xml:"term,attr"`
}

const (
	publiedFmt string = "2006-01-02T15:04:05-07:00"
	hugoTmp    string = `---
title: "%s"
date: %s
draft: %s
author: ["%s"]
categories: [%s]
archives: ["%s", "%s"]
description: "%s"
url: "/%s"
---

%s`
	simpleLink string = "\n<a href=\"$1\">$1</a>  \n"
	embedLink  string = `
<div class="inner-link-wrapper">
  <iframe
    class="hatenablogcard"
    style="width:100%;height:155px;max-width:680px;"
    src="https://hatenablog-parts.com/embed?url=$1"
    width="300" height="150" frameborder="0" scrolling="no">
  </iframe>
</div>
`
	photolifeJPEGLink string = `
<a href="https://f.hatena.ne.jp/$1/$2$3">
	<img src="https://cdn-ak.f.st-hatena.com/images/fotolife/m/$1/$2/$2$3.jpg" alt="$2$3">
</a>`
	photolifePNGLink string = `
<a href="https://f.hatena.ne.jp/$1/$2$3">
	<img src="https://cdn-ak.f.st-hatena.com/images/fotolife/m/$1/$2/$2$3.png" alt="$2$3">
</a>
`
	twitterEmbedLink string = `
<blockquote class="twitter-tweet" >
	<p lang="ja" dir="ltr"></p>
	<a href="https://twitter.com/$1"></a>
</blockquote>
<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
`
)

func (e *Entry) save() bool {
	err := e.prepare()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	file, err := os.Create(saveDir + e.FileName)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer file.Close()

	content := fmt.Sprintf(hugoTmp, e.Title, e.Published, e.Draft, e.Author, e.CategoriesStr, e.PublishedYear, (e.PublishedYear + "-" + e.PublishedMonth), e.Summary, e.Permalink, e.Content)

	writer := bufio.NewWriter(file)
	if _, err := writer.Write([]byte(content)); err != nil {
		fmt.Println(err.Error())
		return false
	}

	writer.Flush()
	return true
}

func (e *Entry) prepare() error {
	for _, link := range e.Links {
		if link.Rel == "alternate" {
			e.Url = link.Href
			break
		}
	}

	e.Draft = strconv.FormatBool(e.Draft == "yes")
	e.Permalink = strings.Replace(e.Url, host, "", 1)
	e.FileName = strings.ReplaceAll(e.Permalink, "/", "_") + ".md"
	for _, c := range e.Categories {
		if c.Term == "" {
			continue
		}
		if e.CategoriesStr != "" {
			e.CategoriesStr = e.CategoriesStr + ","
		}
		e.CategoriesStr = fmt.Sprintf("%s\"%s\"", e.CategoriesStr, c.Term)
	}
	published, err := time.Parse(publiedFmt, e.Published)
	if err != nil {
		return err
	}
	e.PublishedYear = published.Format("2006")
	e.PublishedMonth = published.Format("01")

	e.replaceHatenaSyntax()

	return nil
}

func (e *Entry) replaceHatenaSyntax() {
	e.replaceSimpleLink()
	e.replaceEmbededLink()
	e.replacePhotolifeLink()
	e.replaceTwitterEmbedLink()
	e.removeUnnecessary()
	e.escapeMetadata()
}

func (e *Entry) replaceSimpleLink() {
	reg := regexp.MustCompile(`\[(.*?):title\]`)
	e.Content = reg.ReplaceAllString(e.Content, simpleLink)
}

func (e *Entry) replaceEmbededLink() {
	reg := regexp.MustCompile(`\[(.*?):embed:cite\]`)
	e.Content = reg.ReplaceAllString(e.Content, embedLink)
}

func (e *Entry) replacePhotolifeLink() {
	reg := regexp.MustCompile(`\[f:id:(.*?):(\d{8}?)(\d+)j:.*\]`)
	e.Content = reg.ReplaceAllString(e.Content, photolifeJPEGLink)

	reg = regexp.MustCompile(`\[f:id:(.*?):(\d{8}?)(\d+)p:.*\]`)
	e.Content = reg.ReplaceAllString(e.Content, photolifePNGLink)
}

func (e *Entry) replaceTwitterEmbedLink() {
	reg := regexp.MustCompile(`\[https:\/\/twitter.com\/(.*?):embed\]`)
	e.Content = reg.ReplaceAllString(e.Content, twitterEmbedLink)
}

func (e *Entry) removeUnnecessary() {
	e.Content = strings.ReplaceAll(e.Content, "<!-- more -->", "")
	e.Content = strings.ReplaceAll(e.Content, "[:contents]", "")
}

func (e *Entry) escapeMetadata() {
	e.Title = strings.ReplaceAll(e.Title, "\"", `\"`)
	e.Summary = strings.ReplaceAll(e.Summary, "\"", `\"`)
	e.Author = strings.ReplaceAll(e.Author, "\"", `\"`)
}

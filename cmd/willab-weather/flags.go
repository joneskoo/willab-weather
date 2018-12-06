package main

import (
	"net/url"
	"text/template"
)

type urlFlag struct{ *url.URL }

func (uf *urlFlag) Set(s string) (err error) {
	uf.URL, err = url.Parse(s)
	return err
}

func (uf urlFlag) String() string {
	if uf.URL == nil {
		return "nil"
	}
	return uf.URL.String()
}

type templateFlag struct {
	TemplateText string
	*template.Template
}

func (tf *templateFlag) Set(s string) (err error) {
	t, err := template.New("weather").Parse(s)
	tf.Template = t
	tf.TemplateText = s
	return err
}

func (tf *templateFlag) String() string {
	return ""
}

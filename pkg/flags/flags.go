package flags

import (
	"net/url"
	"text/template"
)

type URLFlag struct{ *url.URL }

func (uf *URLFlag) Set(s string) (err error) {
	uf.URL, err = url.Parse(s)
	return err
}

func (uf URLFlag) String() string {
	if uf.URL == nil {
		return "nil"
	}
	return uf.URL.String()
}

type TemplateFlag struct {
	TemplateText string
	*template.Template
}

func (tf *TemplateFlag) Set(s string) (err error) {
	t, err := template.New("weather").Parse(s)
	tf.Template = t
	tf.TemplateText = s
	return err
}

func (tf *TemplateFlag) String() string {
	return ""
}

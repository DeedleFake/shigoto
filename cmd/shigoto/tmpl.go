package main

import (
	"strings"
	"text/template"

	"github.com/DeedleFake/shigoto"
)

var defaults = map[string]interface{}{
	"sourceName": `{{.Title | slug}}.md`,
	"buildPath":  `{{.Title | slug}}/index.{{.Type | ext}}`,
	"pages":      pagesInfo{Per: 5},
}

func tmplGet(name string, meta ...map[string]interface{}) interface{} {
	for _, meta := range meta {
		v, ok := meta[name]
		if ok {
			if fr, ok := defaults[name].(fromRawer); ok {
				return fr.fromRaw(v)
			}
			return v
		}
	}

	return defaults[name]
}

func metaTmpl(src string, data interface{}) (string, error) {
	snt, err := template.New(src).Funcs(shigoto.StandardFuncs(nil)).Parse(src)
	if err != nil {
		return "", err
	}

	var r strings.Builder
	err = snt.Execute(&r, data)
	return r.String(), err
}

type fromRawer interface {
	fromRaw(raw interface{}) interface{}
}

type pagesInfo struct {
	Tmpl string `yaml:"tmpl"`
	Per  int    `yaml:"per"`
}

func (info pagesInfo) fromRaw(raw interface{}) interface{} {
	rawmap, ok := raw.(map[interface{}]interface{})
	if !ok {
		return nil
	}

	if tmpl, ok := rawmap["tmpl"].(string); ok {
		info.Tmpl = tmpl
	}
	if per, ok := rawmap["per"].(int); ok {
		info.Per = per
	}

	return info
}

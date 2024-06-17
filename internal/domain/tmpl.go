package domain

import (
	"html/template"

	"github.com/catalystgo/logger/log"
)

// ParseTemplate parses a template and returns an error if one occurs.
func ParseTemplate(name string, content []byte) (*template.Template, error) {
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		log.Debugf("parse template %s: %v", name, err)
		return nil, err
	}
	return tmpl, nil
}

// MustParseTemplate parses a template and panics if an error occurs.
func MustParseTemplate(name string, content []byte) *template.Template {
	tmpl, err := ParseTemplate(name, content)
	if err != nil {
		log.Panicf("must parse template %s: %v", name, err)
	}
	return tmpl
}

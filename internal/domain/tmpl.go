package domain

import (
	"bytes"
	"text/template"

	log "github.com/catalystgo/logger/cli"
)

// ParseTemplate parses a template and returns an error if one occurs.
func ParseTemplate(name string, content []byte) (*template.Template, error) {
	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
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

// ExecuteTemplate executes a template with given data and returns an error if one occurs.
func ExecuteTemplate(t *template.Template, data interface{}) ([]byte, error) {
	var b bytes.Buffer
	err := t.Execute(&b, data)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

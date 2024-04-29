package service

import (
	"os"
	"path/filepath"

	"github.com/catalystgo/cli/internal/component"
	"github.com/catalystgo/cli/internal/log"
)

var _ Service = service{}

type Service interface {
	Init(module string, components []component.Component, override bool)
}

type service struct{}

// Init implements CommandService.
func (c service) Init(module string, components []component.Component, override bool) {
	for _, c := range components {
		name := filepath.Join(c.Path(), c.Name())

		log.Debugf("file = %s | exist = %t | override = %t", name, checkFileExists(name), override)

		// Override files content only if override flag is passed
		if checkFileExists(name) {
			if override {
				log.Warnf("overriding file (%s)", name)
			} else {
				log.Warnf("skipping file (%s) => already exist", name)
				continue
			}
		}

		// Get content
		b, err := c.Content()
		if err != nil {
			log.Errorf("read content file (%s) => %v", name, err)
			continue
		}

		// Create sub-directories if needed
		err = os.MkdirAll(c.Path(), os.ModePerm)
		if err != nil {
			log.Errorf("mkdir file (%s) => %v", name, err)
			continue
		}

		// Create file
		f, err := os.Create(name)
		if err != nil {
			log.Errorf("create file (%s) => %v", name, err)
			continue
		}

		// Write content
		_, err = f.Write(b)
		if err != nil {
			log.Errorf("write content file (%s) => %v", name, err)
			continue
		}

		log.Infof("successfully createad %s", name)
	}
}

func New() Service {
	return service{}
}

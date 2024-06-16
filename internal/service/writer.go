package service

import (
	"os"
	"path/filepath"

	"github.com/catalystgo/cli/internal/log"
)

type WriteOpt struct {
	Override    bool // If true, existing files will be overridden. Default is false.
	IgnoreEmpty bool // If true, empty files will not be written. Default is true.
}

type WriteOption func(*WriteOpt)

func WithOverride(b bool) WriteOption {
	return func(w *WriteOpt) {
		w.Override = b
	}
}

func WithIgnoreEmpty(b bool) WriteOption {
	return func(w *WriteOpt) {
		w.IgnoreEmpty = b
	}
}

func parseWriteOpts(opts ...WriteOption) *WriteOpt {
	cfg := &WriteOpt{
		Override:    false,
		IgnoreEmpty: true,
	}

	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

type Writer struct{}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) WriteFile(file string, data []byte, opts ...WriteOption) error {
	cfg := parseWriteOpts(opts...)
	if cfg.IgnoreEmpty && len(data) == 0 {
		log.Warnf("no data to write in file (%s) therefore skipping", file)
		return nil
	}

	err := os.MkdirAll(filepath.Dir(file), os.ModePerm)
	if err != nil {
		log.Errorf("mkdir file (%s) => %v", file, err)
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		log.Errorf("create file (%s) => %v", file, err)
		return err
	}

	defer func() {
		err := f.Close()
		if err != nil {
			log.Errorf("close file (%s) => %v", file, err)
		}
	}()

	_, err = f.Write(data)
	if err != nil {
		log.Errorf("write file (%s) => %v", file, err)
		return err
	}

	log.Infof("successfully createad %s", file)

	return nil
}

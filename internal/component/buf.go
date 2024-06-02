package component

import (
	_ "embed"
)

var (
	//go:embed template/buf/buf.yaml
	bufContent []byte

	//go:embed template/buf/buf.gen.yaml
	bufGenContent []byte

	//go:embed template/buf/buf.gen.vendor.yaml
	bufGenVendorContent []byte
)

type (
	bufComponent          struct{}
	bufGenComponent       struct{}
	bugGenVendorComponent struct{}
)

//////////////////
// BUF COMPONENT
//////////////////

func NewBufComponent() Component {
	return bufComponent{}
}

func (f bufComponent) Content() ([]byte, error) {
	return bufContent, nil
}

func (f bufComponent) Name() string {
	return "buf.yaml"
}

func (f bufComponent) Path() string {
	return "."
}

//////////////////
// BUF GEN COMPONENT
//////////////////

func NewBufGenComponent() Component {
	return bufGenComponent{}
}

func (f bufGenComponent) Content() ([]byte, error) {
	return bufGenContent, nil
}

func (f bufGenComponent) Name() string {
	return "buf.gen.yaml"
}

func (f bufGenComponent) Path() string {
	return "."
}

//////////////////
// BUF GEN VENDOR COMPONENT
//////////////////

func NewBufGenVendorComponent() Component {
	return bugGenVendorComponent{}
}

func (f bugGenVendorComponent) Content() ([]byte, error) {
	return bufGenVendorContent, nil
}

func (f bugGenVendorComponent) Name() string {
	return "buf.gen.vendor.yaml"
}

func (f bugGenVendorComponent) Path() string {
	return "."
}

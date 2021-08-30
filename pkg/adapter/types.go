package adapter

import (
	"io"

	"github.com/maxott/magda-cli/pkg/log"
)

type Adapter interface {
	Get(path string, logger log.Logger) (Payload, error)
	Post(path string, body io.Reader, logger log.Logger) (Payload, error)
	Put(path string, body io.Reader, logger log.Logger) (Payload, error)
	Patch(path string, body io.Reader, logger log.Logger) (Payload, error)
	Delete(path string, logger log.Logger) (Payload, error)

	SkipGateway() bool // experimental!
}

type Payload interface {
	IsObject() bool
	AsObject() map[string]interface{}
	AsArray() []interface{}
	AsBytes() []byte
}

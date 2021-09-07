package adapter

import (
	"io"

	log "go.uber.org/zap"
)

type Adapter interface {
	Get(path string, logger *log.Logger) (Payload, error)
	Post(path string, body io.Reader, logger *log.Logger) (Payload, error)
	Put(path string, body io.Reader, logger *log.Logger) (Payload, error)
	Patch(path string, body io.Reader, logger *log.Logger) (Payload, error)
	Delete(path string, logger *log.Logger) (Payload, error)

	SkipGateway() bool // experimental!
}

type Payload interface {
	// IsObject() bool
	AsType(r interface{}) error
	AsObject() (map[string]interface{}, error)
	AsArray() ([]interface{}, error)
	AsBytes() []byte
}

package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/maxott/magda-cli/pkg/log"
)

func ReplyPrinter(pld Payload, err error) error {
	if err != nil {
		return err
	}

	var b []byte
	var err2 error
	if pld.IsObject() {
		b, err2 = json.MarshalIndent(pld.AsObject(), "", "  ")
	}
	if !pld.IsObject() {
		b, err2 = json.MarshalIndent(pld.AsArray(), "", "  ")
	}
	if err2 != nil {
		return err2
	} else {
		fmt.Printf("%s\n", b)
		return nil
	}
}

func ToPayload(body []byte, contentType string, logger log.Logger) (Payload, error) {
	logger.Debugf("Received content-type '%s'", contentType)	
	var f interface{}
	err := json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}

	switch m := f.(type) {
	case []interface{}:
		return JsonArrPayload{m, body}, nil
	case map[string]interface{}:
		return JsonObjPayload{m, body}, nil
	default:
		return nil, fmt.Errorf("Unknown json type in body")
	}
}

type JsonObjPayload struct {
	payload map[string]interface{}
	bytes   []byte
}

func (p JsonObjPayload) IsObject() bool                   { return true }
func (p JsonObjPayload) AsObject() map[string]interface{} { return p.payload }
func (p JsonObjPayload) AsArray() []interface{}           { return []interface{}{p.payload} }
func (p JsonObjPayload) AsBytes() []byte                  { return p.bytes }

type JsonArrPayload struct {
	payload []interface{}
	bytes   []byte
}

func (JsonArrPayload) IsObject() bool                     { return false }
func (p JsonArrPayload) AsObject() map[string]interface{} { return nil }
func (p JsonArrPayload) AsArray() []interface{}           { return p.payload }
func (p JsonArrPayload) AsBytes() []byte                  { return p.bytes }

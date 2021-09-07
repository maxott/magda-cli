package adapter

import (
	"encoding/json"
	"fmt"
	"errors"
	"io/ioutil"

	log "go.uber.org/zap"
)

type payload struct {
	contentType string
	body   []byte
}

func ToPayload(body []byte, contentType string, logger *log.Logger) (Payload, error) {
	logger.Debug("Received", log.String("content-type", contentType))
	return &payload{body: body, contentType: contentType}, nil
	// var f interface{}
	// err := json.Unmarshal(body, &f)
	// if err != nil {
	// 	return nil, err
	// }

	// switch m := f.(type) {
	// case []interface{}:
	// 	return JsonArrPayload{m, body}, nil
	// case map[string]interface{}:
	// 	return JsonObjPayload{m, body}, nil
	// default:
	// 	return nil, logger.Error(nil, "Unknown json type in body")
	// }
}

func LoadPayloadFromFile(fileName string) (Payload, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return &payload{body: data}, nil
	// var f interface{}
	// err = json.Unmarshal(data, &f)
	// if err != nil {
	// 	return nil, err
	// }
	// m := f.(map[string]interface{})
	// return JsonObjPayload{m, data}, nil
}

func ReplyPrinter(pld Payload, err error) error {
	if err != nil {
		return err
	}

	var f interface{}
	if err = pld.AsType(&f); err != nil {
		return err
	}
	var b []byte
	if b, err = json.MarshalIndent(f, "", "  "); err != nil {
		return err
	} else {
		fmt.Printf("%s\n", b)
		return nil
	}
}

func (p *payload) AsType(r *interface{}) error {
	return json.Unmarshal(p.body, r)
}

func (p *payload) AsObject() (map[string]interface{}, error) {
	var f interface{}
	err := json.Unmarshal(p.body, &f)
	if err != nil {
		return nil, err
	}
	if obj, ok := f.(map[string]interface{}); ok {
		return obj, nil
	} else {
		return nil, errors.New("not an object type")
	}
}

func (p *payload) AsArray() ([]interface{}, error) {
	var f interface{}
	err := json.Unmarshal(p.body, &f)
	if err != nil {
		return nil, err
	}
	switch m := f.(type) {
	case []interface{}:
		return m, nil
	case map[string]interface{}:
		return[]interface{}{m}, nil
	default:
		return nil, errors.New("not an array type")
	}
}

func (p *payload) AsBytes() []byte {
	return p.body
}

// type JsonObjPayload struct {
// 	payload map[string]interface{}
// 	bytes   []byte
// }

// func (p JsonObjPayload) IsObject() bool                   { return true }
// func (p JsonObjPayload) AsObject() map[string]interface{} { return p.payload }
// func (p JsonObjPayload) AsArray() []interface{}           { return []interface{}{p.payload} }
// func (p JsonObjPayload) AsBytes() []byte                  { return p.bytes }

// type JsonArrPayload struct {
// 	payload []interface{}
// 	bytes   []byte
// }

// func (JsonArrPayload) IsObject() bool                     { return false }
// func (p JsonArrPayload) AsObject() map[string]interface{} { return nil }
// func (p JsonArrPayload) AsArray() []interface{}           { return p.payload }
// func (p JsonArrPayload) AsBytes() []byte                  { return p.bytes }

package schema

import (
	"bytes"
	"encoding/json"

	"github.com/maxott/magda-cli/pkg/adapter"
	log "go.uber.org/zap"
)

/**** LIST ****/

type ListRequest struct {
}

func ListRaw(cmd *ListRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := aspectPath(nil, adpt)
	return (*adpt).Get(path, logger)
}

/**** CREATE ****/

type CreateRequest struct {
	Id     string                 `json:"id"`
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"jsonSchema"`
}

func CreateRaw(cmd *CreateRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	r := *cmd

	if body, err := json.MarshalIndent(r, "", "  "); err != nil {
		return nil, err
	} else {
		path := aspectPath(nil, adpt)
		return (*adpt).Post(path, bytes.NewReader(body), logger)
	}
}

/**** READ ****/

type ReadRequest struct {
	Id string
}

func ReadRaw(cmd *ReadRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := aspectPath(&cmd.Id, adpt)
	return (*adpt).Get(path, logger)
}

/**** UPDATE ****/

type UpdateRequest = CreateRequest

func UpdateRaw(cmd *UpdateRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	r := *cmd

	path := aspectPath(&r.Id, adpt)

	if r.Name == "" {
		// get current 'name' first as it is required
		if pld, err := (*adpt).Get(path, logger); err != nil {
			return nil, err
		} else {
			obj, err := pld.AsObject()
			if err != nil {
				logger.Error("no schema found in body", log.Error(err))
				return nil, err
			}
			r.Name = obj["name"].(string)
		}
	}
	if body, err := json.MarshalIndent(r, "", "  "); err != nil {
		return nil, err
	} else {
		// path := aspectPath(&r.Id, adpt)
		return (*adpt).Put(path, bytes.NewReader(body), logger)
	}
}

/**** DELETE ****/

// Not supported

/**** Utils ****/

func aspectPath(id *string, adpt *adapter.Adapter) string {
	path := "/api/v0/registry/aspects"
	if (*adpt).SkipGateway() {
		path = "/v0/aspects"
	}
	if id != nil {
		path = path + "/" + *id
	}
	return path
}

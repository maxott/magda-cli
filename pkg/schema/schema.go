package schema

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/maxott/magda-cli/pkg/adapter"
)

/**** LIST ****/

type ListCmd struct {
}

func ListRaw(cmd *ListCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
	path := aspectPath(nil, adpt)
	return (*adpt).Get(path)
}

/**** CREATE ****/

type CreateCmd struct {
	Id     string                 `json:"id"`
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"jsonSchema"`
}

func CreateRaw(cmd *CreateCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
	r := *cmd

	if body, err := json.MarshalIndent(r, "", "  "); err != nil {
		return nil, err
	} else {
		path := aspectPath(nil, adpt)
		return (*adpt).Post(path, bytes.NewReader(body))
	}
}

/**** READ ****/

type ReadCmd struct {
	Id string
}

func ReadRaw(cmd *ReadCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
	path := aspectPath(&cmd.Id, adpt)
	return (*adpt).Get(path)
}

/**** UPDATE ****/

type UpdateCmd = CreateCmd

func UpdateRaw(cmd *UpdateCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
	r := *cmd

	path := aspectPath(&r.Id, adpt)

	if r.Name == "" {
		// get current 'name' first as it is required
		if pld, err := (*adpt).Get(path); err != nil {
			return nil, err
		} else {
			obj := pld.AsObject()
			if obj == nil {
				return nil, fmt.Errorf("no schema body found")
			}
			r.Name = obj["name"].(string)
		}
	}
	if body, err := json.MarshalIndent(r, "", "  "); err != nil {
		return nil, err
	} else {
		// path := aspectPath(&r.Id, adpt)
		return (*adpt).Put(path, bytes.NewReader(body))
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

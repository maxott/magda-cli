package minion

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
	path := minionPath(nil, adpt)
	return (*adpt).Get(path, logger)
}

/**** CREATE ****/

// {
// 	id: options.id,
// 	name: options.id,
// 	active: true,
// 	url: getWebhookUrl(options),
// 	eventTypes: [
// 	"CreateRecord",
// 	"CreateAspectDefinition",
// 	"CreateRecordAspect",
// 	"PatchRecord",
// 	"PatchAspectDefinition",
// 	"PatchRecordAspect"
// 	],
// 	config: {
// 	aspects: options.aspects,
// 	optionalAspects: options.optionalAspects,
// 	includeEvents: false,
// 	includeRecords: true,
// 	includeAspectDefinitions: false,
// 	dereference: true
// 	}
// 	,
// 	lastEvent: null,
// 	isWaitingForResponse: false,
// 	enabled: true,
// 	lastRetryTime: null,
// 	retryCount: 0,
// 	isRunning: null,
// 	isProcessing: null
// 	}

type EventType string

const (
	CreateRecord           EventType = "CreateRecord"
	CreateAspectDefinition EventType = "CreateAspectDefinition"
	CreateRecordAspect     EventType = "CreateRecordAspect"
	PatchRecord            EventType = "PatchRecord"
	PatchAspectDefinition  EventType = "PatchAspectDefinition"
	PatchRecordAspect      EventType = "PatchRecordAspect"
)

// can't define a const array
var defEventTypes = []EventType{
	CreateRecord, CreateAspectDefinition, CreateRecordAspect,
	PatchRecord, PatchAspectDefinition, PatchRecordAspect,
}

type CreateRequest struct {
	Id              string
	Url             string
	EventTypes      []EventType
	Aspects         []string
	OptionalAspects []string
	RetryCount      int
}

type createPayload struct {
	Id         string       `json:"id"`
	Name       string       `json:"name"`
	Url        string       `json:"url"`
	Active     bool         `json:"active"`
	Enabled    bool         `json:"enabled"`
	EventTypes []EventType  `json:"eventTypes"`
	Config     createConfig `json:"config"`
	RetryCount int          `json:"retryCount"`
}

type createConfig struct {
	Aspects                  []string `json:"aspects"`
	OptionalAspects          []string `json:"optionalAspects"`
	IncludeEvents            bool     `json:"includeEvents"`
	IncludeRecords           bool     `json:"includeRecords"`
	IncludeAspectDefinitions bool     `json:"includeAspectDefinitions"`
	Dereference              bool     `json:"dereference"`
}

func CreateRaw(cmd *CreateRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	config := createConfig{
		Aspects:                  cmd.Aspects,
		OptionalAspects:          cmd.OptionalAspects,
		IncludeEvents:            false,
		IncludeRecords:           true,
		IncludeAspectDefinitions: false,
		Dereference:              true,
	}
	if config.Aspects == nil {
		config.Aspects = make([]string, 0)
	}
	if config.OptionalAspects == nil {
		config.OptionalAspects = make([]string, 0)
	}

	r := createPayload{
		Id: cmd.Id, Name: cmd.Id,
		Url:     cmd.Url,
		Enabled: true, Active: true,
		Config: config,
	}

	if r.EventTypes == nil {
		r.EventTypes = defEventTypes
	}
	if body, err := json.MarshalIndent(r, "", "  "); err != nil {
		return nil, err
	} else {
		path := minionPath(nil, adpt)
		logger.Info("POTS minion", log.ByteString("body", body))
		return (*adpt).Post(path, bytes.NewReader(body), logger)
	}
}

/**** READ ****/

// type ReadCmd struct {
// 	Id string
// }

// func ReadRaw(cmd *ReadCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
// 	path := minionPath(&cmd.Id, adpt)
// 	return (*adpt).Get(path)
// }

/**** UPDATE ****/

// Not supported

// type UpdateCmd = CreateCmd

// func UpdateRaw(cmd *UpdateCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
// 	r := *cmd

// 	path := minionPath(&r.Id, adpt)

// 	if r.Name == "" {
// 		// get current 'name' first as it is required
// 		if pld, err := (*adpt).Get(path); err != nil {
// 			return nil, err
// 		} else {
// 			obj := pld.AsObject()
// 			if obj == nil {
// 				return nil, fmt.Errorf("no minion body found")
// 			}
// 			r.Name = obj["name"].(string)
// 		}
// 	}
// 	if body, err := json.MarshalIndent(r, "", "  "); err != nil {
// 		return nil, err
// 	} else {
// 		path := minionPath(nil, adpt)
// 		return (*adpt).Put(path, bytes.NewReader(body))
// 	}
// }

/**** DELETE ****/

type DeleteRequest struct {
	Id string
}

func DeleteRaw(cmd *DeleteRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := minionPath(&cmd.Id, adpt)
	return (*adpt).Delete(path, logger)
}

/**** Utils ****/

func minionPath(id *string, adpt *adapter.Adapter) string {
	path := "/api/v0/registry/hooks"
	if (*adpt).SkipGateway() {
		path = "/v0/hooks"
	}
	if id != nil {
		path = path + "/" + *id
	}
	return path
}

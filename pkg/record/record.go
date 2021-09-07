package record

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/maxott/magda-cli/pkg/adapter"
	log "go.uber.org/zap"	
)

/**** LIST ****/

type ListRequest struct {
	Aspects   string
	Query     string
	Offset    int
	Limit     int
	PageToken string
}

type ListResult struct {
	HasMore       bool   `json:"hasMore"`
	NextPageToken string `json:"nextPageToken"`
	Records       []struct {
		Aspects   map[string]interface{} `json:"aspects"`
		ID        string                 `json:"id"`
		Name      string                 `json:"name"`
		SourceTag string                 `json:"sourceTag"`
		TenantID  int                    `json:"tenantId"`
	} `json:"records"`
}

func List(cmd *ListRequest, adpt *adapter.Adapter, logger *log.Logger) (ListResult, error) {
	pyl, err := ListRaw(cmd, adpt, logger)
	if err != nil {
		return ListResult{}, err
	}
	res := ListResult{}
	_ = json.Unmarshal(pyl.AsBytes(), &res)
	return res, nil
}

func ListRaw(cmd *ListRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := recordPath(nil, adpt)

	q := []string{}
	if cmd.Aspects != "" {
		q = append(q, "aspect="+url.QueryEscape(cmd.Aspects))
	}
	if cmd.Query != "" {
		q = append(q, "aspectQuery="+url.QueryEscape(cmd.Query))
	}
	if cmd.PageToken != "" {
		q = append(q, "pageToken="+url.QueryEscape(cmd.PageToken))
	}
	if cmd.Offset >= 0 {
		q = append(q, "start="+url.QueryEscape(strconv.Itoa(cmd.Offset)))
	}
	if cmd.Limit >= 0 {
		q = append(q, "limit="+url.QueryEscape(strconv.Itoa(cmd.Limit)))
	}
	if len(q) > 0 {
		path = path + "?" + strings.Join(q, "&")
	}
	//fmt.Printf("PATH: %s\n", path)
	return (*adpt).Get(path, logger)
}

/**** CREATE ****/

type CreateRequest struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Aspects   Aspects `json:"aspects"`
	SourceTag string  `json:"sourceTag,omitempty"`
}

type Aspects map[string]Aspect
type Aspect map[string]interface{}

func CreateRaw(cmd *CreateRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	if (*cmd).Id == "" {
		(*cmd).Id = uuid.New().String()
	}

	body, err := json.MarshalIndent(*cmd, "", "  ")
	if err != nil {
		logger.Error("error marshalling body.", log.Error(err))
		return nil, err
	}
  // fmt.Printf("RECORD %+v - %s\n", cmd, body)

	path := recordPath(nil, adpt)
	return (*adpt).Post(path, bytes.NewReader(body), logger)
}

/**** READ ****/

type ReadRequest struct {
	Id         string
	AddAspects string
	Aspect     string
}

func ReadRaw(cmd *ReadRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := recordPath(&cmd.Id, adpt)
	if cmd.AddAspects != "" {
		path = path + "?aspect=" + cmd.AddAspects
	} else if cmd.Aspect != "" {
		path = path + "/aspects/" + cmd.Aspect
	} else {
		// display summary
		path = recordPath(nil, adpt) + "/summary/" + cmd.Id
	}
	return (*adpt).Get(path, logger)
}

/**** UPDATE ****/

type UpdateRequest = CreateRequest

func UpdateRaw(cmd *UpdateRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	r := *cmd

	path := recordPath(&r.Id, adpt)
	if r.Name == "" {
		// get current 'name' first as it is required
		pld, err := (*adpt).Get(path, logger)
		if err != nil {
			return nil, err
		}
		obj, err := pld.AsObject()
		if err != nil {
			logger.Error("no record body found", log.Error(err))
			return nil, err
		}
		r.Name = obj["name"].(string)
	}
	body, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		logger.Error("error marshalling body.", log.Error(err))
		return nil, err
	}
	return (*adpt).Put(path, bytes.NewReader(body), logger)
}

/**** PATCH ASPECT ********/

type PatchAspectRequest struct {
	Id         string
	Aspect     string
	Patch      []interface{}
}

func PatchAspectRaw(cmd *PatchAspectRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := recordPath(&cmd.Id, adpt) + "/aspects/" + cmd.Aspect
	body, err := json.MarshalIndent(cmd.Patch, "", "  ")
	if err != nil {
		logger.Error("marshalling body", log.Error(err))
		return nil, err
	}
	return (*adpt).Patch(path, bytes.NewReader(body), logger)
}

/**** DELETE ****/

type DeleteRequest struct {
	Id         string
	AspectName string
}

func DeleteRaw(cmd *DeleteRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := recordPath(&cmd.Id, adpt)
	if cmd.AspectName != "" {
		path = path + "/aspects/" + cmd.AspectName
	}
	return (*adpt).Delete(path, logger)
}

/**** HISTORY ****/

type HistoryRequest struct {
	Id        string
	EventId   string
	Offset    int
	Limit     int
	PageToken string
}

func HistoryRaw(cmd *HistoryRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := recordPath(&cmd.Id, adpt) + "/history"
	if cmd.EventId != "" {
		path = path + "/" + cmd.EventId
	}

	q := []string{}
	if cmd.PageToken != "" {
		q = append(q, "pageToken="+url.QueryEscape(cmd.PageToken))
	}
	if cmd.Offset >= 0 {
		q = append(q, "start="+url.QueryEscape(strconv.Itoa(cmd.Offset)))
	}
	if cmd.Limit >= 0 {
		q = append(q, "limit="+url.QueryEscape(strconv.Itoa(cmd.Limit)))
	}
	if len(q) > 0 {
		path = path + "?" + strings.Join(q, "&")
	}
	// fmt.Printf("PATH: %s\n", path)
	return (*adpt).Get(path, logger)
}

/**** UTILS ****/

func recordPath(id *string, adpt *adapter.Adapter) string {
	path := "/api/v0/registry/records"
	if (*adpt).SkipGateway() {
		path = "/v0/records"
	}
	if id != nil {
		path = path + "/" + *id
	}
	return path
}

package record

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/maxott/magda-cli/pkg/adapter"
)

/**** LIST ****/

type ListCmd struct {
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

func List(cmd *ListCmd, adpt *adapter.Adapter) (ListResult, error) {
	pyl, err := ListRaw(cmd, adpt)
	if err != nil {
		return ListResult{}, err
	}
	res := ListResult{}
	_ = json.Unmarshal(pyl.AsBytes(), &res)
	return res, nil
}

func ListRaw(cmd *ListCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
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
	return (*adpt).Get(path)
}

/**** CREATE ****/

type CreateCmd struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Aspects   Aspects `json:"aspects"`
	SourceTag string  `json:"sourceTag,omitempty"`
}

type Aspects map[string]map[string]interface{}

func CreateRaw(cmd *CreateCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
	r := *cmd
	if r.Id == "" {
		r.Id = uuid.New().String()
	}

	body, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling body. - %s", err)
	}
	// fmt.Printf("RECORD %s\n", body)

	path := recordPath(nil, adpt)
	return (*adpt).Post(path, bytes.NewReader(body))
}

/**** READ ****/

type ReadCmd struct {
	Id         string
	AddAspects string
	Aspect     string
}

func ReadRaw(cmd *ReadCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
	path := recordPath(&cmd.Id, adpt)
	if cmd.AddAspects != "" {
		path = path + "?aspect=" + cmd.AddAspects
	} else if cmd.Aspect != "" {
		path = path + "/aspects/" + cmd.Aspect
	} else {
		// display summary
		path = recordPath(nil, adpt) + "/summary/" + cmd.Id
	}
	return (*adpt).Get(path)
}

/**** UPDATE ****/

type UpdateCmd = CreateCmd

func UpdateRaw(cmd *UpdateCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
	r := *cmd

	path := recordPath(&r.Id, adpt)
	if r.Name == "" {
		// get current 'name' first as it is required
		pld, err := (*adpt).Get(path)
		if err != nil {
			return nil, err
		}
		obj := pld.AsObject()
		if obj == nil {
			return nil, fmt.Errorf("no record body found")
		}
		r.Name = obj["name"].(string)
	}
	body, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling body. - %s", err)
	}
	return (*adpt).Put(path, bytes.NewReader(body))
}

/**** DELETE ****/

type DeleteCmd struct {
	Id         string
	AspectName string
}

func DeleteRaw(cmd *DeleteCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
	path := recordPath(&cmd.Id, adpt)
	if cmd.AspectName != "" {
		path = path + "/aspects/" + cmd.AspectName
	}
	return (*adpt).Delete(path)
}

/**** HISTORY ****/

type HistoryCmd struct {
	Id        string
	EventId   string
	Offset    int
	Limit     int
	PageToken string
}

func HistoryRaw(cmd *HistoryCmd, adpt *adapter.Adapter) (adapter.JsonPayload, error) {
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
	return (*adpt).Get(path)
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

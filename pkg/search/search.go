package search

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

/**** DATASETS ****/

type DatasetsRequest struct {
	Aspects   string
	Query     string
	Offset    int
	Limit     int
	PageToken string
}

type DatasetResult struct {
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

func Dataset(cmd *DatasetRequest, adpt *adapter.Adapter, logger *log.Logger) (DatasetResult, error) {
	pyl, err := DatasetRaw(cmd, adpt, logger)
	if err != nil {
		return DatasetResult{}, err
	}
	res := DatasetResult{}
	_ = json.Unmarshal(pyl.AsBytes(), &res)
	return res, nil
}

func DatasetRaw(cmd *DatasetRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
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

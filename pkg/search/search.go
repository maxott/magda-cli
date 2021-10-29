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
	Query     string
	Offset    int
	Limit     int
}

type DatasetResult struct {
	HitCount       int   `json:"hitCount"`
	DataSets       []struct {
		Title					string	`json:"title"`
		Description				string	`json:"description"`
		Issued					string `json:"issued"`
		Modified				string `json:"modified"`
		Languages				[]string `json:"languages"`
		Publisher				string `json:"publisher"`
		AccrualPeriodicity			string `json:"accrualPeriodicity"`
		AccrualPeriodicityRecurrenceRule	string `json:"accrualPeriodicityRecurrenceRule"`
		Themes					[]string `json:"themes"`
		Keywords				[]string `json:"keywords"`
		ContactPoint				string `json:"contactPoint"`
		LandingPage				string `json:"landingPage"`
		DefaultLicense				string `json:"defaultLicense"`
	} `json:"dataSets"`
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

	if cmd.Query != "" {
		q = append(q, "query="+url.QueryEscape(cmd.Query))
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
	path := "/api/v0/search/datasets"
	if (*adpt).SkipGateway() {
		path = "/v0/datasets"
	}
	if id != nil {
		path = path + "/" + *id
	}
	return path
}

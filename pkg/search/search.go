package search

import (
	"encoding/json"
	//"fmt"
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/maxott/magda-cli/pkg/adapter"
	log "go.uber.org/zap"
)

/**** DATASETS ****/

type DatasetRequest struct {
	Query     string
	Offset    int
	Limit     int
	Publisher string
}

type DatasetResult struct {
	HitCount int `json:"hitCount"`
	DataSets []struct {
		Title                            string   `json:"title"`
		Description                      string   `json:"description"`
		Issued                           string   `json:"issued"`
		Modified                         string   `json:"modified"`
		Languages                        []string `json:"languages"`
		Publisher                        string   `json:"publisher"`
		AccrualPeriodicity               string   `json:"accrualPeriodicity"`
		AccrualPeriodicityRecurrenceRule string   `json:"accrualPeriodicityRecurrenceRule"`
		Themes                           []string `json:"themes"`
		Keywords                         []string `json:"keywords"`
		ContactPoint                     string   `json:"contactPoint"`
		LandingPage                      string   `json:"landingPage"`
		DefaultLicense                   string   `json:"defaultLicense"`
	} `json:"dataSets"`
}

func Dataset(ctxt context.Context, cmd *DatasetRequest, adpt *adapter.Adapter, logger *log.Logger) (DatasetResult, error) {
	pyl, err := DatasetRaw(ctxt, cmd, adpt, logger)
	if err != nil {
		return DatasetResult{}, err
	}
	res := DatasetResult{}
	_ = json.Unmarshal(pyl.AsBytes(), &res)
	return res, nil
}

func DatasetRaw(ctxt context.Context, cmd *DatasetRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := searchPath(nil, adpt)

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
	if cmd.Publisher != "" {
		q = append(q, "publisher="+url.QueryEscape(cmd.Publisher))
	}

	if len(q) > 0 {
		path = path + "?" + strings.Join(q, "&")
	}

	return (*adpt).Get(ctxt, path, logger)
}

/**** UTILS ****/

func searchPath(id *string, adpt *adapter.Adapter) string {
	path := "/api/v0/search/datasets"
	if id != nil {
		path = path + "/" + *id
	}
	return path
}

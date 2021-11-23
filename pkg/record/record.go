package record

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	AndQuery  []QueryTerm
	OrQuery   []QueryTerm
	Offset    int
	Limit     int
	PageToken string
}

type QueryTerm struct {
	Path  string      `json:"path"`
	Op    QueryOp     `json:"op"`
	Value interface{} `json:"value"`
	urlQ  string      // used by cli as it is already 'assembled'
}

func NewQueryTermS(urlQ string) QueryTerm {
	return QueryTerm{urlQ: urlQ}
}

// Create a QueryTerm from a map with keys "path", "op", and "value".
func NewQueryTermI(r map[string]interface{}, logger *func(msg string, fields ...log.Field)) *QueryTerm {
	var p string
	var o QueryOp
	var v interface{}
	var ok bool
	if p, ok = r["path"].(string); !ok {
		if logger != nil {
			(*logger)("missing field", log.String("field", "path"), log.Reflect("in", r))
		}
		return nil
	}
	var os string
	if os, ok = r["op"].(string); !ok {
		if logger != nil {
			(*logger)("missing field", log.String("field", "op"), log.Reflect("in", r))
		}
		return nil
	}
	if o = toQueryOp(os); o == UnknownOp {
		if logger != nil {
			(*logger)("unknown operator", log.String("op", "os"), log.Reflect("in", r))
		}
		return nil
	}
	if v, ok = r["value"]; !ok {
		if logger != nil {
			(*logger)("missing field", log.String("field", "value"), log.Reflect("in", r))
		}
		return nil
	}
	return &QueryTerm{p, o, v, ""}
}

type QueryOp string

const (
	Equal        QueryOp = "=" //  equal
	NotEqual     QueryOp = "!" // not equal
	MatchPattern QueryOp = "?" // matches a pattern, case insensitive. Use Postgresql ILIKE operator.
	// e.g. :?%rating% will match the field contains keyword rating
	// e.g. :?rating% will match the field starts with keyword rating
	NotMatchPattern  QueryOp = "!?" // does not match a pattern, case insensitive. Use Postgresql NOT ILIKE operator
	MatchRegExp      QueryOp = "~"  // matches POSIX regular expression, case insensitive. Use Postgresql ~* operator
	NotMatchRegExp   QueryOp = "!~" // does not match POSIX regular expression, case insensitive. Use Postgresql !~* operator
	GreaterThan      QueryOp = ">"  // greater than
	GreaterEqualThan QueryOp = ">=" // greater than or equal to
	LessThan         QueryOp = "<"  // less than
	LessEqualThen    QueryOp = "<=" // less than or equal
	UnknownOp        QueryOp = "??" // should only be used when a string can't be converted into a legal operator
)

func toQueryOp(s string) QueryOp {
	switch s {
	case "=":
		return Equal
	case "!":
		return NotEqual
	case "?":
		return MatchPattern
	case "!?":
		return NotMatchPattern
	case "~":
		return MatchRegExp
	case "!~":
		return NotMatchRegExp
	case ">":
		return GreaterThan
	case ">=":
		return GreaterEqualThan
	case "<":
		return LessThan
	case "<=":
		return LessEqualThen
	default:
		return UnknownOp
	}

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

func List(ctxt context.Context, cmd *ListRequest, adpt *adapter.Adapter, logger *log.Logger) (ListResult, error) {
	pyl, err := ListRaw(ctxt, cmd, adpt, logger)
	if err != nil {
		return ListResult{}, err
	}
	res := ListResult{}
	_ = json.Unmarshal(pyl.AsBytes(), &res)
	return res, nil
}

func ListRaw(ctxt context.Context, cmd *ListRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := recordPath(nil, adpt)

	pa := []string{}
	if cmd.Aspects != "" {
		pa = append(pa, "aspect="+url.QueryEscape(cmd.Aspects))
	}
	if len(cmd.AndQuery) > 0 {
		for _, q := range cmd.AndQuery {
			pa = append(pa, "aspectQuery="+q.asUrlQuery())
		}
	}
	if len(cmd.OrQuery) > 0 {
		for _, q := range cmd.OrQuery {
			pa = append(pa, "aspectOrQuery="+q.asUrlQuery())
		}
	}
	if cmd.PageToken != "" {
		pa = append(pa, "pageToken="+url.QueryEscape(cmd.PageToken))
	}
	if cmd.Offset >= 0 {
		pa = append(pa, "start="+url.QueryEscape(strconv.Itoa(cmd.Offset)))
	}
	if cmd.Limit >= 0 {
		pa = append(pa, "limit="+url.QueryEscape(strconv.Itoa(cmd.Limit)))
	}
	if len(pa) > 0 {
		path = path + "?" + strings.Join(pa, "&")
	}
	//fmt.Printf("PATH: %s\n", path)
	return (*adpt).Get(ctxt, path, logger)
}

func (t *QueryTerm) asUrlQuery() string {
	if t.urlQ != "" {
		return t.urlQ
	}
	v := fmt.Sprint(t.Value)
	v = strings.ReplaceAll(v, ":", "%3A") // ':' is the separation character, so it needs to be escaped
	op := t.Op
	if op == "=" {
		op = ""
	}
	q := fmt.Sprintf("%s:%s%s", t.Path, op, v)
	return url.QueryEscape(q)
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

func CreateRaw(ctxt context.Context, cmd *CreateRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
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
	return (*adpt).Post(ctxt, path, bytes.NewReader(body), logger)
}

/**** READ ****/

type ReadRequest struct {
	Id         string
	AddAspects string
	Aspect     string
}

func ReadRaw(ctxt context.Context, cmd *ReadRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := recordPath(&cmd.Id, adpt)
	if cmd.AddAspects != "" {
		path = path + "?aspect=" + cmd.AddAspects
	} else if cmd.Aspect != "" {
		path = path + "/aspects/" + cmd.Aspect
	} else {
		// display summary
		path = recordPath(nil, adpt) + "/summary/" + cmd.Id
	}
	return (*adpt).Get(ctxt, path, logger)
}

/**** UPDATE ****/

type UpdateRequest = CreateRequest

func UpdateRaw(ctxt context.Context, cmd *UpdateRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	r := *cmd

	path := recordPath(&r.Id, adpt)
	if r.Name == "" {
		// get current 'name' first as it is required
		pld, err := (*adpt).Get(ctxt, path, logger)
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
	return (*adpt).Put(ctxt, path, bytes.NewReader(body), logger)
}

/**** DELETE ****/

type DeleteRequest struct {
	Id         string
	AspectName string
}

func DeleteRaw(ctxt context.Context, cmd *DeleteRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := recordPath(&cmd.Id, adpt)
	if cmd.AspectName != "" {
		path = path + "/aspects/" + cmd.AspectName
	}
	return (*adpt).Delete(ctxt, path, logger)
}

/**** HISTORY ****/

type HistoryRequest struct {
	Id        string
	EventId   string
	Offset    int
	Limit     int
	PageToken string
}

func HistoryRaw(ctxt context.Context, cmd *HistoryRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
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
	return (*adpt).Get(ctxt, path, logger)
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

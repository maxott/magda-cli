// Program to create, update & delete aspect schemas in Magda
package adapter

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Adapter interface {
	Get(path string) (JsonPayload, error)
	Post(path string, body io.Reader) (JsonPayload, error)
	Put(path string, body io.Reader) (JsonPayload, error)
	Patch(path string, body io.Reader) (JsonPayload, error)
	Delete(path string) (JsonPayload, error)

	SkipGateway() bool // experimental!
}

type ConnectionCtxt struct {
	Host        string
	TenantID    string
	AuthID      string
	AuthKey     string
	UseTLS      bool
	SkipGateway bool
}

// type ReplyHandlerF func(JsonPayload) error

type JsonPayload interface {
	IsObject() bool
	AsObject() map[string]interface{}
	AsArray() []interface{}
	AsBytes() []byte
}

func LoadJsonFromFile(fileName string) (JsonPayload, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	var f interface{}
	err = json.Unmarshal(data, &f)
	if err != nil {
		return nil, err
	}
	m := f.(map[string]interface{})
	return JsonObjPayload{m, data}, nil
}

func ReplyPrinter(pld JsonPayload, err error) error {
	if (err != nil) {
		return err
	}

	var b []byte
	var err2 error
	if pld.IsObject() {
		b, err2 = json.MarshalIndent(pld.AsObject(), "", "  ")
	}
	if !pld.IsObject() {
		b, err2 = json.MarshalIndent(pld.AsArray(), "", "  ")
	}
	if err2 != nil {
		return err2
	} else {
		fmt.Printf("%s\n", b)
		return nil
	}
}

type JsonObjPayload struct {
	payload map[string]interface{}
	bytes []byte
}

func (p JsonObjPayload) IsObject() bool                   { return true }
func (p JsonObjPayload) AsObject() map[string]interface{} { return p.payload }
func (p JsonObjPayload) AsArray() []interface{}           { return []interface{}{p.payload} }
func (p JsonObjPayload) AsBytes() []byte                  { return p.bytes }

type JsonArrPayload struct {
	payload []interface{}
	bytes []byte
}

func (JsonArrPayload) IsObject() bool                     { return false }
func (p JsonArrPayload) AsObject() map[string]interface{} { return nil }
func (p JsonArrPayload) AsArray() []interface{}           { return p.payload }
func (p JsonArrPayload) AsBytes() []byte                  { return p.bytes }

func RestAdapter(ctxt ConnectionCtxt) Adapter {
	return restAdapter{ctxt}
}

type restAdapter struct {
	ctxt ConnectionCtxt
}

func (a restAdapter) Get(path string) (JsonPayload, error) {
	return connect("GET", path, nil, &a.ctxt)
}

func (a restAdapter) Post(path string, body io.Reader) (JsonPayload, error) {
	return connect("POST", path, body, &a.ctxt)
}

func (a restAdapter) Put(path string, body io.Reader) (JsonPayload, error) {
	return connect("PUT", path, body, &a.ctxt)
}

func (a restAdapter) Patch(path string, body io.Reader) (JsonPayload, error) {
	return connect("PATCH", path, body, &a.ctxt)
}

func (a restAdapter) Delete(path string) (JsonPayload, error) {
	return connect("DELETE", path, nil, &a.ctxt)
}

func (a restAdapter) SkipGateway() bool {
	return a.ctxt.SkipGateway
}

func connect(
	method string,
	path string,
	body io.Reader,
	ctxt *ConnectionCtxt,
) (JsonPayload, error) {
	if ctxt.Host == "" {
		log.Fatal("required flag --host not provided, try --help")
	}
	protocol := "http://"
	if ctxt.UseTLS {
		protocol = "https://"
	}
	path = protocol + ctxt.Host + path

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	if ctxt.TenantID != "" {
		req.Header.Set("X-Magda-Tenant-Id", ctxt.TenantID)
	}
	if ctxt.AuthID != "" {
		req.Header.Set("X-Magda-API-Key-Id", ctxt.AuthID)
	}
	if ctxt.AuthKey != "" {
		req.Header.Set("X-Magda-API-Key", ctxt.AuthKey)
	}
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	if resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "Error: %v\n", resp.Status)
		if len(respBody) > 0 {
			fmt.Fprintf(os.Stderr, "%s\n", respBody)
		}
		os.Exit(1)
	}

	return toJSON(respBody)
}

func toJSON(body []byte) (JsonPayload, error) {
	var f interface{}
	err := json.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}

	switch m := f.(type) {
	case []interface{}:
		return JsonArrPayload{m, body}, nil
	case map[string]interface{}:
		return JsonObjPayload{m, body}, nil
	default:
		return nil, fmt.Errorf("unknown json type in body")
	}
}
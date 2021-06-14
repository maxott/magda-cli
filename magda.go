// Program to create, update & delete aspect schemas in Magda
package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	// "io/reader"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

var (
	app = kingpin.New("magda-cli", "Managing records & schemas in Magda.")

	host        = app.Flag("host", "DNS name/IP of Magda host [MAGDA_HOST]").Short('H').Envar("MAGDA_HOST").String()
	tenantID    = app.Flag("tenantID", "Tenant ID [TENANT_ID]").Envar("TENANT_ID").String()
	authID      = app.Flag("authID", "Authorization Key ID [AUTH_ID]").Envar("AUTH_ID").String()
	authKey     = app.Flag("authKey", "Authorization Key [AUTH_KEY]").Envar("AUTH_KEY").String()
	useTLS      = app.Flag("useTLS", "Use https").Default("false").Bool()
	skipGateway = app.Flag("skipGateway", "Skip gateway server and call registry server directly [SKIP_GATEWAY]]").
			Default("false").Envar("SKIP_GATEWAY").Bool()
)

type Command struct {
	id string
}

type ReplyHandlerF func(JsonObjPayload, JsonArrPayload) error

type JsonObjPayload map[string]interface{}
type JsonArrPayload []interface{}

func App() *kingpin.Application {
	return app
}

func Get(path string) error {
	return GetP(path, ReplyPrinter)
}

func GetP(path string, handler ReplyHandlerF) error {
	return connect("GET", path, nil, handler)
}

func Post(path string, body io.Reader) error {
	return PostP( path, body, ReplyPrinter)
}

func PostP(path string, body io.Reader, handler ReplyHandlerF) error {
	return connect("POST", path, body, handler)
}

func Put(path string, body io.Reader) error {
	return PutP( path, body, ReplyPrinter)
}

func PutP(path string, body io.Reader, handler ReplyHandlerF) error {
	return connect("PUT", path, body, handler)
}


func Delete(path string) error {
	return DeleteP( path, ReplyPrinter)
}

func DeleteP(path string, handler ReplyHandlerF) error {
	return connect("DELETE", path, nil, handler)
}

func connect(method string, path string, body io.Reader, replyHandler ReplyHandlerF) error {
	if *host == "" {
		log.Fatal("required flag --host not provided, try --help")
	}
	protocol := "http://"
	if *useTLS {
		protocol = "https://"
	}
	path = protocol + *host + path

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	if tenantID != nil {
		req.Header.Set("X-Magda-Tenant-Id", *tenantID)
	}
	if authID != nil {
		req.Header.Set("X-Magda-API-Key-Id", *authID)
	}
	if authKey != nil {
		req.Header.Set("X-Magda-API-Key", *authKey)
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

	//b, err := ppBody(respBody)
	err = handleReply(respBody, replyHandler)
	if err != nil {
		log.Fatal("Error parsing body. ", err)
	}
	//fmt.Printf("%s\n", b)
	return nil
}

func handleReply(body []byte, printer ReplyHandlerF) error {
	var f interface{}
	err := json.Unmarshal(body, &f)
	if err != nil {
		return err
	}

	switch f.(type) {
	case []interface{}:
		m := f.([]interface{})
		return printer(nil, m)
	case map[string]interface{}:
		m := f.(map[string]interface{})
		return printer(m, nil)
	default:
		return errors.New("unknown json type in body")
	}
}

// func ppBody(body []byte) ([]byte, error) {
// 	var f interface{}
// 	err := json.Unmarshal(body, &f)
// 	if err != nil {
// 		return nil, err
// 	}

// 	switch f.(type) {
// 	case []interface{}:
// 		m := f.([]interface{})
// 		return json.MarshalIndent(m, "", "  ")
// 	case map[string]interface{}:
// 		m := f.(map[string]interface{})
// 		return json.MarshalIndent(m, "", "  ")
// 	default:
// 		return nil, errors.New("unknown json type in body")
// 	}
// }

func ReplyPrinter(obj JsonObjPayload, arr JsonArrPayload) error {
	var b []byte
	var err error
	if (obj != nil) {
		b, err = json.MarshalIndent(obj, "", "  ")
	}
	if (arr != nil) {
		b, err = json.MarshalIndent(arr, "", "  ")
		// if (err != nil) {
		// 	return err
		// }
		// fmt.Printf("%s\n", b)
	}
	if (err != nil) {
		return err
	}
	fmt.Printf("%s\n", b)
	return nil
}

func LoadJsonFromFile(fileName string) (JsonObjPayload, error) {
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
	return m, nil
}

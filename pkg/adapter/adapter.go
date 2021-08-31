// Program to create, update & delete aspect schemas in Magda
package adapter

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/maxott/magda-cli/pkg/log"
)

type ConnectionCtxt struct {
	Host        string
	TenantID    string
	AuthID      string
	AuthKey     string
	JwtToken    string
	UseTLS      bool
	SkipGateway bool
}

func CreateJwtToken(userID *string, signingSecret *string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": *userID,
		"iat":    time.Now().Unix(),
	})
	return token.SignedString([]byte(*signingSecret))
}

func LoadJsonFromFile(fileName string) (Payload, error) {
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

func RestAdapter(ctxt ConnectionCtxt) Adapter {
	return restAdapter{ctxt}
}

type restAdapter struct {
	ctxt ConnectionCtxt
}

func (a restAdapter) Get(path string, logger log.Logger) (Payload, error) {
	return connect("GET", path, nil, &a.ctxt, logger)
}

func (a restAdapter) Post(path string, body io.Reader, logger log.Logger) (Payload, error) {
	return connect("POST", path, body, &a.ctxt, logger)
}

func (a restAdapter) Put(path string, body io.Reader, logger log.Logger) (Payload, error) {
	return connect("PUT", path, body, &a.ctxt, logger)
}

func (a restAdapter) Patch(path string, body io.Reader, logger log.Logger) (Payload, error) {
	return connect("PATCH", path, body, &a.ctxt, logger)
}

func (a restAdapter) Delete(path string, logger log.Logger) (Payload, error) {
	return connect("DELETE", path, nil, &a.ctxt, logger)
}

func (a restAdapter) SkipGateway() bool {
	return a.ctxt.SkipGateway
}

func connect(
	method string,
	path string,
	body io.Reader,
	ctxt *ConnectionCtxt,
	logger log.Logger,
) (Payload, error) {
	logger = logger.With("method", method).With("path", path)
	if ctxt.Host == "" {
		return nil, logger.Error(nil, "required flag --host not provided, try --help")
	}
	protocol := "http://"
	if ctxt.UseTLS {
		protocol = "https://"
	}
	path = protocol + ctxt.Host + path

	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, logger.Error(err, "Error reading request. ")
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
	if ctxt.JwtToken != "" {
		req.Header.Set("X-Magda-Session", ctxt.JwtToken)
	}

	client := &http.Client{Timeout: time.Second * 10}
	logger.Debug("Calling magda registry")
	resp, err := client.Do(req)
	if err != nil {
		return nil, logger.Error(err, "HTTP request failed.")
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, logger.Error(err, "Accessing respone body failed.")
	}

	if resp.StatusCode >= 300 {
		if len(respBody) > 0 {
			logger = logger.With("body", string(respBody))
		}
		return nil, logger.With("statusCode", resp.StatusCode).Error(nil, "Error response")
	}
	contentType := resp.Header.Get("Content-Type") 
	return ToPayload(respBody, contentType, logger)
}



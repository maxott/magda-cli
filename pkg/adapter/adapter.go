// Program to create, update & delete aspect schemas in Magda
package adapter

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	log "go.uber.org/zap"
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

func RestAdapter(ctxt ConnectionCtxt) Adapter {
	return restAdapter{ctxt}
}

type restAdapter struct {
	ctxt ConnectionCtxt
}

func (a restAdapter) Get(path string, logger *log.Logger) (Payload, error) {
	return connect("GET", path, nil, &a.ctxt, logger)
}

func (a restAdapter) Post(path string, body io.Reader, logger *log.Logger) (Payload, error) {
	return connect("POST", path, body, &a.ctxt, logger)
}

func (a restAdapter) Put(path string, body io.Reader, logger *log.Logger) (Payload, error) {
	return connect("PUT", path, body, &a.ctxt, logger)
}

func (a restAdapter) Patch(path string, body io.Reader, logger *log.Logger) (Payload, error) {
	return connect("PATCH", path, body, &a.ctxt, logger)
}

func (a restAdapter) Delete(path string, logger *log.Logger) (Payload, error) {
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
	logger *log.Logger,
) (Payload, error) {
	logger = logger.With(log.String("method", method),log.String("path", path) )
	if ctxt.Host == "" {
		logger.Error("Missing 'host'")
		return nil, fmt.Errorf("missing magda host")
	}
	protocol := "http://"
	if ctxt.UseTLS {
		protocol = "https://"
	}
	url := protocol + ctxt.Host + path
	logger = logger.With(log.String("url", url))
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logger.Error("Creating http request", log.Error(err))
		return nil, err
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
		logger.Error("HTTP request failed.", log.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Accessing response body failed.", log.Error(err))
		return nil, err
	}

	if resp.StatusCode >= 300 {
		if len(respBody) > 0 {
			logger = logger.With(log.ByteString("body", respBody))
		}
		logger.Error("HTTP response", log.Int("statusCode", resp.StatusCode))
		return nil, fmt.Errorf("Error response")
	}
	contentType := resp.Header.Get("Content-Type") 
	return ToPayload(respBody, contentType, logger)
}



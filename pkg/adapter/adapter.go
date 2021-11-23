// Program to create, update & delete aspect schemas in Magda
package adapter

import (
	"context"
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

func RestAdapter(connCtxt ConnectionCtxt) Adapter {
	return restAdapter{connCtxt}
}

type IAdapterError interface {
	Error() string
	Path() string
}
type AdapterError struct {
	path string
}

func (e *AdapterError) Path() string { return e.path }

func (e *AdapterError) Error() string { return "Generic magda adapter error" }

type MissingHostError struct {
	AdapterError
}

func (e MissingHostError) Error() string { return "Missing host name" }

type ResourceNotFoundError struct {
	AdapterError
}

func (e ResourceNotFoundError) Error() string { return "Resource not found" }

type UnauthorizedError struct {
	AdapterError
}

func (e *UnauthorizedError) Error() string { return "Unauthorized access" }

type MagdaError struct {
	AdapterError
	StatusCode int
	Message    string
}

func (e *MagdaError) Error() string { return e.Message }

type ClientError struct {
	AdapterError
	err error
}

func (e *ClientError) Error() string {
	return fmt.Sprintf("while connecting to Magda registry - %s", e.err.Error())
}

type restAdapter struct {
	ctxt ConnectionCtxt
}

func (a restAdapter) Get(ctxt context.Context, path string, logger *log.Logger) (Payload, error) {
	return connect(ctxt, "GET", path, nil, &a.ctxt, logger)
}

func (a restAdapter) Post(ctxt context.Context, path string, body io.Reader, logger *log.Logger) (Payload, error) {
	return connect(ctxt, "POST", path, body, &a.ctxt, logger)
}

func (a restAdapter) Put(ctxt context.Context, path string, body io.Reader, logger *log.Logger) (Payload, error) {
	return connect(ctxt, "PUT", path, body, &a.ctxt, logger)
}

func (a restAdapter) Patch(ctxt context.Context, path string, body io.Reader, logger *log.Logger) (Payload, error) {
	return connect(ctxt, "PATCH", path, body, &a.ctxt, logger)
}

func (a restAdapter) Delete(ctxt context.Context, path string, logger *log.Logger) (Payload, error) {
	return connect(ctxt, "DELETE", path, nil, &a.ctxt, logger)
}

func (a restAdapter) SkipGateway() bool {
	return a.ctxt.SkipGateway
}

func connect(
	ctxt context.Context,
	method string,
	path string,
	body io.Reader,
	connCtxt *ConnectionCtxt,
	logger *log.Logger,
) (Payload, error) {
	logger = logger.With(log.String("method", method), log.String("path", path))
	if connCtxt.Host == "" {
		logger.Error("Missing 'host'")
		return nil, &MissingHostError{AdapterError{path}}
	}
	protocol := "http://"
	if connCtxt.UseTLS {
		protocol = "https://"
	}
	url := protocol + connCtxt.Host + path
	logger = logger.With(log.String("url", url))
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logger.Error("Creating http request", log.Error(err))
		return nil, &ClientError{AdapterError{path}, err}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	if connCtxt.TenantID != "" {
		req.Header.Set("X-Magda-Tenant-Id", connCtxt.TenantID)
	}
	if connCtxt.AuthID != "" {
		req.Header.Set("X-Magda-API-Key-Id", connCtxt.AuthID)
	}
	if connCtxt.AuthKey != "" {
		req.Header.Set("X-Magda-API-Key", connCtxt.AuthKey)
	}
	if connCtxt.JwtToken != "" {
		req.Header.Set("X-Magda-Session", connCtxt.JwtToken)
	}

	client := &http.Client{Timeout: time.Second * 10}
	logger.Debug("Calling magda registry")
	resp, err := client.Do(req)
	if err != nil {
		logger.Warn("HTTP request failed.", log.Error(err))
		return nil, &ClientError{AdapterError{path}, err}
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Warn("Accessing response body failed.", log.Error(err))
		return nil, &ClientError{AdapterError{path}, err}
	}

	if resp.StatusCode >= 300 {
		if len(respBody) > 0 {
			logger = logger.With(log.ByteString("body", respBody))
		}
		logger.Warn("HTTP response", log.Int("statusCode", resp.StatusCode))
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, &ResourceNotFoundError{AdapterError{path}}
		case http.StatusUnauthorized:
			return nil, &UnauthorizedError{AdapterError{path}}
		default:
			return nil, &MagdaError{
				AdapterError{path},
				resp.StatusCode,
				string(respBody),
			}
		}

		//ResourceNotFoundError
	}
	contentType := resp.Header.Get("Content-Type")
	return ToPayload(respBody, contentType, logger)
}

package record

import (
	"encoding/json"
	_ "fmt"
	_ "regexp"
	"testing"
)

func TestListQuery(t *testing.T) {
	qt := QueryTerm{
		Path:  "foo",
		Op:    Equal,
		Value: "a",
	}
	p := qt.asUrlQuery()
	exp := "foo%3Aa"
	if exp != p {
		t.Errorf("expected '%s', but got '%s'", exp, p)
	}
}

func TestListQuery2(t *testing.T) {
	qt := QueryTerm{
		Path:  "foo",
		Op:    Equal,
		Value: "a:b",
	}
	p := qt.asUrlQuery()
	exp := "foo%3Aa%253Ab"
	if exp != p {
		t.Errorf("expected '%s', but got '%s'", exp, p)
	}
}

func TestListQueryJson(t *testing.T) {
	j := `{
		"path": "asp.foo",
		"op": "=",
		"value": "a:b"
	}`
	qt := QueryTerm{}
	err := json.Unmarshal([]byte(j), &qt)
	if err != nil {
		t.Errorf("while unmarshal query - %v", err)
		return
	}

	p := qt.asUrlQuery()
	exp := "asp.foo%3Aa%253Ab"
	if exp != p {
		t.Errorf("expected '%s', but got '%s'", exp, p)
	}
}

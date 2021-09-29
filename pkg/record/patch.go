package record

import (
	"bytes"
	"encoding/json"

	"github.com/maxott/magda-cli/pkg/adapter"
	log "go.uber.org/zap"
)

/**** PATCH ASPECT ********/

type PatchAspectRequest struct {
	Id     string
	Aspect string
	Patch  []PatchOp
}

type PatchOp interface {
}

func PatchAspectRaw(cmd *PatchAspectRequest, adpt *adapter.Adapter, logger *log.Logger) (adapter.Payload, error) {
	path := recordPath(&cmd.Id, adpt) + "/aspects/" + cmd.Aspect
	body, err := json.MarshalIndent(cmd.Patch, "", "  ")
	if err != nil {
		logger.Error("marshalling body", log.Error(err))
		return nil, err
	}
	return (*adpt).Patch(path, bytes.NewReader(body), logger)
}

type patch struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

// { "op": "add", "path": "/biscuits/1", "value": { "name": "Ginger Nut" } }
func PatchAddOp(path string, value interface{}) PatchOp {
	return patch{
		Op:    "add",
		Path:  path,
		Value: value,
	}
}

// { "op": "remove", "path": "/biscuits" }
func PatchRemoveOp(path string) PatchOp {
	return patch{
		Op:   "remove",
		Path: path,
	}
}

// { "op": "replace", "path": "/biscuits/0/name", "value": "Chocolate Digestive" }
func PatchReplaceOp(path string, value interface{}) PatchOp {
	return patch{
		Op:    "replace",
		Path:  path,
		Value: value,
	}
}

type patch2 struct {
	Op   string `json:"op"`
	From string `json:"from"`
	Path string `json:"path"`
}

// { "op": "copy", "from": "/biscuits/0", "path": "/best_biscuit" }
func PatchCopyOp(from string, path string) PatchOp {
	return patch2{
		Op:   "copy",
		From: from,
		Path: path,
	}
}

// { "op": "move", "from": "/biscuits", "path": "/cookies" }
func PatchMoveOp(from string, path string) PatchOp {
	return patch2{
		Op:   "move",
		From: from,
		Path: path,
	}
}

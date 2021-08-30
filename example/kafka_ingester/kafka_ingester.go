package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	// "reflect"

	"github.com/maxott/magda-cli/pkg/record"
	"github.com/maxott/magda-cli/pkg/adapter"
	
	jsonpatch "github.com/evanphx/json-patch"
	kafka "github.com/segmentio/kafka-go"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("kafka-ingester", "Ingest Kafka messages into Magda.")

	brocker    = app.Flag("k:brocker", "Address of Kafka broker (e.g. localhost:8888) [KAFKA_BROKER]").Short('b').Envar("KAFKA_BROKER").String()
	topic    = app.Flag("k:topic", "Kafka topic to listen to [KAFKA_TOPIC]").Short('t').Envar("KAFKA_TOPIC").String()
	groupID = app.Flag("k:groupID", "Kafka consumer group [KAFKA_GROUP_ID]").Short('g').Envar("KAFKA_GROUP_ID").String()
	offset = app.Flag("k:offset", "Kafka message offset [KAFKA_OFFSET]").Short('o').Envar("KAFKA_OFFSET").Int64()

	patchFile = app.Flag("x:patchFile", "Optional JSON Patch file to transform message [PATCH_FILE]").Short('p').Envar("PATCH_FILE").String()


	recordID        = app.Flag("m:id", "Record ID to append kafka record").Short('i').Required().String()
	aspectName      = app.Flag("m:aspectNmae", "Name of aspect to append kafka record").Short('a').Required().String()

	host        = app.Flag("m:host", "DNS name/IP of Magda host [MAGDA_HOST]").Short('H').Envar("MAGDA_HOST").String()
	tenantID    = app.Flag("m:tenantID", "Tenant ID [MAGDA_TENANT_ID]").Envar("MAGDA_TENANT_ID").String()
	authID      = app.Flag("m:authID", "Authorization Key ID [MAGDA_AUTH_ID]").Envar("MAGDA_AUTH_ID").String()
	authKey     = app.Flag("m:authKey", "Authorization Key [MAGDA_AUTH_KEY]").Envar("MAGDA_AUTH_KEY").String()
	useTLS      = app.Flag("m:useTLS", "Use https").Default("false").Bool()

	skipGateway = app.Flag("skipGateway", "Skip gateway server and call registry server directly [MAGDA_SKIP_GATEWAY]]").
			Default("false").Envar("MAGDA_SKIP_GATEWAY").Bool()

	verbose = app.Flag("verbose", "Add more logging").Short('v').Default("false").Bool()
)

func run() error {
	r := init_kafka()
	patch := init_patch()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		if err = process(m.Value, patch, m.Offset); err != nil {
			return err
		}
	}
}

func init_kafka() *kafka.Reader {
	if *brocker == "" {
		app.Fatalf("Missing --k:brocker")
	}
	if *topic == "" {
		app.Fatalf("Missing --k:topic")
	}

	cfg := kafka.ReaderConfig{
		Brokers:   []string{*brocker},
		Topic:     *topic,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	}
	if (groupID != nil) {
		cfg.GroupID = *groupID
	}
	r := kafka.NewReader(cfg)
	if offset != nil {
		r.SetOffset(*offset)
	}
	return r
}

func init_patch() jsonpatch.Patch {
	var patch jsonpatch.Patch = nil
	var err error
	if *patchFile != "" {
		if patch, err = loadPatch(*patchFile); err != nil {
			app.Fatalf("Error while loading patch file '%s' - %+v", *patchFile, err)
		}
	}
	return patch
}

func loadPatch(fileName string) (jsonpatch.Patch, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return jsonpatch.DecodePatch(data)
}

func magdaAdapter() *adapter.Adapter {
	adapter := adapter.RestAdapter(adapter.ConnectionCtxt{
		Host: *host, TenantID: *tenantID, AuthID: *authID, AuthKey: *authKey, UseTLS: *useTLS, SkipGateway: *skipGateway,
	})
	return &adapter
}

type PatchOp struct {
	Op string `json:"op"`
	Path string `json:"path"`
	Value map[string]interface{} `json:"value,omitempty"`
}

func process(data []byte, patch jsonpatch.Patch, msgID int64) error {
	rec, err := convert(data, patch)
	if err != nil {
		return err
	}

	// { "op": "add", "path": "/biscuits/1", "value": { "name": "Ginger Nut" } }

	cmd := record.PatchAspectCmd{
		Id: *recordID,
		Aspect: *aspectName,
		Patch: []interface{}{	
			PatchOp{
				Op: "add",
				Path: "/requests/-",
				Value: rec,
			},
		},
	}
	_, err = record.PatchAspectRaw(&cmd, magdaAdapter())
	if err == nil && *verbose {
		log.Printf("Successfully added message %d", msgID)
	}
	return err
}

func convert(data []byte, patch jsonpatch.Patch) (map[string]interface{}, error) {
	d := data
	if patch != nil {
		var err error
		d, err = patch.ApplyIndent(d, "  ")
		if err != nil {
			return nil, fmt.Errorf("while patching json - %v", err)
		}
	} 
	if *verbose {
		log.Printf("%s", d)
	}

	var m map[string]interface{}
	if err := json.Unmarshal(d, &m); err != nil {
		return nil, fmt.Errorf("while decoding kafka value")
	}
	return m , nil
}

func main() {
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	//patch()
	if err := run(); err != nil {
		app.Fatalf("ERROR %v", err)
	}
}

package main

import (
	"fmt"
	"errors"
	"log"
	"encoding/json"
	"bytes"

	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	cmd := App().Command("schema", "Managing aspect schemas")
	cliSchemaList(cmd)
	cliSchemaCreate(cmd)
	cliSchemaRead(cmd)
	cliSchemaUpdate(cmd)
}

/**** LIST ****/

func cliSchemaList(topCmd *kingpin.CmdClause) {
	topCmd.Command("list", "List all aspect schemas").Action(func(_ *kingpin.ParseContext) error {
		return schemaList()
	})
}

func schemaList() error {
	path := aspectPath(nil)
	fmt.Printf("LIST PATH %s\n", path)
	return Get(path)
}

/**** CREATE ****/

type SchemaCreate struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	SchemaFile string  `json:"-"`
	Schema     JsonObjPayload `json:"jsonSchema"`
}

func cliSchemaCreate(topCmd *kingpin.CmdClause) {
	r := &SchemaCreate{}
	c := topCmd.Command("create", "Creates a new schema").Action(func(_ *kingpin.ParseContext) error {
		return schemaCreate(r)
	})
	c.Flag("name", "Descriptive name").
		Short('n').
		Required().
		StringVar(&r.Name)
		cliAddSchemaCUFlags(r, c)
}

func cliAddSchemaCUFlags(r *SchemaCreate, c *kingpin.CmdClause) {
	c.Flag("id", "Schema ID").
		Short('i').
		Required().
		StringVar(&r.Id)
	c.Flag("schemaFile", "File containing schema/aspect decalration").
		Short('f').
		Required().
		ExistingFileVar(&r.SchemaFile)
}

func schemaCreate(args *SchemaCreate) error {
	r := *args

	body := createCUBody(&r)

	path := aspectPath(nil)
	// return Post(path, bytes.NewReader(body))
	return PostP(path, bytes.NewReader(body), func(obj JsonObjPayload, arr JsonArrPayload) error {
		fmt.Printf("Successfully create schema '%s'\n", r.Id)
		return nil
	})
}

func createCUBody(r *SchemaCreate) []byte {
	js, err := LoadJsonFromFile(r.SchemaFile)
	r.Schema = js
	if err != nil {
		App().Fatalf("failed to load & verify '%s' - %s", r.SchemaFile, err)
	}
	body, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling body. ", err)
	}
	return body
}

/**** READ ****/

type SchemaRead struct {
	Id string
}

func cliSchemaRead(topCmd *kingpin.CmdClause) {
	r := &SchemaRead{}
	c := topCmd.Command("read", "Read the content of a record").Action(func(_ *kingpin.ParseContext) error {
		return schemaRead(r)
	})
	c.Flag("id", "Record ID").
		Short('i').
		Required().
		StringVar(&r.Id)
}

func schemaRead(args *SchemaRead) error {
	path := aspectPath(&args.Id)
	return Get(path)
}

/**** UPDATE ****/

func cliSchemaUpdate(topCmd *kingpin.CmdClause) {
	r := &SchemaCreate{}
	c := topCmd.Command("update", "Update existing schema").Action(func(_ *kingpin.ParseContext) error {
		return schemaUpdate(r)
	})
	c.Flag("name", "Descriptive name").
		Short('n').
		StringVar(&r.Name)
		cliAddSchemaCUFlags(r, c)
}

func schemaUpdate(args *SchemaCreate) error {
	r := *args
	
	path := aspectPath(&r.Id)
	
	f := func() error {
		body := createCUBody(&r)
		// return Put(path, bytes.NewReader(body))
		return PutP(path, bytes.NewReader(body), func(obj JsonObjPayload, arr JsonArrPayload) error {
			fmt.Printf("Successfully updated schema '%s'\n", r.Id)
			return nil
		})
	}

	if (r.Name == "") {
		// get current 'name' first as it is required
		return GetP(path, func(obj JsonObjPayload, arr JsonArrPayload) error {
			if (obj == nil) {
				return errors.New("No schema body found")
			}
			r.Name = obj["name"].(string)
			return f()
		})
	} else {
		return f()
	}
}

/**** DELETE ****/

// Not supported 

/**** Utils ****/
func aspectPath(id *string) string {
	path := "/api/v0/registry/aspects"
	if *skipGateway {
		path = "/v0/aspects"
	}
	if id != nil {
		path = path + "/" + *id
	}
	return path
}

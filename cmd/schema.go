package cmd

import (
	"fmt"

	"github.com/maxott/magda-cli/pkg/adapter"
	"github.com/maxott/magda-cli/pkg/schema"
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
	cmd := &schema.ListCmd{}
	topCmd.Command("list", "List all aspect schemas").Action(func(_ *kingpin.ParseContext) error {
		return adapter.ReplyPrinter(schema.ListRaw(cmd, Adapter()))
	})
}

/**** CREATE ****/

type SchemaCreate struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	SchemaFile string `json:"-"`
}

func cliSchemaCreate(topCmd *kingpin.CmdClause) {
	r := &SchemaCreate{}
	c := topCmd.Command("create", "Creates a new schema").Action(func(_ *kingpin.ParseContext) error {
		if js, err := adapter.LoadJsonFromFile(r.SchemaFile); err != nil {
			return fmt.Errorf("failed to load & verify '%s' - %s", r.SchemaFile, err)
		} else {
			cmd := schema.CreateCmd{
				Id: r.Id, Name: r.Name, Schema: js.AsObject(),
			}
			if _, err := schema.CreateRaw(&cmd, Adapter()); err == nil {
				fmt.Printf("Successfully create schema '%s'\n", r.Id)
				return nil
			} else {
				return err
			}
		}
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
	c.Flag("schema-file", "File containing schema/aspect decalration").
		Short('f').
		Required().
		ExistingFileVar(&r.SchemaFile)
}

/**** READ ****/

func cliSchemaRead(topCmd *kingpin.CmdClause) {
	r := &schema.ReadCmd{}
	c := topCmd.Command("read", "Read the content of a schema").Action(func(_ *kingpin.ParseContext) error {
		return adapter.ReplyPrinter(schema.ReadRaw(r, Adapter()))
	})
	c.Flag("id", "Record ID").
		Short('i').
		Required().
		StringVar(&r.Id)
}

/**** UPDATE ****/

func cliSchemaUpdate(topCmd *kingpin.CmdClause) {
	r := &SchemaCreate{}
	c := topCmd.Command("update", "Update existing schema").Action(func(_ *kingpin.ParseContext) error {
		if js, err := adapter.LoadJsonFromFile(r.SchemaFile); err != nil {
			return fmt.Errorf("failed to load & verify '%s' - %s", r.SchemaFile, err)
		} else {
			cmd := schema.UpdateCmd{
				Id: r.Id, Name: r.Name, Schema: js.AsObject(),
			}
			if _, err := schema.UpdateRaw(&cmd, Adapter()); err == nil {
				fmt.Printf("Successfully updated schema '%s'\n", r.Id)
				return nil
			} else {
				return err
			}
		}
	})
	c.Flag("name", "Descriptive name").
		Short('n').
		StringVar(&r.Name)
	cliAddSchemaCUFlags(r, c)
}

/**** DELETE ****/

// Not supported

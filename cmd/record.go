package cmd

import (
	"fmt"

	"github.com/maxott/magda-cli/pkg/adapter"
	"github.com/maxott/magda-cli/pkg/record"
	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	cmd := App().Command("record", "Managing magda records")
	cliRecordList(cmd)
	cliRecordRead(cmd)
	cliRecordCreate(cmd)
	cliRecordUpdate(cmd)
	cliRecordDelete(cmd)
	cliRecordHistory(cmd)
}

/**** LIST ****/

func cliRecordList(topCmd *kingpin.CmdClause) {
	r := &record.ListCmd{Offset: -1, Limit: -1}
	c := topCmd.Command("list", "List some records").Action(func(_ *kingpin.ParseContext) error {
		// if rec, err := record.List(r, Adapter()); err != nil {
		// 	return err;
		// } else {
		// 	fmt.Printf("List: %+v\n", rec)
		// 	return nil
		// }

		return adapter.ReplyPrinter(record.ListRaw(r, Adapter()))
	})
	c.Flag("aspects", "The aspects for which to retrieve data").
		Short('a').
		StringVar(&r.Aspects)
	c.Flag("query", "Record Name").
		Short('q').
		StringVar(&r.Query)
	c.Flag("offset", "Index of first record retrieved").
		Short('o').
		IntVar(&r.Offset)
	c.Flag("limit", "The maximumm number of records to retrieve").
		Short('l').
		IntVar(&r.Limit)
	c.Flag("page-token", "Token that identifies the start of a page of results").
		Short('t').
		StringVar(&r.PageToken)
}

/**** CREATE ****/

type CreateCmd struct {
	Id         string         `json:"id"`
	Name       string         `json:"name"`
	Aspects    record.Aspects `json:"aspects"`
	SourceTag  string         `json:"sourceTag,omitempty"`
	AspectName string         `json:"-"`
	AspectFile string         `json:"-"`
}

func cliRecordCreate(topCmd *kingpin.CmdClause) {
	r := &CreateCmd{}
	c := topCmd.Command("create", "Creates a new record").Action(func(_ *kingpin.ParseContext) error {
		addAspects(r)
		cmd := record.CreateCmd{
			Id: r.Id, Name: r.Name, Aspects: r.Aspects, SourceTag: r.SourceTag,
		}
		if _, err := record.CreateRaw(&cmd, Adapter()); err == nil {
			fmt.Printf("Successfully create record '%s'\n", r.Id)
			return nil
		} else {
			return err
		}
	})
	c.Flag("id", "Record ID (defaults to UUID)").
		Short('i').
		StringVar(&r.Id)
	c.Flag("name", "Record Name").
		Short('n').
		Required().
		StringVar(&r.Name)
	cliAddAspectFlags(r, c)
}

func cliAddAspectFlags(r *CreateCmd, c *kingpin.CmdClause) {
	c.Flag("aspect-name", "Name of aspect to add (requires --aspectFile)").
		Short('a').
		StringVar(&r.AspectName)
	c.Flag("aspect-file", "File containing aspect data").
		Short('f').
		ExistingFileVar(&r.AspectFile)
}

func addAspects(r *CreateCmd) {
	r.Aspects = record.Aspects{}
	if r.AspectName != "" || r.AspectFile != "" {
		if r.AspectName == "" {
			App().Fatalf("required flag --aspect-name not provided, try --help")
		}
		if r.AspectFile == "" {
			App().Fatalf("required flag --aspect-file not provided, try --help")
		}
		adata, err := adapter.LoadJsonFromFile(r.AspectFile)
		if err != nil {
			App().Fatalf("failed to load & verify '%s' - %s", r.AspectFile, err)
		}
		r.Aspects[r.AspectName] = adata.AsObject()
	}
}

/**** READ ****/

func cliRecordRead(topCmd *kingpin.CmdClause) {
	r := &record.ReadCmd{}
	c := topCmd.Command("read", "Read the content of a record").Action(func(_ *kingpin.ParseContext) error {
		return adapter.ReplyPrinter(record.ReadRaw(r, Adapter()))
	})
	c.Flag("id", "Record ID").
		Short('i').
		Required().
		StringVar(&r.Id)
	c.Flag("add-aspects", "Add aspects to record listing (comma separated)").
		StringVar(&r.AddAspects)
	c.Flag("aspect", "Show only this aspect of the record as result").
		Short('a').
		StringVar(&r.Aspect)
}

/**** UPDATE ****/

func cliRecordUpdate(topCmd *kingpin.CmdClause) {
	r := &CreateCmd{}
	c := topCmd.Command("update", "Update an existing record").Action(func(_ *kingpin.ParseContext) error {
		addAspects(r)
		cmd := record.UpdateCmd{
			Id: r.Id, Name: r.Name, Aspects: r.Aspects, SourceTag: r.SourceTag,
		}
		if _, err := record.UpdateRaw(&cmd, Adapter()); err == nil {
			fmt.Printf("Successfully updated record '%s'\n", r.Id)
			return nil
		} else {
			return err
		}
	})
	c.Flag("id", "Record ID (defaults to UUID)").
		Short('i').
		Required().
		StringVar(&r.Id)
	c.Flag("name", "Record Name").
		Short('n').
		StringVar(&r.Name)
	cliAddAspectFlags(r, c)
}

/**** DELETE ****/

func cliRecordDelete(topCmd *kingpin.CmdClause) {
	r := &record.DeleteCmd{}
	c := topCmd.Command("delete", "Delete a record or one of it's aspects").Action(func(_ *kingpin.ParseContext) error {
		if _, err := record.DeleteRaw(r, Adapter()); err == nil {
			fmt.Printf("Successfully deleted record '%s'\n", r.Id)
			return nil
		} else {
			return err
		}
	})
	c.Flag("id", "Record ID").
		Short('i').
		Required().
		StringVar(&r.Id)
	c.Flag("aspect", "Only delete this aspect").
		Short('a').
		StringVar(&r.AspectName)
}

/**** HISTORY ****/

func cliRecordHistory(topCmd *kingpin.CmdClause) {
	r := &record.HistoryCmd{Offset: -1, Limit: -1}
	c := topCmd.Command("history", "Get a list of all events for a record").Action(func(_ *kingpin.ParseContext) error {
		return adapter.ReplyPrinter(record.HistoryRaw(r, Adapter()))
	})
	c.Flag("id", "Record ID").
		Short('i').
		Required().
		StringVar(&r.Id)
	c.Flag("event-id", "Only show event with event-id").
		Short('e').
		StringVar(&r.EventId)
	c.Flag("offset", "Index of first record retrieved").
		Short('o').
		IntVar(&r.Offset)
	c.Flag("limit", "The maximumm number of records to retrieve").
		Short('l').
		IntVar(&r.Limit)
	c.Flag("page-token", "Token that identifies the start of a page of results").
		Short('t').
		StringVar(&r.PageToken)
}
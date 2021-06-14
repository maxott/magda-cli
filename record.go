package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"errors"

	"github.com/google/uuid"
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

type RecordList struct {
	Aspects string
	Query string
	Offset int
	Limit int
	PageToken string
}

func cliRecordList(topCmd *kingpin.CmdClause) {
	r := &RecordList{Offset: -1, Limit: -1}
	c := topCmd.Command("list", "List some records").Action(func(_ *kingpin.ParseContext) error {
		return recordList(r)
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
	c.Flag("pageToken", "Token that identifies the start of a page of results").
		Short('t').
		StringVar(&r.PageToken)
}

func recordList(args *RecordList) error {
	path := recordPath(nil)

	q := []string{}
	if args.Aspects != "" {
		q = append(q, "aspect=" + url.QueryEscape(args.Aspects))
	}
	if args.Query != "" {
		q = append(q, "aspectQuery=" + url.QueryEscape(args.Query))
	}
	if args.PageToken != "" {
		q = append(q, "pageToken=" + url.QueryEscape(args.PageToken))
	}
	if args.Offset >= 0 {
		q = append(q, "start=" + url.QueryEscape(strconv.Itoa(args.Offset)))
	}
	if args.Limit >= 0 {
		q = append(q, "limit=" + url.QueryEscape(strconv.Itoa(args.Limit)))
	}
	if len(q) > 0 {
		path = path + "?" + strings.Join(q,"&")
	}
	//fmt.Printf("PATH: %s\n", path)
	return Get(path)
}

/**** CREATE ****/

type RecordCreate struct {
	Id         string  `json:"id"`
	Name       string  `json:"name"`
	Aspects    Aspects `json:"aspects"`
	SourceTag  string  `json:"sourceTag,omitempty"`
	AspectName string  `json:"-"`
	AspectFile string  `json:"-"`
}

type Aspects map[string]JsonObjPayload

func cliRecordCreate(topCmd *kingpin.CmdClause) {
	r := &RecordCreate{}
	c := topCmd.Command("create", "Creates a new record").Action(func(_ *kingpin.ParseContext) error {
		return recordCreate(r)
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

func cliAddAspectFlags(r *RecordCreate, c *kingpin.CmdClause) {
	c.Flag("aspectName", "Name of aspect to add (requires --aspectFile)").
		Short('a').
		StringVar(&r.AspectName)
	c.Flag("aspectFile", "File containing aspect data").
		Short('f').
		ExistingFileVar(&r.AspectFile)
}

func recordCreate(args *RecordCreate) error {
	r := *args
	if r.Id == "" {
		r.Id = uuid.New().String()
	}

	addAspects(&r)
	body, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling body. ", err)
	}
	// fmt.Printf("RECORD %s\n", body)

	path := recordPath(nil)
	return PostP(path, bytes.NewReader(body), func(obj JsonObjPayload, arr JsonArrPayload) error {
		fmt.Printf("Successfully create record '%s'\n", r.Id)
		return nil
	})
}

func addAspects(r *RecordCreate) {
	r.Aspects = Aspects{}
	if r.AspectName != "" || r.AspectFile != "" {
		if r.AspectName == "" {
			App().Fatalf("required flag --aspectNmae not provided, try --help")
		}
		if r.AspectFile == "" {
			App().Fatalf("required flag --aspectFile not provided, try --help")
		}
		adata, err := LoadJsonFromFile(r.AspectFile)
		if err != nil {
			App().Fatalf("failed to load & verify '%s' - %s", r.AspectFile, err)
		}
		r.Aspects[r.AspectName] = adata
	}
}

/**** READ ****/

type RecordRead struct {
	Id string
	AddAspects string
	Aspect string
}

func cliRecordRead(topCmd *kingpin.CmdClause) {
	r := &RecordRead{}
	c := topCmd.Command("read", "Read the content of a record").Action(func(_ *kingpin.ParseContext) error {
		return recordRead(r)
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

func recordRead(args *RecordRead) error {
	path := recordPath(&args.Id)
	if args.AddAspects != "" {
		path = path + "?aspect=" + args.AddAspects
	} else if args.Aspect != "" {
		path = path + "/aspects/" + args.Aspect
	} else {
		// display summary
		path = recordPath(nil) + "/summary/" + args.Id
	}
	return Get(path)
}

/**** UPDATE ****/

// type RecordUpdate struct {
// 	Id         string  `json:"id"`
// 	Name       string  `json:"name"`
// 	Aspects    Aspects `json:"aspects"`
// 	SourceTag  string  `json:"sourceTag,omitempty"`
// 	AspectName string  `json:"-"`
// 	AspectFile string  `json:"-"`
// }

func cliRecordUpdate(topCmd *kingpin.CmdClause) {
	r := &RecordCreate{}
	c := topCmd.Command("update", "Update an existing record").Action(func(_ *kingpin.ParseContext) error {
		return recordUpdate(r)
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

func recordUpdate(args *RecordCreate) error {
	r := *args
	
	path := recordPath(&r.Id)
	
	f := func() error {
		addAspects(&r)

		body, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			log.Fatal("Error marshalling body. ", err)
		}
		return PutP(path, bytes.NewReader(body), func(obj JsonObjPayload, arr JsonArrPayload) error {
			fmt.Printf("Successfully updated record '%s'\n", r.Id)
			return nil
		})
	}

	if (r.Name == "") {
		// get current 'name' first as it is required
		return GetP(path, func(obj JsonObjPayload, arr JsonArrPayload) error {
			if (obj == nil) {
				return errors.New("No record body found")
			}
			r.Name = obj["name"].(string)
			return f()
		})
	} else {
		return f()
	}
}

/**** DELETE ****/

type RecordDelete struct {
	Id string
	AspectName string
}

func cliRecordDelete(topCmd *kingpin.CmdClause) {
	r := &RecordDelete{}
	c := topCmd.Command("delete", "Delete a record or one of it's aspects").Action(func(_ *kingpin.ParseContext) error {
		return recordDelete(r)
	})
	c.Flag("id", "Record ID").
		Short('i').
		Required().
		StringVar(&r.Id)
	c.Flag("aspect", "Only delete this aspect").
		Short('a').
		StringVar(&r.AspectName)
}

func recordDelete(args *RecordDelete) error {
	path := recordPath(&args.Id)
	if args.AspectName != "" {
		path = path + "/aspects/" + args.AspectName
	}
	return Delete(path)
}

/**** HISTORY ****/

type RecordHistory struct {
	Id string
	EventId string
	Offset int
	Limit int
	PageToken string
}

func cliRecordHistory(topCmd *kingpin.CmdClause) {
	r := &RecordHistory{Offset: -1, Limit: -1}
	c := topCmd.Command("history", "Get a list of all events for a record").Action(func(_ *kingpin.ParseContext) error {
		return recordHistory(r)
	})
	c.Flag("id", "Record ID").
		Short('i').
		Required().
		StringVar(&r.Id)
	c.Flag("event-id", "Only show event wiht event-id").
		Short('e').
		StringVar(&r.EventId)
	c.Flag("offset", "Index of first record retrieved").
		Short('o').
		IntVar(&r.Offset)
	c.Flag("limit", "The maximumm number of records to retrieve").
		Short('l').
		IntVar(&r.Limit)
	c.Flag("pageToken", "Token that identifies the start of a page of results").
		Short('t').
		StringVar(&r.PageToken)
}

func recordHistory(args *RecordHistory) error {
	path := recordPath(&args.Id) + "/history"
	if (args.EventId != "") {
		path = path + "/" + args.EventId
	}

	q := []string{}
	if args.PageToken != "" {
		q = append(q, "pageToken=" + url.QueryEscape(args.PageToken))
	}
	if args.Offset >= 0 {
		q = append(q, "start=" + url.QueryEscape(strconv.Itoa(args.Offset)))
	}
	if args.Limit >= 0 {
		q = append(q, "limit=" + url.QueryEscape(strconv.Itoa(args.Limit)))
	}
	if len(q) > 0 {
		path = path + "?" + strings.Join(q,"&")
	}
	// fmt.Printf("PATH: %s\n", path)
	return Get(path)
}

/**** UTILS ****/

func recordPath(id *string) string {
	path := "/api/v0/registry/records"
	if *skipGateway {
		path = "/v0/records"
	}
	if id != nil {
		path = path + "/" + *id
	}
	return path
}

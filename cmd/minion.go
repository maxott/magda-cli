package cmd

import (
	"fmt"
	"strings"

	"github.com/maxott/magda-cli/pkg/adapter"
	"github.com/maxott/magda-cli/pkg/minion"
	//log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	cmd := App().Command("minion", "Managing minion registration")
	cliMinionList(cmd)
	cliMinionCreate(cmd)
	// cliMinionRead(cmd)
	// cliMinionUpdate(cmd)
	cliMinionDelete(cmd)

}

/**** LIST ****/

func cliMinionList(topCmd *kingpin.CmdClause) {
	cmd := &minion.ListRequest{}
	topCmd.Command("list", "List all minion registration").Action(func(_ *kingpin.ParseContext) error {
		return adapter.ReplyPrinter(minion.ListRaw(cmd, Adapter(), Logger()))
	})
}

/**** CREATE ****/

func cliMinionCreate(topCmd *kingpin.CmdClause) {
	r := &minion.CreateRequest{}
	var aspects string
	var optAspects string

	c := topCmd.Command("create", "Creates a new minion").Action(func(_ *kingpin.ParseContext) error {
		r.Aspects = strings.Split(aspects, ",")
		if optAspects != "" {
			r.OptionalAspects = strings.Split(optAspects, ",")
		}
		if _, err := minion.CreateRaw(r, Adapter(), Logger()); err == nil {
			fmt.Printf("Successfully create minion hook '%s'\n", r.Id)
			return nil
		} else {
			return err
		}
	})
	c.Flag("id", "Minion ID").
		Short('i').
		Required().
		StringVar(&r.Id)
	c.Flag("url", "Callback URL").
		Short('u').
		Required().
		StringVar(&r.Url)
	c.Flag("aspects", "Comma separated aspects to listen for").
		Short('a').
		Required().
		StringVar(&aspects)
	c.Flag("optional-aspects", "Optional comma separated aspects to listen for").
		StringVar(&optAspects)

}

/**** READ ****/

// func cliMinionRead(topCmd *kingpin.CmdClause) {
// 	r := &minion.ReadCmd{}
// 	c := topCmd.Command("read", "Read the content of a minion").Action(func(_ *kingpin.ParseContext) error {
// 		return adapter.ReplyPrinter(minion.ReadRaw(r, Adapter()))
// 	})
// 	c.Flag("id", "Record ID").
// 		Short('i').
// 		Required().
// 		StringVar(&r.Id)
// }

/**** UPDATE ****/

// func cliMinionUpdate(topCmd *kingpin.CmdClause) {
// 	r := &MinionCreate{}
// 	c := topCmd.Command("update", "Update existing minion").Action(func(_ *kingpin.ParseContext) error {
// 		if js, err := adapter.LoadJsonFromFile(r.MinionFile); err != nil {
// 			return fmt.Errorf("failed to load & verify '%s' - %s", r.MinionFile, err)
// 		} else {
// 			cmd := minion.UpdateCmd{
// 				Id: r.Id, Name: r.Name, Minion: js.AsObject(),
// 			}
// 			if _, err := minion.UpdateRaw(&cmd, Adapter()); err == nil {
// 				fmt.Printf("Successfully updated minion '%s'\n", r.Id)
// 				return nil
// 			} else {
// 				return err
// 			}
// 		}
// 	})
// 	c.Flag("name", "Descriptive name").
// 		Short('n').
// 		StringVar(&r.Name)
// 	cliAddMinionCUFlags(r, c)
// }

/**** DELETE ****/

func cliMinionDelete(topCmd *kingpin.CmdClause) {
	r := &minion.DeleteRequest{}
	c := topCmd.Command("delete", "Delete a minion hook").Action(func(_ *kingpin.ParseContext) error {
		if _, err := minion.DeleteRaw(r, Adapter(), Logger()); err == nil {
			fmt.Printf("Successfully deleted hook '%s'\n", r.Id)
			return nil
		} else {
			return err
		}
	})
	c.Flag("id", "Minion ID").
		Short('i').
		Required().
		StringVar(&r.Id)
}

// Not supported

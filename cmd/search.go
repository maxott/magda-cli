package cmd

import (
	"context"

	"github.com/maxott/magda-cli/pkg/adapter"
	"github.com/maxott/magda-cli/pkg/search"
	"gopkg.in/alecthomas/kingpin.v2"
)

func init() {
	cmd := App().Command("search", "Magda full-text search")
	cliSearchDatasets(cmd)
}

/**** SEARCH DATASETS ****/

func cliSearchDatasets(topCmd *kingpin.CmdClause) {
	r := &search.DatasetRequest{Offset: -1, Limit: -1}
	c := topCmd.Command("datasets", "Fulltext search dataset datasets").Action(func(_ *kingpin.ParseContext) error {
		if pyld, err := search.DatasetRaw(context.Background(), r, Adapter(), Logger()); err != nil {
			return err
		} else {
			return adapter.ReplyPrinter(pyld, *useYaml)
		}
	})
	c.Flag("query", "Full text search query").
		Short('q').
		StringVar(&r.Query)
	c.Flag("offset", "Index of first dataset retrieved").
		Short('o').
		IntVar(&r.Offset)
	c.Flag("limit", "The maximumm number of datasets to retrieve").
		Short('l').
		IntVar(&r.Limit)
	c.Flag("publisher", "Filter search query by names of organisations").
		Short('p').
		StringVar(&r.Publisher)
}

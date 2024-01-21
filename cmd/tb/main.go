package main

import (
	"context"
	"os"

	"github.com/trichner/toolbox/cmd/csv2json"
	"github.com/trichner/toolbox/cmd/sheet2json"
	"github.com/trichner/toolbox/pkg/cmdreg"

	"github.com/trichner/toolbox/cmd/sql2json"

	"github.com/trichner/toolbox/cmd/jiracli"
	"github.com/trichner/toolbox/cmd/json2sheet"
	"github.com/trichner/toolbox/cmd/kraki"
)

func main() {
	r := cmdreg.New(cmdreg.WithProgramName("tb"))

	r.RegisterFunc("csv2json", csv2json.Exec)
	r.RegisterFunc("jiracli", jiracli.Exec)
	r.RegisterFunc("json2sheet", json2sheet.Exec)
	r.RegisterFunc("kraki", kraki.Exec)
	r.RegisterFunc("sheet2json", sheet2json.Exec, cmdreg.WithCompletion(sheet2json.Completions()))
	r.RegisterFunc("sql2json", sql2json.Exec)

	r.RegisterFunc("help", help(r))

	ctx := context.Background()
	r.Exec(ctx, os.Args)
}

func help(r *cmdreg.CommandRegistry) cmdreg.CommandFunc {
	return func(_ context.Context, args []string) {
		r.PrintHelp(os.Stdout)
	}
}

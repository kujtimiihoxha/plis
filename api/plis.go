package api

import (
	"github.com/kujtimiihoxha/plis/cmd"
	"github.com/spf13/cobra"
)

type PlisAPI struct {
	cmd *cobra.Command
}

func (p *PlisAPI) Help() {
	p.cmd.Help()
}
func NewPlisAPI(cmd *cobra.Command) *PlisAPI {
	return &PlisAPI{
		cmd: cmd,
	}
}

func (p *PlisAPI) RunPlisCmd(pCmd string, s []string) error {
	c := []string{
		pCmd,
	}
	c = append(c, s...)
	cmd.RootCmd.SetArgs(c)
	return cmd.RootCmd.Execute()
}

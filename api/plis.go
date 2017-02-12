package api

import (
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

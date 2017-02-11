package api

import (
	"github.com/spf13/cobra"
)

type PlisApi struct {
	cmd *cobra.Command
}

func (p *PlisApi) Help() {
	p.cmd.Help()
}
func NewPlisApi(cmd *cobra.Command) *PlisApi {
	return &PlisApi{
		cmd: cmd,
	}
}

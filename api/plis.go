package api

import (
	"github.com/kujtimiihoxha/plis/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
func (p *PlisAPI) ForceOverride(b bool) {
	viper.Set("plis.fs.force_rewrite", b)
}
func (p *PlisAPI) RunPlisTool(tool string, s []string) error {
	c := []string{
		tool,
	}
	c = append(c, s...)
	cmd.RootCmd.SetArgs(c)
	return cmd.RootCmd.Execute()
}

package api

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
	"testing"

	"github.com/kujtimiihoxha/plis/cmd"
	"strings"
)

type testWriter struct {
	test string
}

func TestPlisAPI_Help(t *testing.T) {
	c := &cobra.Command{
		Use:   "plis",
		Short: "Test description",
	}
	tW := NewTestWritter()
	c.SetOutput(tW)
	plisAPI := NewPlisAPI(c)
	plisAPI.Help()
	Convey("Check if help command is executed", t, func() {
		So(strings.Contains(tW.test, c.Short), ShouldBeTrue)
	})

}
func TestPlisAPI_RunPlisCmd(t *testing.T) {
	c := &cobra.Command{
		Use:   "test",
		Short: "Test description",
	}
	c.Flags().String("tst", "", "Test Fla")
	tW := NewTestWritter()
	c.SetOutput(tW)
	cmd.RootCmd.SetOutput(tW)
	cmd.RootCmd.AddCommand(c)
	plisAPI := NewPlisAPI(cmd.RootCmd)
	args := []string{
		"--tst",
		"hello",
	}
	Convey("Check if command is called", t, func() {
		err := plisAPI.RunPlisCmd("test", args)
		Convey("Check if command executes without errors", func() {
			So(err, ShouldBeNil)
		})
		Convey("Check if command executes and flags are set", func() {
			So(c.Flags().Lookup("tst").Value.String(), ShouldEqual, "hello")
		})
	})

}
func NewTestWritter() *testWriter {
	return &testWriter{
		test: "",
	}
}
func (w *testWriter) Write(data []byte) (n int, err error) {
	w.test += string(data)
	return len(data), nil
}

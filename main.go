package main

import (
	"github.com/agui2200/wingman-store/cmd"
	"github.com/spf13/cobra"
)

func main() {
	c := &cobra.Command{Use: "store"}
	c.AddCommand(
		cmd.InitCommand(),
		cmd.GenerateCommand(),
	)
	_ = c.Execute()
}

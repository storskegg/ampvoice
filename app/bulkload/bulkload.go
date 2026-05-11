package bulkload

import (
    "fmt"

    "github.com/spf13/cobra"
)

func Run() error {
    return cmdRoot.Execute()
}

func runRootE(cmd *cobra.Command, args []string) error {
    return fmt.Errorf("not implemented")
}

package bulkload

import "github.com/spf13/cobra"

var cmdRoot = &cobra.Command{
    Use:   "bulkload",
    Short: "Bulk load data into the database",
    Long:  "Bulk load data into the database from a CSV file",
    RunE:  runRootE,
}

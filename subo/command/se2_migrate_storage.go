package command

import (
    "context"
    "strings"

    "github.com/spf13/cobra"

    "github.com/suborbital/subo/subo/util"
)

// SE2MigrateStorageCommand returns the dev command.
func SE2MigrateStorageCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "migrate [storage root]",
        Short: "reformat storage to match SE2 specification",
        Long:  `Update storage format to match the Suborbital Extension Engine (SE2) storage specification`,
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            location := args[0]

            switch scheme(strings.ToLower(location)) {
            case "":
                return util.FileStoreMigration(location)
            case "gs":
                fallthrough
            case "s3":
                return util.BlobStoreMigration(context.Background(), location)
            }

            return nil
        },
    }

    return cmd
}

func scheme(location string) string {
    idx := strings.Index(location, "://")
    if idx >= 2 {
        return strings.ToLower(location[:strings.Index(location, "://")])
    }

    return ""
}

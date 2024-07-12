package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func main() {
	var (
		newVersion string
		channel    string
		veneerFile string
	)

	rootCmd := &cobra.Command{
		Use:   "veneer-update",
		Short: "Prints an updated veneer file which adds a minor or a patch version to the entries list",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				b   []byte
				err error
			)

			b, err = updateVeneer(veneerFile, channel, newVersion)
			if err != nil {
				panic(err)
			}

			fmt.Println(string(b))
		},
	}

	rootCmd.Flags().StringVar(&newVersion, "new-version", "", "New version (e.g. 0.10.0)")
	rootCmd.Flags().StringVar(&channel, "channel", "", "Channel to update (e.g. stable-v0)")
	rootCmd.Flags().StringVar(&veneerFile, "veneer-file", "", "Path to veneer file")

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func updateVeneer(veneerFile, channel, newVersionStr string) ([]byte, error) {
	var t EverestBasicTemplate
	err := t.readFromFile(veneerFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	var r release
	err = r.create(t.currentVersion(channel), newVersionStr)
	if err != nil {
		return nil, fmt.Errorf("%s: invalid version format: %w", newVersionStr, err)
	}

	err = t.update(r, channel)
	if err != nil {
		return nil, fmt.Errorf("failed to update template: %w", err)
	}

	return t.toByteArray()
}

package commands

import (
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

func Man() *core.Command {
	manCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "man",
			Aliases: []string{"manpages"},
			Short:   "Generate manpages for ionosctl",
			Long: `WARNING: This command is only supported on Linux.

The 'man' command allows you to generate manpages for ionosctl in a given directory. By default, the manpages will be compressed using gzip, but you can skip this step by using the '--skip-compression' flag.
In order to install the manpages, there are a few steps you need to follow:
- Decide where you would like to install the manpages. You can check which directories are available to you by running 'manpath'. If you want to install the manpages to a directory that is not listed, you can add a new entry to '~/.manpath' (see 'man 5 manpath' for how to do it). The directory must contain subdirectories for each section (e.g. 'man1', 'man5', etc.).
- Copy the manpages to the installation directory, in the 'man1' section.
- Run 'sudo mandb' to update the 'man' internal database.

After following these steps, you should be able to use 'man ionosctl' to access the manpages.`,
			TraverseChildren: true,
			Example:          `ionosctl man --target-dir /tmp/ionosctl-man`,
			PreRunE: func(cmd *cobra.Command, args []string) error {
				if runtime.GOOS != "linux" {
					return fmt.Errorf("manpages generation is only supported on Linux")
				}

				targetDir, _ := cmd.Flags().GetString(constants.FlagTargetDir)
				if !filepath.IsAbs(targetDir) {
					return fmt.Errorf("target-dir must be an absolute path")
				}

				return nil
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				targetDir, _ := cmd.Flags().GetString(constants.FlagTargetDir)
				skipCompression, _ := cmd.Flags().GetBool(constants.FlagSkipCompression)

				_, _ = fmt.Fprintf(cmd.OutOrStdout(), jsontabwriter.GenerateVerboseOutput("Checking if target directory for generation already exists"))
				if err := handleExistingManpagesTargetDir(cmd, targetDir); err != nil {
					return err
				}

				if err := os.MkdirAll(targetDir, 0700); err != nil {
					return fmt.Errorf("error creating target directory %s: %w", targetDir, err)
				}

				_, _ = fmt.Fprintf(cmd.OutOrStdout(), jsontabwriter.GenerateVerboseOutput("Generating manpages"))
				if err := doc.GenManTree(cmd.Root(), nil, targetDir); err != nil {
					return fmt.Errorf("error generating manpages: %v", err)
				}

				if skipCompression {
					_, _ = fmt.Fprintf(cmd.OutOrStdout(), jsontabwriter.GenerateLogOutput("Manpages successfully generated."))
					return nil
				}

				if err := compressManpages(targetDir); err != nil {
					return fmt.Errorf("error compressing manpages: %v", err)
				}

				_, _ = fmt.Fprintf(cmd.OutOrStdout(), jsontabwriter.GenerateLogOutput("Manpages successfully generated and compressed."))
				return nil
			},
		},
	}

	manCmd.Command.Flags().String(constants.FlagTargetDir, "/tmp/ionosctl-man", "Target directory where manpages will be generated. Must be an absolute path")
	manCmd.Command.Flags().Bool(constants.FlagSkipCompression, false, "Skip compressing manpages with gzip, just generate them")

	return manCmd
}

func handleExistingManpagesTargetDir(c *cobra.Command, targetDir string) error {
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return nil
	}

	if !confirm.FAsk(
		c.InOrStdin(),
		fmt.Sprintf("Target directory %s already exists. Do you want to replace it", targetDir),
		viper.GetBool(constants.ArgForce),
	) {
		return fmt.Errorf(confirm.UserDenied)
	}

	if err := os.RemoveAll(targetDir); err != nil {
		return fmt.Errorf("error deleting target directory %s: %w", targetDir, err)
	}

	return nil
}

func compressManpages(genDir string) error {
	files, err := os.ReadDir(genDir)
	if err != nil {
		return fmt.Errorf("error opening manpages target directory %s: %v", genDir, err)
	}

	for _, file := range files {
		uncompressedFilePath := fmt.Sprintf("%s/%s", genDir, file.Name())
		fileContent, err := os.ReadFile(uncompressedFilePath)
		if err != nil {
			return fmt.Errorf("error reading uncompressed manpage file %s: %v", file.Name(), err)
		}

		compressedFilePath := fmt.Sprintf("%s.gz", uncompressedFilePath)
		if err = gzipManFile(fileContent, compressedFilePath); err != nil {
			return fmt.Errorf("error compressing manpage file: %v", err)
		}

		if err = os.Remove(uncompressedFilePath); err != nil {
			return fmt.Errorf("error removing uncompressed manpage file %s: %v", file.Name(), err)
		}
	}

	return nil
}

func gzipManFile(fileContent []byte, newFileName string) error {
	gzipFile, err := os.OpenFile(newFileName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error creating gzipped manpage file %s: %v", newFileName, err)
	}
	defer gzipFile.Close()

	gzipWriter := gzip.NewWriter(gzipFile)
	_, err = gzipWriter.Write(fileContent)
	if err != nil {
		return fmt.Errorf("error writing to gzipped manpage file %s: %v", newFileName, err)
	}
	defer gzipWriter.Close()

	return nil
}

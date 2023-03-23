// Package doc generates Markdown files and organizes a directory structure that follows the command hierarchy.
//
// The WriteDocs function is the main entry point for generating the documentation. It recursively processes all
// subcommands and creates the appropriate files and directories based on the command structure, following this rule:
//
// - For commands with no namespace (e.g., `ionosctl version`, `ionosctl login`), files are placed in `docs/subcommands/cli-setup`
// - For commands with a one-level deep namespace (e.g., `ionosctl server list`, `ionosctl datacenter list`), files are placed in `docs/subcommands/compute`
// - For commands with deeper namespaces (e.g., `ionosctl dbaas mongo cluster create`), files are placed in corresponding subdirectories (e.g., `docs/subcommands/dbaas/mongo/cluster`)
//
// The GenerateSummary function is another entry point, which can create a summary.md file containing the table of contents for the generated documentation.
package doc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

const rootCmdName = "ionosctl"

func GenerateSummary(dir string) error {
	f, err := os.Create(filepath.Join(dir, "summary.md"))
	if err != nil {
		return err
	}
	defer f.Close()

	buf := new(bytes.Buffer)

	buf.WriteString("# Table of contents\n\n")
	buf.WriteString("* [Introduction](README.md)\n")
	buf.WriteString("* [Changelog](/CHANGELOG.md)\n\n")
	buf.WriteString("## Subcommands\n\n")

	err = generateDirectoryContent(filepath.Join(dir, "subcommands"), buf, "")
	if err != nil {
		return err
	}
	_, err = buf.WriteTo(f)
	return err
}

func WriteDocs(cmd *core.Command, dir string) error {
	for _, c := range cmd.SubCommands() {
		if c.Command.HasParent() {
			if !c.Command.IsAvailableCommand() {
				continue
			}
		}
		if err := WriteDocs(c, dir); err != nil {
			return err
		}
	}
	if err := createStructure(cmd, dir); err != nil {
		return err
	}
	return nil
}

func createStructure(cmd *core.Command, dir string) error {
	var file, filename string
	if cmd != nil {
		if cmd.Command.HasParent() && cmd.Command.Runnable() {
			name := strings.ReplaceAll(cmd.Command.CommandPath(), rootCmdName+" ", "")
			name = strings.ReplaceAll(name, " ", "-")
			filename = fmt.Sprintf("%s.md", name)
			subdir := determineSubdir(name)
			dir = filepath.Join(dir, subdir)
		} else {
			return nil
		}
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
		file = filepath.Join(dir, filename)
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		defer f.Close()
		err = writeDoc(cmd, f)
		if err != nil {
			return err
		}
	}
	return nil
}

// determineSubdir is a hack to respect the current doc tree style
func determineSubdir(name string) string {
	segments := strings.Split(name, "-")
	switch len(segments) {
	case 0, 1:
		return "cli-setup"
	case 2:
		return "compute"
	default:
		return filepath.Join(segments[0], strings.Join(segments[1:], "/"))
	}
}

func generateDirectoryContent(dir string, buf *bytes.Buffer, prefix string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		name := file.Name()
		title := strings.ToTitle(strings.ReplaceAll(name, "-", " "))

		if file.IsDir() {
			subdir := filepath.Join(dir, name)
			buf.WriteString(fmt.Sprintf("%s* %s\n", prefix, title))
			err = generateDirectoryContent(subdir, buf, prefix+"    ")
			if err != nil {
				return err
			}
			continue
		}

		if filepath.Ext(name) == ".md" {
			nameWithoutExt := strings.TrimSuffix(name, filepath.Ext(name))
			title = strings.ToTitle(strings.ReplaceAll(nameWithoutExt, "-", " "))
			link := filepath.Join("subcommands", strings.ReplaceAll(dir, "\\", "/"), name)
			buf.WriteString(fmt.Sprintf("%s* [%s](%s)\n", prefix, title, link))
		}
	}
	return nil
}

func writeDoc(cmd *core.Command, w io.Writer) error {
	cmd.Command.InitDefaultHelpCmd()
	cmd.Command.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.Command.CommandPath()

	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("description: %s\n", cmd.Command.Short))
	buf.WriteString("---\n\n")

	// Customize title
	title := strings.Title(strings.ReplaceAll(cmd.Command.CommandPath(), rootCmdName+" ", ""))
	title = customizeTitle(title, "-", "")
	title = customizeTitle(title, " ", "")

	buf.WriteString(fmt.Sprintf("# %s\n\n", title))

	buf.WriteString("## Usage\n\n")
	if cmd.Command.Runnable() {
		buf.WriteString(fmt.Sprintf("```text\n%s\n```\n\n", cmd.Command.UseLine()))
	}
	if cmd.Command.HasAvailableSubCommands() {
		buf.WriteString(fmt.Sprintf("```text\n%s [command]\n```\n\n", cmd.Command.CommandPath()))
	}

	if len(cmd.Command.Aliases) > 0 || len(cmd.Command.Parent().Aliases) > 0 {
		buf.WriteString("## Aliases\n\n")
		// Write available aliases for all 3 levels of Command
		if cmd.Command.Parent().Parent() != nil {
			writeCmdAliases(&core.Command{Command: cmd.Command.Parent().Parent()}, buf)
		}
		if cmd.Command.Parent() != nil {
			writeCmdAliases(&core.Command{Command: cmd.Command.Parent()}, buf)
		}
		writeCmdAliases(cmd, buf)
	}

	buf.WriteString("## Description\n\n")
	if len(cmd.Command.Long) > 0 {
		buf.WriteString(cmd.Command.Long + "\n\n")
	} else if len(cmd.Command.Short) > 0 {
		buf.WriteString(cmd.Command.Short + "\n\n")
	}

	flags := cmd.Command.Flags()
	if flags.HasAvailableFlags() {
		buf.WriteString("## Options\n\n```text\n")
		flags.SortFlags = true
		// create new buffer to replace user info
		newbuf := new(bytes.Buffer)
		flags.SetOutput(newbuf)
		flags.PrintDefaults()
		// get $XDG_CONFIG_HOME from environment
		xdgConfig, _ := os.UserConfigDir()
		// replace with constant $XDG_CONFIG_HOME
		buf.Write(bytes.ReplaceAll(newbuf.Bytes(), []byte(xdgConfig), []byte("$XDG_CONFIG_HOME")))
		buf.WriteString("```\n\n")
	}

	if len(cmd.Command.Example) > 0 {
		buf.WriteString("## Examples\n\n")
		buf.WriteString(fmt.Sprintf("```text\n%s\n```\n\n", cmd.Command.Example))
	}

	var link string
	if hasSeeAlso(cmd) {
		children := cmd.Command.Commands()
		buf.WriteString("## Related commands\n\n")
		buf.WriteString("| Command | Description |\n")
		buf.WriteString("| :--- | :--- |\n")
		for _, child := range children {
			if !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			cname := name + " " + child.Name()
			if !child.HasAvailableSubCommands() {
				link = child.Name() + ".md"
			} else {
				link = child.Name() + "/"
			}
			buf.WriteString(fmt.Sprintf("| [%s](%s) | %s |\n", cname, link, child.Short))
		}
		buf.WriteString("\n")
	}

	_, err := buf.WriteTo(w)
	return err
}

func writeCmdAliases(cmd *core.Command, buf *bytes.Buffer) {
	if cmd != nil {
		if len(cmd.Command.Aliases) > 0 {
			buf.WriteString(fmt.Sprintf("For `%s` command:\n\n", cmd.Command.Name()))
			buf.WriteString(fmt.Sprintf("```text\n%s\n```\n\n", cmd.Command.Aliases))
		}
	}
	return
}

func customizeTitle(title, old, new string) string {
	if strings.Contains(title, old) {
		title = strings.ReplaceAll(title, old, new)
	}
	return title
}

func hasSeeAlso(cmd *core.Command) bool {
	if cmd.Command.HasParent() && cmd.Command.HasAvailableSubCommands() {
		return true
	}
	for _, c := range cmd.SubCommands() {
		if !c.Command.IsAvailableCommand() || c.Command.IsAdditionalHelpTopicCommand() {
			continue
		}
		return true
	}
	return false
}

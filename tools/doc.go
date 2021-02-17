package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ionos-cloud/ionosctl/commands"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/spf13/cobra"
)

func main() {
	dir := os.Getenv("DOCS_OUT")
	if dir == "" {
		fmt.Printf("DOCS_OUT environment variable not set.\n")
		os.Exit(1)
	}
	if _, err := os.Stat(dir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err := WriteDocs(commands.GetRootCmd(), dir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	for _, cmd := range commands.GetRootCmd().ChildCommands() {
		err := WriteDocs(cmd, dir)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}

// Generate Markdown documentation based on information described in commands.
// Using WriteDocs function, it will be created one structure in the path specified.
// For example, for `ionosctl volume` commands, we will have the following structure generated:
// $DOCS_OUT/volume/
// ├── attach
// │   ├── get.md
// │   ├── list.md
// │   └── README.md
// ├── create.md
// ├── delete.md
// ├── detach.md
// ├── get.md
// ├── list.md
// ├── README.md
// └── update.md
// For each command, an Markdown file is generated with the following fields:
// # Usage
// # Description
// # Options
// # Examples
// # See also
// depending if these fields are set in the command.

const rootCmdName = "ionosctl"

func WriteDocs(cmd *builder.Command, dir string) error {
	// Exit if there's an error
	for _, c := range cmd.ChildCommands() {
		if c.Command.HasParent() {
			if !c.Command.IsAvailableCommand() || c.Command.Parent().Name() != rootCmdName {
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

func createStructure(cmd *builder.Command, dir string) error {
	var file string
	if cmd.Command.HasAvailableSubCommands() && cmd.Command.Name() != rootCmdName {
		dirParent := fmt.Sprintf("%s/%s", dir, cmd.Command.Use)
		err := os.Mkdir(dirParent, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			return err
		}
		// Add a README.md file about the parent cmd
		file = filepath.Join(dirParent, "README.md")
		f, err := os.Create(file)
		if err != nil {
			return err
		}
		defer f.Close()

		err = writeDoc(cmd, f)
		if err != nil {
			return err
		}
		for _, child := range cmd.ChildCommands() {
			if !child.Command.IsAvailableCommand() || child.Command.IsAdditionalHelpTopicCommand() {
				continue
			}
			if err := createStructure(child, dirParent); err != nil {
				return err
			}
		}
	} else {
		if cmd.Command.Name() == rootCmdName {
			file = filepath.Join(dir, "README.md")
		} else {
			file = filepath.Join(dir, fmt.Sprintf("%s.md", cmd.Command.Name()))
		}
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

func writeDoc(cmd *builder.Command, w io.Writer) error {
	cmd.Command.InitDefaultHelpCmd()
	cmd.Command.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.Command.CommandPath()

	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("description: %s\n", cmd.Command.Short))
	buf.WriteString("---\n\n")

	buf.WriteString(fmt.Sprintf("# %s\n\n", strings.Title(cmd.Command.Name())))

	buf.WriteString("## Usage\n\n")
	if cmd.Command.Runnable() {
		buf.WriteString(fmt.Sprintf("```text\n%s\n```\n\n", cmd.Command.UseLine()))
	}
	if cmd.Command.HasAvailableSubCommands() {
		buf.WriteString(fmt.Sprintf("```text\n%s [command]\n```\n\n", cmd.Command.CommandPath()))
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
		buf.WriteString("## See also\n\n")
		if cmd.Command.HasParent() {
			parent := cmd.Command.Parent()
			pname := parent.CommandPath()
			if !cmd.Command.HasAvailableSubCommands() {
				link = "./"
			} else {
				link = "../"
			}
			buf.WriteString(fmt.Sprintf("* [%s](%s)\n", pname, link))
			cmd.Command.VisitParents(func(c *cobra.Command) {
				if c.DisableAutoGenTag {
					cmd.Command.DisableAutoGenTag = c.DisableAutoGenTag
				}
			})
		}

		children := cmd.Command.Commands()
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
			buf.WriteString(fmt.Sprintf("* [%s](%s)\n", cname, link))
		}
		buf.WriteString("\n")
	}

	_, err := buf.WriteTo(w)
	return err
}

func hasSeeAlso(cmd *builder.Command) bool {
	if cmd.Command.HasParent() {
		return true
	}
	for _, c := range cmd.ChildCommands() {
		if !c.Command.IsAvailableCommand() || c.Command.IsAdditionalHelpTopicCommand() {
			continue
		}
		return true
	}
	return false
}

// Package doc generates Markdown files and organizes a directory structure that follows the command hierarchy.
//
// The WriteDocs function is the main entry point for generating the documentation. It recursively processes all
// subcommands and creates the appropriate files and directories based on the command structure,
// following the rules defined in determineSubdir
//
// The GenerateSummary function is another entry point, which can create a summary.md file containing the table of contents for the generated documentation.
package doc

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/pflag"
)

const rootCmdName = "ionosctl"

type SubdirRule struct {
	// Either a static prefix (e.g. []string{"dbaas","in-memory-db"})
	// or a glob/regex.
	Prefix []string

	// How to build the target path.
	// Use “{1}”, “{2+}” etc to refer to prefix and tail segments.
	// e.g. "Database-as-a-Service/{1}/replicas/{2+}"
	//  - {i}   → the i’th segment
	//  - {i+}  → segments[i]…end, joined by filepath.Separator
	Template string
}

var subdirRules = []SubdirRule{
	// CLI Setup
	{Prefix: []string{"version"}, Template: "CLI Setup/{0+}"},
	{Prefix: []string{"completion"}, Template: "CLI Setup/{0+}"},
	{Prefix: []string{"man"}, Template: "CLI Setup/{0+}"},
	{Prefix: []string{"login"}, Template: "CLI Setup/{0+}"},

	// Authentication
	{Prefix: []string{"token"}, Template: "Authentication/{0+}"},

	// Interactive Shell
	{Prefix: []string{"shell"}, Template: "Interactive Shell/{0+}"},

	// Database-as-a-Service / In-Memory-DB
	// cmdPathHyphens: "dbaas-in-memory-db-…"
	// segs:          ["dbaas","in","memory","db",…]
	{Prefix: []string{"dbaas", "in", "memory", "db"},
		Template: "Database-as-a-Service/In-Memory-DB/{4+}"},

	// All other DBaaS
	{Prefix: []string{"dbaas"}, Template: "Database-as-a-Service/{1+}"},

	// Non-compute APIs
	{Prefix: []string{"apigateway"}, Template: "API Gateway/{1+}"},
	{Prefix: []string{"applicationloadbalancer"}, Template: "Application-Load-Balancer/{1+}"},
	{Prefix: []string{"backupunit"}, Template: "Managed-Backup/{1+}"},
	// certificate-manager → ["certificate","manager"]
	{Prefix: []string{"certmanager"}, Template: "Certificate-Manager/{2+}"},
	// container-registry → ["container","registry"]
	{Prefix: []string{"container", "registry"}, Template: "Container-Registry/{2+}"},
	{Prefix: []string{"natgateway"}, Template: "NAT-Gateway/{1+}"},
	{Prefix: []string{"networkloadbalancer"}, Template: "Network-Load-Balancer/{1+}"},
	{Prefix: []string{"k8s"}, Template: "Managed-Kubernetes/{1+}"},
	{Prefix: []string{"user"}, Template: "User-Management/{1+}"},
	{Prefix: []string{"dns"}, Template: "DNS/{1+}"},
	{Prefix: []string{"cdn"}, Template: "CDN/{1+}"},
	{Prefix: []string{"kafka"}, Template: "Kafka/{1+}"},
	{Prefix: []string{"config"}, Template: "CLI Setup/{1+}"},
	{Prefix: []string{"vm", "autoscaling"}, Template: "VM Autoscaling/{2+}"},
	{Prefix: []string{"vpn"}, Template: "VPN Gateway/{1+}"},
	{Prefix: []string{"logging", "service"}, Template: "Logging-Service/{2+}"},

	// fallback → Compute Engine
	{Prefix: []string{}, Template: "Compute Engine/{0+}"},
}

// determineSubdir is a hack to support the old tree structure...
func determineSubdir(cmdPathHyphens string) string {
	segs := strings.Split(cmdPathHyphens, "-")

	// 2) Try each rule in order
	for _, rule := range subdirRules {
		if len(segs) < len(rule.Prefix) {
			continue
		}
		// Does segs[0:len(rule.Prefix)] == rule.Prefix?
		match := true
		for i, p := range rule.Prefix {
			if segs[i] != p {
				match = false
				break
			}
		}
		if !match {
			continue
		}

		return fillTemplate(rule.Template, segs)
	}

	return filepath.Join("Compute Engine", filepath.Join(segs...))
}

func fillTemplate(tmpl string, segs []string) string {
	tmpl = regexp.MustCompile(`\{(\d+)\+\}`).
		ReplaceAllStringFunc(tmpl, func(m string) string {
			idx, _ := strconv.Atoi(m[1 : len(m)-2])
			return strings.Join(segs[idx:], string(filepath.Separator))
		})
	tmpl = regexp.MustCompile(`\{(\d+)\}`).
		ReplaceAllStringFunc(tmpl, func(m string) string {
			idx, _ := strconv.Atoi(m[1 : len(m)-1])
			return segs[idx]
		})
	return tmpl
}

func GenerateSummary(dir string) error {
	f, err := os.Create(filepath.Join(dir, "summary.md"))
	if err != nil {
		return err
	}
	defer f.Close()

	buf := new(bytes.Buffer)

	_, err = buf.WriteString("# Table of contents\n\n* [Introduction](README.md)\n* [Changelog](/CHANGELOG.md)\n\n## Subcommands\n\n")
	if err != nil {
		return err
	}

	err = generateDirectoryContent(filepath.Join(dir, "subcommands"), buf, "")
	if err != nil {
		return err
	}

	buf.WriteString("\n\n## Legal\n\n---\n\n* [Privacy policy](https://www.ionos.com/terms-gtc/terms-privacy/)\n* [Imprint](https://www.ionos.de/impressum)\n")
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(f)
	if err != nil {
		return err
	}

	return nil
}

func generateDirectoryContent(dir string, buf *bytes.Buffer, prefix string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		name := file.Name()

		if file.IsDir() {
			subdir := filepath.Join(dir, name)
			buf.WriteString(fmt.Sprintf("%s* %s\n", prefix, strings.ReplaceAll(name, "-", " ")))
			err = generateDirectoryContent(subdir, buf, prefix+"    ")
			if err != nil {
				return err
			}
			continue
		}

		if filepath.Ext(name) == ".md" {
			nameWithoutExt := strings.TrimSuffix(name, filepath.Ext(name))
			title := strings.ReplaceAll(nameWithoutExt, "-", " ")
			link := filepath.Join("subcommands", strings.ReplaceAll(strings.TrimPrefix(dir, "docs/subcommands/"), "\\", "/"), file.Name())
			escapedLink := url.PathEscape(link)
			buf.WriteString(fmt.Sprintf("%s* [%s](%s)\n", prefix, title, escapedLink))
		}
	}
	return nil
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
			cmdNameWithHyphens := strings.ReplaceAll(name, " ", "-")
			subdir := determineSubdir(cmdNameWithHyphens)
			dir = filepath.Join(dir, subdir)
		} else {
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(dir), os.ModePerm); err != nil {
			return err
		}
		filename = filepath.Base(dir) + ".md"
		file = filepath.Join(filepath.Dir(dir), filename)
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

func writeDoc(cmd *core.Command, w io.Writer) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic occurred for command path %s: %+v\n", cmd.Command.CommandPath(), r)
		}
	}()

	cmd.Command.InitDefaultHelpCmd()
	cmd.Command.InitDefaultHelpFlag()

	buf := new(bytes.Buffer)
	name := cmd.Command.CommandPath()

	buf.WriteString("---\n")
	buf.WriteString(fmt.Sprintf("description: \"%s\"\n", cmd.Command.Short))
	buf.WriteString("---\n\n")

	// Customize title
	title := strings.Title(strings.ReplaceAll(cmd.Command.CommandPath(), rootCmdName+" ", ""))
	title = StrReplaceIfContains(title, "-", "")
	title = StrReplaceIfContains(title, " ", "")

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

		// Create new buffer to replace user info
		newbuf := new(bytes.Buffer)
		flags.SetOutput(newbuf)
		flags.VisitAll(
			func(flag *pflag.Flag) {
				handler := getStrategyForFlag(flag.Name)
				// If a custom default value handler is specified, use it to modify the default of this flag for docs
				if handler != nil {
					flag.DefValue = handler(flag.Usage, flag.DefValue)
				}
			},
		)
		flags.PrintDefaults()

		// Get $XDG_CONFIG_HOME from environment
		xdgConfig, _ := os.UserConfigDir()
		// Replace with constant $XDG_CONFIG_HOME
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

func StrReplaceIfContains(title, old, new string) string {
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

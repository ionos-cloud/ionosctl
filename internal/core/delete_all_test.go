package core

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// scriptedConfirmer answers each FAsk call from a queue of booleans.
// A true override (e.g. --force) short-circuits to true without consuming an answer.
type scriptedConfirmer struct{ answers []bool }

func (s *scriptedConfirmer) Ask(_ io.Reader, _ string, overrides ...bool) bool {
	for _, o := range overrides {
		if o {
			return true
		}
	}
	if len(s.answers) == 0 {
		return false
	}
	a := s.answers[0]
	s.answers = s.answers[1:]
	return a
}

func newTestCmdConfig(out io.Writer) *CommandConfig {
	c := &CommandConfig{
		Command: &Command{
			Command: &cobra.Command{Use: "test"},
		},
		NS:       "test",
		Resource: "thing",
	}
	c.Command.Command.SetOut(out)
	c.Command.Command.SetIn(strings.NewReader(""))
	return c
}

func TestDeleteAll(t *testing.T) {
	type res struct {
		id   string
		name string
	}
	all := []res{{"a", "alpha"}, {"b", "beta"}, {"c", "gamma"}}

	baseOpts := func(deleted *[]string, failIDs map[string]bool) DeleteAllOptions[res] {
		return DeleteAllOptions[res]{
			Resource: "thing",
			List:     func() ([]res, error) { return all, nil },
			Summary:  func(r res) string { return fmt.Sprintf("%s (id: %s)", r.name, r.id) },
			ID:       func(r res) string { return r.id },
			Delete: func(r res) error {
				if failIDs[r.id] {
					return fmt.Errorf("boom %s", r.id)
				}
				*deleted = append(*deleted, r.id)
				return nil
			},
		}
	}

	t.Run("empty list errors", func(t *testing.T) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		c := newTestCmdConfig(&bytes.Buffer{})
		var deleted []string
		opts := baseOpts(&deleted, nil)
		opts.List = func() ([]res, error) { return nil, nil }
		err := DeleteAll(c, opts)
		assert.Error(t, err)
		assert.Empty(t, deleted)
	})

	t.Run("list error propagates", func(t *testing.T) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		c := newTestCmdConfig(&bytes.Buffer{})
		var deleted []string
		opts := baseOpts(&deleted, nil)
		opts.List = func() ([]res, error) { return nil, fmt.Errorf("list failed") }
		err := DeleteAll(c, opts)
		assert.Error(t, err)
	})

	t.Run("force deletes all without prompting", func(t *testing.T) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		confirm.SetStrategy(&scriptedConfirmer{answers: nil}) // no answers; force must bypass
		defer confirm.SetStrategy(nil)

		c := newTestCmdConfig(&bytes.Buffer{})
		var deleted []string
		err := DeleteAll(c, baseOpts(&deleted, nil))
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b", "c"}, deleted)
	})

	t.Run("answering n skips only that item and continues", func(t *testing.T) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		// yes, no, yes -> delete a, skip b, delete c
		confirm.SetStrategy(&scriptedConfirmer{answers: []bool{true, false, true}})
		defer confirm.SetStrategy(nil)

		c := newTestCmdConfig(&bytes.Buffer{})
		var deleted []string
		err := DeleteAll(c, baseOpts(&deleted, nil))
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "c"}, deleted, "b must be skipped, a and c still deleted")
	})

	t.Run("delete error is aggregated and does not stop the rest", func(t *testing.T) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)

		c := newTestCmdConfig(&bytes.Buffer{})
		var deleted []string
		err := DeleteAll(c, baseOpts(&deleted, map[string]bool{"b": true}))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "b")
		assert.Equal(t, []string{"a", "c"}, deleted, "a and c still deleted despite b failing")
	})

	t.Run("preview lists every resource", func(t *testing.T) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)

		buf := &bytes.Buffer{}
		c := newTestCmdConfig(buf)
		var deleted []string
		_ = DeleteAll(c, baseOpts(&deleted, nil))
		out := buf.String()
		assert.Contains(t, out, "alpha (id: a)")
		assert.Contains(t, out, "beta (id: b)")
		assert.Contains(t, out, "gamma (id: c)")
		assert.Contains(t, out, "Done: 3 deleted, 0 skipped, 0 failed")
	})
}

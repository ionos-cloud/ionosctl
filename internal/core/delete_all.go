package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/spf13/viper"
)

// DeleteAllError is returned by DeleteAll when one or more deletions fail.
// Its message is a concise count so the terminal stays readable; the per-item
// compact reasons are available to programmatic callers via Unwrap.
type DeleteAllError struct {
	Resource     string
	Failed, Total int
	Errs         error // joined one-line reasons, one per failed item
}

func (e *DeleteAllError) Error() string {
	return fmt.Sprintf("%d of %d %s(s) could not be deleted", e.Failed, e.Total, e.Resource)
}

func (e *DeleteAllError) Unwrap() error { return e.Errs }

// shortAPIError collapses an IONOS API error to "STATUS - message; message".
// IONOS services share the {"httpStatus":N,"messages":[{"message":...}]} error
// envelope, so it is parsed out of the error text regardless of which SDK
// produced it. Falls back to whitespace-collapsed, truncated text otherwise.
func shortAPIError(err error) string {
	s := err.Error()
	if i := strings.IndexByte(s, '{'); i >= 0 {
		var env struct {
			HTTPStatus int `json:"httpStatus"`
			Messages   []struct {
				Message string `json:"message"`
			} `json:"messages"`
		}
		if json.Unmarshal([]byte(s[i:]), &env) == nil && len(env.Messages) > 0 {
			msgs := make([]string, len(env.Messages))
			for j, m := range env.Messages {
				msgs[j] = m.Message
			}
			return fmt.Sprintf("%d - %s", env.HTTPStatus, strings.Join(msgs, "; "))
		}
	}
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > 200 {
		s = s[:200] + "…"
	}
	return s
}

// DeleteAllOptions configures a consistent "delete --all" flow for a resource type.
//
// All ionosctl resources should drive their --all deletion through DeleteAll so that
// behaviour is uniform: a preview of what will be deleted, a per-item confirmation that
// skips (never aborts the rest) when the user answers 'n', error aggregation that
// continues on failure, and a per-item plus summary output.
type DeleteAllOptions[T any] struct {
	// Resource is the singular human-readable label used in messages, e.g. "datacenter".
	Resource string
	// List fetches all candidate resources. Apply any --name/location filters here.
	List func() ([]T, error)
	// Summary returns a rich one-line description of a resource for the preview list
	// and the confirmation prompt, e.g. "myDC (id: abc, location: de/txl, desc: prod)".
	Summary func(T) string
	// ID returns the resource identifier, used in error messages.
	ID func(T) string
	// Delete removes a single resource.
	Delete func(T) error
}

// DeleteAll drives the canonical "delete --all" flow:
//
//  1. List candidates; error out if none are found.
//  2. Print a preview of every resource that will be deleted.
//  3. Ask for confirmation per item. Answering 'n' skips ONLY that item and
//     continues with the rest; it never cancels the whole operation. The global
//     --force flag skips all prompts.
//  4. Delete each confirmed item, aggregating errors and continuing on failure.
//  5. Print per-item success and a final summary (deleted / skipped / failed).
//
// It returns the joined error of all failed deletions, or nil if none failed.
func DeleteAll[T any](c *CommandConfig, o DeleteAllOptions[T]) error {
	c.Verbose("Deleting all %ss!", o.Resource)

	items, err := o.List()
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return fmt.Errorf("no %ss found", o.Resource)
	}

	c.Msg("%d %s(s) found to delete:", len(items), o.Resource)
	for _, it := range items {
		c.Msg("  - %s", o.Summary(it))
	}

	force := viper.GetBool(constants.ArgForce)

	var errs error
	var deleted, skipped, failed int
	for _, it := range items {
		if !confirm.FAsk(c.Command.Command.InOrStdin(),
			fmt.Sprintf("delete %s %s", o.Resource, o.Summary(it)), force) {
			skipped++
			continue // skip THIS item only - never abort the remaining items
		}

		if err := o.Delete(it); err != nil {
			failed++
			short := shortAPIError(err)
			c.Msg("✗ %s %s — %s", o.Resource, o.ID(it), short)
			c.Verbose("failed to delete %s %s: %v", o.Resource, o.ID(it), err) // raw detail only with --verbose
			errs = errors.Join(errs, fmt.Errorf("%s %s: %s", o.Resource, o.ID(it), short))
			continue
		}

		deleted++
		c.Msg("✓ %s %s deleted", o.Resource, o.ID(it))
	}

	c.Msg("Done: %d deleted, %d skipped, %d failed", deleted, skipped, failed)
	if failed > 0 {
		return &DeleteAllError{Resource: o.Resource, Failed: failed, Total: len(items), Errs: errs}
	}
	return nil
}

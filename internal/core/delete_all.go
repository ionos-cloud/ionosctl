package core

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/spf13/viper"
)

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
			errs = errors.Join(errs, fmt.Errorf(constants.ErrDeleteAll, o.Resource, o.ID(it), err))
			continue
		}

		deleted++
		c.Msg("%s %s deleted", o.Resource, o.ID(it))
	}

	c.Msg("Done: %d deleted, %d skipped, %d failed", deleted, skipped, failed)
	return errs
}

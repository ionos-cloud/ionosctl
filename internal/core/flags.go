package core

import (
	"strconv"

	"github.com/spf13/pflag"
)

// flags wraps the executing command's pflag set with single-value typed getters.
//
// pflag's GetX methods return (val, err) and are strict: they error if the flag's
// declared type doesn't match the getter (e.g. .Int32 on a flag registered with
// AddIntFlag). The command code, however, was written against viper, which coerced
// freely - reading an int flag as int32, an IP flag as a string, etc. - and never
// surfaced an error. To keep every existing read working unchanged while dropping
// viper, these getters reproduce that lenient behavior:
//
//   - native fast path: if the flag's type matches the getter, return pflag's value.
//   - coercion fallback: on a type mismatch, convert from the flag's own
//     Value.String() (best effort; an unparseable value yields the zero value).
//   - undefined flag: Lookup returns nil, so the zero value is returned (treated as
//     unset) - matching viper, where an unknown key read as zero.
//
// Like viper, these never panic on a read. Changed and every native GetX remain
// available through the embedded *pflag.FlagSet for callers that want strictness.
type flags struct{ *pflag.FlagSet }

func (f flags) String(name string) string {
	if v, err := f.GetString(name); err == nil {
		return v
	}
	if fl := f.Lookup(name); fl != nil {
		return fl.Value.String() // coerce any flag type to its string form
	}
	return ""
}

func (f flags) Bool(name string) bool {
	if v, err := f.GetBool(name); err == nil {
		return v
	}
	if fl := f.Lookup(name); fl != nil {
		b, _ := strconv.ParseBool(fl.Value.String())
		return b
	}
	return false
}

func (f flags) Int(name string) int {
	if v, err := f.GetInt(name); err == nil {
		return v
	}
	if fl := f.Lookup(name); fl != nil {
		n, _ := strconv.Atoi(fl.Value.String())
		return n
	}
	return 0
}

func (f flags) Int32(name string) int32 {
	if v, err := f.GetInt32(name); err == nil {
		return v
	}
	if fl := f.Lookup(name); fl != nil {
		n, _ := strconv.ParseInt(fl.Value.String(), 10, 32)
		return int32(n)
	}
	return 0
}

func (f flags) StringSlice(name string) []string {
	if v, err := f.GetStringSlice(name); err == nil {
		return v
	}
	if fl := f.Lookup(name); fl != nil {
		return []string{fl.Value.String()} // coerce a scalar flag to a single-element slice
	}
	return nil
}

func (f flags) IntSlice(name string) []int {
	if v, err := f.GetIntSlice(name); err == nil {
		return v
	}
	if fl := f.Lookup(name); fl != nil {
		if n, err := strconv.Atoi(fl.Value.String()); err == nil {
			return []int{n}
		}
	}
	return nil
}

func (f flags) StringToString(name string) map[string]string {
	if v, err := f.GetStringToString(name); err == nil {
		return v
	}
	return map[string]string{}
}

// FlagsOf wraps a raw pflag set in the typed, single-value view. Use it where only
// a *cobra.Command is in scope (e.g. inside flag-completion closures) rather than a
// *core.Command/CommandConfig: core.FlagsOf(cmd.Flags()).String(name).
func FlagsOf(fs *pflag.FlagSet) flags { return flags{fs} }

// Flags returns a typed, single-value view over the command's flag set.
func (c *Command) Flags() flags { return flags{c.Command.Flags()} }

// Flags returns a typed, single-value view over the command's flag set.
func (c *CommandConfig) Flags() flags { return flags{c.Command.Command.Flags()} }

// Flags returns a typed, single-value view over the command's flag set.
func (c *PreCommandConfig) Flags() flags { return flags{c.Command.Command.Flags()} }

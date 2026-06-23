package core

import (
	"errors"

	"github.com/spf13/pflag"
)

// flags wraps the executing command's pflag set with single-value typed getters.
//
// pflag's GetX methods return (val, err). At read time that err can only mean one of
// two programmer-side conditions, never bad user input (malformed input fails earlier
// during cobra's parse via flag.Set). See pflag's getFlagType
// (vendor/github.com/spf13/pflag/flag.go):
//
//   - flag not defined: Lookup returns nil. pflag reports this as a typed
//     *pflag.NotExistError. We treat it as "unset" and return the zero value - this
//     matches the old viper-based reads (an unset/unknown key returned zero), so a
//     command can read a flag it didn't set without exploding, and tests need only
//     register the flags they actually inject.
//
//   - type mismatch: flag.Value.Type() != the requested type (e.g. .Bool on a string
//     flag). This is always a real bug - the wrong getter for a flag's declared type -
//     which viper silently masked via loose coercion. We panic to fail loud in dev/CI.
//
// The conversion func (third getFlagType path) parses the flag's own Value.String()
// output, so it round-trips for every built-in type and our custom uuid/set flags
// (both declared "string"); a failure there is unreachable in normal use and also
// panics. Changed and every native GetX remain available through the embedded
// *pflag.FlagSet.
type flags struct{ *pflag.FlagSet }

func (f flags) String(name string) string        { return get(f.GetString(name)) }
func (f flags) Bool(name string) bool            { return get(f.GetBool(name)) }
func (f flags) Int(name string) int              { return get(f.GetInt(name)) }
func (f flags) Int32(name string) int32          { return get(f.GetInt32(name)) }
func (f flags) StringSlice(name string) []string { return get(f.GetStringSlice(name)) }
func (f flags) IntSlice(name string) []int       { return get(f.GetIntSlice(name)) }
func (f flags) StringToString(name string) map[string]string {
	return get(f.GetStringToString(name))
}

// get collapses pflag's (val, err) pair to a single value. An undefined flag yields
// the zero value (treated as unset); any other err is a real bug and panics. See the
// flags doc for the rationale. Internal helper, not a call-site API.
func get[T any](v T, err error) T {
	if err != nil {
		var notExist *pflag.NotExistError
		if errors.As(err, &notExist) {
			return v // zero value; flag not defined -> treat as unset
		}
		panic(err)
	}
	return v
}

// Flags returns a typed, single-value view over the command's flag set.
func (c *Command) Flags() flags { return flags{c.Command.Flags()} }

// Flags returns a typed, single-value view over the command's flag set.
func (c *CommandConfig) Flags() flags { return flags{c.Command.Command.Flags()} }

// Flags returns a typed, single-value view over the command's flag set.
func (c *PreCommandConfig) Flags() flags { return flags{c.Command.Command.Flags()} }

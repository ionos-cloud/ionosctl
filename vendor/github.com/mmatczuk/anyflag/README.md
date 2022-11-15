# Anyflag

[![Build Status](https://github.com/mmatczuk/anyflag/actions/workflows/go.yml/badge.svg)](https://github.com/mmatczuk/anyflag/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mmatczuk/anyflag)](https://goreportcard.com/report/github.com/mmatczuk/anyflag)

Anyflag is an implementation of [Cobra](https://github.com/spf13/cobra) `pflag.Value` and `pflag.SliceValue` interfaces using Go Generics.

To bind your custom type to a flag, all you have to do is specify the value type and parser function, and you are done, no boilerplate.  

It supports any type including, but not limited to: enums, maps, slices, structs, struct pointers. 

## Installation

```bash
go get github.com/mmatczuk/anyflag
```

## Examples

This example shows how `anytype` can be used for JSON encoded maps.

```go
func parseJSONMap(val string) (map[string]interface{}, error) {
	var m map[string]interface{}
	return m, json.Unmarshal([]byte(val), &m)
}

func newCommand() *cobra.Command {
	var m map[string]interface{}

	cmd := &cobra.Command{
		Use: "json-map",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(cmd.OutOrStdout(), m)
		},
	}

	fs := cmd.Flags()
	value := anyflag.NewValue[map[string]interface{}](nil, &m, parseJSONMap)

	fs.VarP(value, "map", "", "map")

	return cmd
}

func main() {
	newCommand().Execute()
}
```

More examples can be found in [examples](examples) directory.

## License

This project is based on [spf13/pflag](https://github.com/spf13/pflag) licensed under the BSD 3-Clause "New" or "Revised" License

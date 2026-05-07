package engine

import (
	"fmt"
	"strings"

	"github.com/containeroo/tinyflags/internal/argparse"
	"github.com/containeroo/tinyflags/internal/core"
)

// runArgParserFSM initializes the argument parser and runs it.
// It returns any remaining positional arguments and a parsing error if any.
func runArgParserFSM(fs *FlagSet, args []string) ([]string, error) {
	return argparse.Parse(argparse.Config{
		ContinueOnError:   fs.errorHandling == ContinueOnError,
		LookupStaticFlag:  fs.lookupStaticFlag,
		LookupShortFlag:   fs.lookupShortFlag,
		LookupDynamicFlag: fs.lookupDynamicFlag,
		HandleUnknownFlag: fs.unknownFlag,
	}, args)
}

func (f *FlagSet) lookupStaticFlag(name string) *core.BaseFlag {
	return f.staticFlagsMap[name]
}

func (f *FlagSet) lookupShortFlag(short string) *core.BaseFlag {
	for _, fl := range f.staticFlagsMap {
		if fl.Short == short {
			return fl
		}
	}
	return nil
}

func (f *FlagSet) lookupDynamicFlag(name string, raw string) (core.DynamicValue, string, error) {
	parts := strings.Split(name, ".")
	if len(parts) != 3 {
		return nil, "", fmt.Errorf("invalid dynamic flag: --%s", name)
	}
	groupName, id, field := parts[0], parts[1], parts[2]

	group, ok := f.dynamicGroupsMap[groupName]
	if !ok {
		return nil, "", fmt.Errorf("unknown dynamic group %q in flag %s", groupName, raw)
	}

	item, ok := group.Items()[field]
	if !ok {
		return nil, "", fmt.Errorf("unknown dynamic field %q in flag %s", field, raw)
	}

	return item.Value, id, nil
}

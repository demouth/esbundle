package cli

import (
	"strings"

	"github.com/demouth/esbundle/pkg/api"
)

type BuildResult struct {
}

func runImpl(osArgs []string) int {

	// build

	buildOptions, err := parseOptionsForRun(osArgs)
	if err != nil {
		return 1
	}

	_ = api.Build(*buildOptions)

	return 1

	// transform
}

func parseOptionsForRun(osArgs []string) (*api.BuildOptions, error) {
	entryPoints := make([]string, 0)
	for _, arg := range osArgs {
		if !strings.HasPrefix(arg, "-") {
			options := api.BuildOptions{
				EntryPoints: entryPoints,
			}
			parseOptionsImpl(osArgs, &options)
			return &options, nil
		}
	}
	return &api.BuildOptions{
		EntryPoints: entryPoints,
	}, nil
}

func parseOptionsImpl(osArgs []string, buildOpts *api.BuildOptions) {
	for _, arg := range osArgs {
		switch {
		case strings.HasPrefix(arg, "--outfile=") && buildOpts != nil:
			buildOpts.Outfile = arg[len("--outfile="):]

		case !strings.HasPrefix(arg, "-"):
			buildOpts.EntryPoints = append(buildOpts.EntryPoints, arg)
		}
	}
}

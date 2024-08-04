package api

type BuildOptions struct {
	Outfile     string
	EntryPoints []string
}

func Build(options BuildOptions) BuildResult {
	ctx := contextImpl(options)
	rebuild := ctx.Rebuild()
	return rebuild
}

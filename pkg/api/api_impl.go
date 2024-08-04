package api

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"github.com/demouth/esbundle/internal/bundler"
	"github.com/demouth/esbundle/internal/cache"
	"github.com/demouth/esbundle/internal/config"
	"github.com/demouth/esbundle/internal/fs"
	"github.com/demouth/esbundle/internal/graph"
	"github.com/demouth/esbundle/internal/linker"
)

type internalContext struct {
	args rebuildArgs
}
type BuildResult struct {
}
type rebuildArgs struct {
	caches      *cache.CacheSet
	entryPoints []bundler.EntryPoint
	options     config.Options
}

func contextImpl(buildOpts BuildOptions) *internalContext {
	realFS, err := fs.RealFS()
	if err != nil {
		panic(err)
	}

	options, entryPoints := validateBuildOptions(buildOpts, realFS)
	args := rebuildArgs{
		caches:      cache.MakeCacheSet(),
		entryPoints: entryPoints,
		options:     options,
	}
	return &internalContext{
		args: args,
	}
}

func (c *internalContext) Rebuild() BuildResult {
	return c.rebuild()
}

func (c *internalContext) rebuild() BuildResult {
	args := c.args
	return rebuildImpl(args)
}

func rebuildImpl(args rebuildArgs) BuildResult {
	bundle := bundler.ScanBundle(args.caches, args.entryPoints, args.options)
	_ = bundle
	results := bundle.Compile(linker.Link)

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(results))
	for _, result := range results {
		go func(result graph.OutputFile) {
			if err := ioutil.WriteFile(result.AbsPath, result.Contents, 0666); err != nil {
				panic(err)
			}
			waitGroup.Done()
		}(result)
	}
	waitGroup.Wait()

	return BuildResult{}
}

func validateBuildOptions(
	buildOpts BuildOptions,
	realFS fs.FS,
) (
	options config.Options,
	entryPoints []bundler.EntryPoint,
) {
	for _, entryPoint := range buildOpts.EntryPoints {
		entryPoints = append(entryPoints, bundler.EntryPoint{
			InputPath: entryPoint,
		})
	}

	var err error
	options.AbsOutputFile, err = filepath.Abs(buildOpts.Outfile)
	if err != nil {
		panic(err)
	}
	options.AbsOutputDir = realFS.Cwd()
	return
}

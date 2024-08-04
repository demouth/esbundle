package bundler

import (
	"sync"

	"github.com/demouth/esbundle/internal/cache"
	"github.com/demouth/esbundle/internal/config"
	"github.com/demouth/esbundle/internal/graph"
	"github.com/demouth/esbundle/internal/logger"
	"github.com/demouth/esbundle/internal/resolver"
)

type scanner struct {
	caches  *cache.CacheSet
	results []parseResult

	resultChannel chan parseResult

	options config.Options

	remaining int
}

type parseResult struct {
	file scannerFile
	ok   bool
}

type scannerFile struct {
	inputFile graph.InputFile
}

func (s *scanner) addEntryPoints(entryPoints []EntryPoint) []graph.EntryPoint {

	type entryPointInfo struct {
		results []resolver.ResolveResult
	}
	entryMetas := make([]graph.EntryPoint, 0, len(entryPoints)+1)
	entryPointInfos := make([]entryPointInfo, len(entryPoints))

	entryPointWaitGroup := sync.WaitGroup{}
	entryPointWaitGroup.Add(len(entryPoints))
	for i, entryentryPoint := range entryPoints {
		go func(i int, entryPoint EntryPoint) {
			resolveResult := RunOnResolvePlugins()
			entryPointInfos[i] = entryPointInfo{
				results: []resolver.ResolveResult{*resolveResult},
			}
			entryPointWaitGroup.Done()
		}(i, entryentryPoint)
	}
	entryPointWaitGroup.Wait()

	for _, info := range entryPointInfos {
		if info.results == nil {
			continue
		}
		sourceIndex := s.maybeParseFile()

		entryMetas = append(entryMetas, graph.EntryPoint{
			SourceIndex: sourceIndex,
		})
	}
	return entryMetas
}

func (s *scanner) scanAllDependencies() {
	for s.remaining > 0 {
		result := <-s.resultChannel
		s.remaining--

		if !result.ok {
			continue
		}
		s.results = append(s.results, result)
	}
}

func (s *scanner) processScannedFiles(entryPointMeta []graph.EntryPoint) []scannerFile {
	files := make([]scannerFile, len(s.results))
	for sourceIndex, result := range s.results {
		files[sourceIndex] = result.file
	}
	return files
}

func (s *scanner) maybeParseFile() uint32 {
	s.remaining++
	sourceIndex := s.allocateSourceIndex()
	go parseFile(parseArgs{
		caches:  s.caches,
		results: s.resultChannel,
		KeyPath: "entry_point.js",
	})
	return sourceIndex
}

func (s *scanner) allocateSourceIndex() uint32 {
	return s.caches.SourceIndexCache.Get()
}

type parseArgs struct {
	results chan parseResult
	caches  *cache.CacheSet
	KeyPath string
}

func parseFile(args parseArgs) {
	source := logger.Source{
		KeyPath: args.KeyPath,
	}

	runOnLoadPlugins(&source)

	result := parseResult{}

	var loader uint8 = 0
	switch loader {
	case 0:
		ast, ok := args.caches.JSCache.Parse(source)
		result.file.inputFile.Repr = &graph.JSRepr{
			AST: ast,
		}
		result.ok = ok
	}

	args.results <- result
}

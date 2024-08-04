package bundler

import (
	"io/ioutil"

	"github.com/demouth/esbundle/internal/cache"
	"github.com/demouth/esbundle/internal/config"
	"github.com/demouth/esbundle/internal/graph"
	"github.com/demouth/esbundle/internal/js_ast"
	"github.com/demouth/esbundle/internal/js_parser"
	"github.com/demouth/esbundle/internal/linker"
	"github.com/demouth/esbundle/internal/logger"
	"github.com/demouth/esbundle/internal/resolver"
)

type Bundle struct {
	files       []scannerFile
	entryPoints []graph.EntryPoint
	options     config.Options
}

type EntryPoint struct {
	InputPath string
}

func ScanBundle(
	caches *cache.CacheSet,
	entryPoints []EntryPoint,
	options config.Options,
) Bundle {
	s := scanner{
		caches:        caches,
		options:       options,
		results:       make([]parseResult, 0),
		resultChannel: make(chan parseResult),
	}

	// TODO:<runtime>
	// s.results = append(s.results, parseResult{})
	// s.remaining++
	// go func() {
	// 	ast := globalRuntimeCache.parseRuntime()
	// 	_ = ast
	// 	s.resultChannel <- parseResult{}
	// }()

	// // TODO
	// s.results = append(s.results, parseResult{})
	// s.remaining++

	entryPointMeta := s.addEntryPoints(entryPoints)

	s.scanAllDependencies()

	files := s.processScannedFiles(entryPointMeta)

	return Bundle{
		files:       files,
		entryPoints: entryPointMeta,
		options:     options,
	}
}

func (b *Bundle) Compile(link linker.Linker) []graph.OutputFile {
	options := b.options
	files := make([]graph.InputFile, len(b.files))
	for i, file := range b.files {
		files[i] = file.inputFile
	}

	var resultGroups [][]graph.OutputFile
	resultGroups = [][]graph.OutputFile{
		link(&options, files, b.entryPoints),
	}

	var outputFiles []graph.OutputFile
	for _, group := range resultGroups {
		outputFiles = append(outputFiles, group...)
	}

	return outputFiles
}

var globalRuntimeCache = runtimeCache{}

type runtimeCache struct {
}

func (cache *runtimeCache) parseRuntime() (runtimeAST js_ast.AST) {
	// TODO
	source := logger.Source{}
	runtimeAST = js_parser.Parse(source)
	return
}

func RunOnResolvePlugins() *resolver.ResolveResult {
	// TODO
	return &resolver.ResolveResult{}
}

func runOnLoadPlugins(source *logger.Source) bool {
	buffer, err := ioutil.ReadFile(source.KeyPath)
	if err != nil {
		return false
	}
	source.Contents = string(buffer)
	return true
}

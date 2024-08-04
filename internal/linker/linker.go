package linker

import (
	"path/filepath"
	"sync"

	"github.com/demouth/esbundle/internal/config"
	"github.com/demouth/esbundle/internal/graph"
	"github.com/demouth/esbundle/internal/helpers"
	"github.com/demouth/esbundle/internal/js_printer"
)

type Linker func(
	options *config.Options,
	inputFiles []graph.InputFile,
	entryPoints []graph.EntryPoint,
) []graph.OutputFile

type linkerContext struct {
	options *config.Options
	graph   graph.LinkerGraph
	chunks  []chunkInfo
}

type intermediateOutput struct {
	joiner helpers.Joiner
}

type compileResultJS struct {
	js_printer.PrintResult
}
type partRange struct {
	sourceIndex uint32
}

type chunkInfo struct {
	chunkRepr chunkRepr

	intermediateOutput intermediateOutput

	// finalTemplate []config.PathTemplate
}

type chunkRepr interface {
	isChunk()
}

type chunkReprJS struct {
	partsInChunkInOrder []partRange
}

func (c *chunkReprJS) isChunk() {}

func Link(
	options *config.Options,
	inputFiles []graph.InputFile,
	entryPoints []graph.EntryPoint,
) []graph.OutputFile {
	c := linkerContext{
		options: options,
		graph: graph.CloneLinkerGraph(
			inputFiles,
			entryPoints,
		),
	}
	c.computeChunks()
	c.computeCrossChunkDependencies()
	return c.generateChunksInParallel()
}

func (c *linkerContext) computeChunks() {

	jsChunks := make([]chunkInfo, 0)

	for _, entryPoint := range c.graph.EntryPoints() {
		file := c.graph.Files[entryPoint.SourceIndex]
		chunk := chunkInfo{}
		switch file.InputFile.Repr.(type) {
		case *graph.JSRepr:
			chunk.chunkRepr = &chunkReprJS{}
			jsChunks = append(jsChunks, chunk)
		}
	}

	sortedChunks := make([]chunkInfo, 0)
	for _, chunk := range jsChunks {
		sortedChunks = append(sortedChunks, chunk)
	}

	for _, chunk := range sortedChunks {
		if chunkRepr, ok := chunk.chunkRepr.(*chunkReprJS); ok {
			chunkRepr.partsInChunkInOrder = c.findImportedPartsInJSOrder()
		}
	}

	// for chunkIndex := range sortedChunks {
	// 	chunk := &sortedChunks[chunkIndex]

	// 	var template []config.PathTemplate

	// 	base := filepath.Base(c.options.AbsOutputFile)
	// 	ext := filepath.Ext(c.options.AbsOutputFile)

	// 	chunk.finalTemplate = config.PathTemplate{
	// 		Dir: base + "." + ext,
	// 	}
	// }

	c.chunks = sortedChunks
}

func (c *linkerContext) computeCrossChunkDependencies() {
}

func (c *linkerContext) generateChunksInParallel() []graph.OutputFile {

	// Generate chunks

	generateWaitGroup := sync.WaitGroup{}
	generateWaitGroup.Add(len(c.chunks))
	for chunkIndex := range c.chunks {
		switch c.chunks[chunkIndex].chunkRepr.(type) {
		case *chunkReprJS:
			go c.generateChunkJS(chunkIndex, &generateWaitGroup)
		}
	}
	generateWaitGroup.Wait()

	// Generate final output files

	var resultsWaitGroup sync.WaitGroup
	results := make([][]graph.OutputFile, len(c.chunks))
	resultsWaitGroup.Add(len(c.chunks))
	for chunkIndex, chunk := range c.chunks {
		go func(chunkIndex int, chunk chunkInfo) {
			var outputFiles []graph.OutputFile

			outputContentsJoiner := c.substituteFinalPaths(chunk.intermediateOutput)
			outputContents := outputContentsJoiner.Done()

			// TODO
			base := filepath.Base(c.options.AbsOutputFile)
			base = base[:len(base)-len(filepath.Ext(base))]

			outputFiles = append(outputFiles, graph.OutputFile{
				AbsPath:  c.options.AbsOutputDir + "/" + base + ".js",
				Contents: outputContents,
			})
			results[chunkIndex] = outputFiles
			resultsWaitGroup.Done()
		}(chunkIndex, chunk)
	}
	resultsWaitGroup.Wait()

	outputFiles := make([]graph.OutputFile, 0)
	for _, result := range results {
		outputFiles = append(outputFiles, result...)
	}

	return outputFiles
}

func (c *linkerContext) generateChunkJS(chunkIndex int, generateWaitGroup *sync.WaitGroup) {
	// js_printer.Print()

	chunk := &c.chunks[chunkIndex]
	chunkRepr := chunk.chunkRepr.(*chunkReprJS)
	compileResults := make([]compileResultJS, 0)

	waitGroup := sync.WaitGroup{}
	for _, partRange := range chunkRepr.partsInChunkInOrder {

		compileResults = append(compileResults, compileResultJS{})
		compileResult := &compileResults[len(compileResults)-1]

		waitGroup.Add(1)
		go c.generateCodeForFileInChunkJS(
			&waitGroup,
			partRange,
			compileResult,
		)
	}
	waitGroup.Wait()

	j := helpers.Joiner{}
	for _, compileResult := range compileResults {
		j.AddBytes(compileResult.JS)
	}

	chunk.intermediateOutput = c.breakJoinerIntoPieces(j)

	generateWaitGroup.Done()
}

func (c *linkerContext) generateCodeForFileInChunkJS(
	waitGroup *sync.WaitGroup,
	partRange partRange,
	result *compileResultJS,
) {
	file := &c.graph.Files[partRange.sourceIndex]
	repr := file.InputFile.Repr.(*graph.JSRepr)

	printOptions := js_printer.Options{}
	tree := repr.AST

	result.PrintResult = js_printer.Print(tree, printOptions)

	// TODO

	waitGroup.Done()
}

func (c *linkerContext) findImportedPartsInJSOrder() []partRange {
	jsPartsPrefix := []partRange{}
	for i, _ := range c.graph.Files {
		jsPartsPrefix = append(jsPartsPrefix, partRange{
			sourceIndex: uint32(i),
		})
	}
	return jsPartsPrefix
}

func (c *linkerContext) substituteFinalPaths(intermediateOutput intermediateOutput) helpers.Joiner {
	return intermediateOutput.joiner
}

func (c *linkerContext) breakJoinerIntoPieces(j helpers.Joiner) intermediateOutput {
	return intermediateOutput{joiner: j}
}

package graph

type LinkerFile struct {
	InputFile InputFile
}

type EntryPoint struct {
	SourceIndex uint32
}

type LinkerGraph struct {
	Files       []LinkerFile
	entryPoints []EntryPoint
}

func (g *LinkerGraph) EntryPoints() []EntryPoint {
	return g.entryPoints
}

func CloneLinkerGraph(
	inputFiles []InputFile,
	originalEntryPoints []EntryPoint,
) LinkerGraph {
	files := make([]LinkerFile, len(inputFiles))
	for i, file := range inputFiles {
		switch repr := file.Repr.(type) {
		case *JSRepr:
			{
				clone := *repr
				repr := &clone
				files[i].InputFile.Repr = repr
			}
		}
	}

	return LinkerGraph{
		Files:       files,
		entryPoints: originalEntryPoints,
	}
}

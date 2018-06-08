package generate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/dave/jennifer/jen"
)

// type definitions used by the generate process
type GenerateFunc func(target string) (
	message string,
	err error,
	sourceFiles map[string]*jen.File)
type TargetsCollection map[string]GenerateFunc

// a header to affix to all files generated
const GENERATED_HEADER = `
//
//
// THIS FILE HAS BEEN GENERATED BY sameriver-generate
//
//
// DO NOT MODIFY BY HAND UNLESS YOU WANNA HAVE A GOOD TIME WHEN THE NEXT
// GENERATION DESTROYS WHAT YOU WROTE. UNLESS YOU KNOW HOW TO HAVE A GOOD TIME
//
//
`

// struct to hold data related to the generation
type GenerateProcess struct {
	engineDir        string
	gameDir          string
	outputDir        string
	sourceFiles      map[string]*jen.File
	messages         map[string]string
	errors           map[string]string
	rootTargets      TargetsCollection
	targetsProcessed []string
}

func NewGenerateProcess(
	engineDir string, gameDir string, outputDir string) *GenerateProcess {

	g := GenerateProcess{}
	g.engineDir = engineDir
	g.gameDir = gameDir
	g.outputDir = outputDir
	g.sourceFiles = make(map[string]*jen.File)
	g.messages = make(map[string]string)
	g.errors = make(map[string]string)
	return &g
}

// used to run the targets
func (g *GenerateProcess) Run(targets TargetsCollection) {
	for target, f := range targets {
		fmt.Printf("----- running target: %s -----\n", target)
		message, err, sourceFiles := f(target)
		g.messages[target] = message
		if err != nil {
			g.errors[target] = fmt.Sprintf("%v", err)
		}
		for filename, contents := range sourceFiles {
			g.sourceFiles[filename] = contents
		}
		g.targetsProcessed = append(g.targetsProcessed, target)
	}
}

// used to display a summary at the end
func (g *GenerateProcess) PrintReport() {
	fmt.Println("GENERATE PROCESS REPORT:\n===\n")
	for _, target := range g.targetsProcessed {
		fmt.Printf("## %s\n", target)
		msg := g.messages[target]
		fmt.Printf("message: %s\n", msg)
		err := g.errors[target]
		if err != "" {
			fmt.Printf("error: %s\n", err)
		}
		fmt.Println()
	}
}

func (g *GenerateProcess) OutputFiles() {
	fmt.Println("Generated source file output in progress...")
	defer fmt.Println("Finished output of generated source files.")

	for filename, file := range g.sourceFiles {
		// open the file to write
		outputFileName := fmt.Sprintf("%s/%s", g.outputDir, filename)
		f, err := os.Create(outputFileName)
		if err != nil {
			panic(err)
		}
		// add header and package declaration
		contents := fmt.Sprintf("%s\n%#v", GENERATED_HEADER, file)
		// write the file
		contentsBytes := []byte(contents)
		bytesWritten, err := f.Write(contentsBytes)
		if bytesWritten != len(contentsBytes) {
			panic(errors.New(fmt.Sprintf("bytes written to %s (%d) did not "+
				"match byte-count of file contents (%d) -- this is a weird "+
				"and rare error probably",
				outputFileName, bytesWritten, len(contentsBytes))))
		}
		// flush changes to disk
		f.Sync()
		fmt.Printf("%s written\n", outputFileName)
	}
}

// copy the files from ${gameDir}/sameriver/ into the outputDir
func (g *GenerateProcess) CopyFiles() {
	fmt.Printf("Copying all files from %s to %s...\n",
		g.gameDir, g.outputDir)
	files, err := ioutil.ReadDir(g.gameDir)
	if err != nil {
		panic(err)
	}
	for _, fileinfo := range files {
		filePath := path.Join(g.gameDir, fileinfo.Name())
		fmt.Printf("Copying %s...\n", filePath)
		err = exec.Command("cp", filePath, g.outputDir).Run()
		if err != nil {
			panic(err)
		}
	}
}

func (g *GenerateProcess) HadErrors() bool {
	return len(g.errors) > 0
}

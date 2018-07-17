package generate

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/dave/jennifer/jen"
)

// type definitions used by the generate process
type GeneratedFile struct {
	File    *jen.File
	Imports []string
}
type GenerateFunc func(target string) GenerateOutput
type GenerateOutput struct {
	message              string
	err                  error
	generatedSourceFiles map[string]GeneratedFile
}
type TargetsCollection map[string]GenerateFunc

// struct to hold data related to the generation
type GenerateProcess struct {
	engineDir        string
	gameDir          string
	outputDir        string
	sourceFiles      map[string]GeneratedFile
	messages         map[string]string
	errors           map[string]string
	rootTargets      TargetsCollection
	targetsProcessed []string
	debug            bool
}

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

// Init the struct
func NewGenerateProcess(
	engineDir string, gameDir string,
	outputDir string, debug bool) *GenerateProcess {

	g := GenerateProcess{}
	g.engineDir = engineDir
	g.gameDir = gameDir
	g.outputDir = outputDir
	g.sourceFiles = make(map[string]GeneratedFile)
	g.messages = make(map[string]string)
	g.errors = make(map[string]string)
	g.debug = debug
	return &g
}

// used to run the targets
func (g *GenerateProcess) Run(targets TargetsCollection) {
	for target, f := range targets {
		fmt.Printf("----- running target: %s -----\n", target)
		output := f(target)
		g.messages[target] = output.message
		if output.err != nil {
			g.errors[target] = fmt.Sprintf("%v", output.err)
		}
		for filename, generatedSourceFile := range output.generatedSourceFiles {
			g.sourceFiles[filename] = generatedSourceFile
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

	for filename, generateFile := range g.sourceFiles {
		// open the file to write
		outputFileName := fmt.Sprintf("%s/GENERATED_%s", g.outputDir, filename)
		f, err := os.Create(outputFileName)
		if err != nil {
			panic(err)
		}
		// add header and package declaration
		contents := fmt.Sprintf("%s\n", GENERATED_HEADER)
		// add the file contents
		contents += fmt.Sprintf("%#v", generateFile.File)
		// parse the generated file to find out what imports it has
		rawFile := fmt.Sprintf("%#v", generateFile.File)
		if g.debug {
			fmt.Println("==================================================")
			fmt.Printf("raw file for %s: \n\n%s\n\n", filename, rawFile)
			fmt.Println("==================================================")
		}
		importsToAdd := getImportStringsFromFileAsString(rawFile)
		nImportsAlready := len(importsToAdd)
		// add any imports to be added which are not already in the
		// generated file
		for _, importStr := range generateFile.Imports {
			importedAlready := false
			for _, importAlready := range importsToAdd {
				if importStr == importAlready {
					importedAlready = true
				}
			}
			if !importedAlready {
				importsToAdd = append(importsToAdd, importStr)
			}
		}
		// if we have imports to add
		if len(importsToAdd) > 0 {
			if nImportsAlready == 0 {
				// there is no import statement already. We can add ours
				// directly after the package statement
				packageStatementRegexp := regexp.MustCompile("\npackage .+\n")
				regexpPosition := packageStatementRegexp.
					FindStringIndex(contents)
				contents = contents[:regexpPosition[1]+1] +
					importsBlockFromStrings(importsToAdd) +
					contents[regexpPosition[1]:]
			} else {
				// replace the import statement/block of the generated file with
				// our own import block containing what was already there plus what
				// was part of the specifications of the GeneratedFile (usually the
				// imports from the custom files)
				var importStatementRegexp *regexp.Regexp
				if nImportsAlready == 1 {
					importStatementRegexp = regexp.MustCompile(
						"\nimport .+\n")
				} else {
					importStatementRegexp = regexp.MustCompile(
						"\nimport \\(\n[^\\)]+\\)\n")
				}
				regexpPosition := importStatementRegexp.
					FindStringIndex(contents)
				contents = contents[:regexpPosition[0]+1] +
					importsBlockFromStrings(importsToAdd) +
					contents[regexpPosition[1]:]
			}
		}
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
	fmt.Printf("Copying all files from %s to %s with prefix...\n",
		g.gameDir, g.outputDir)
	files, err := ioutil.ReadDir(g.gameDir)
	if err != nil {
		panic(err)
	}
	for _, fileinfo := range files {
		srcPath := path.Join(
			g.gameDir, fileinfo.Name())
		destPath := path.Join(
			g.outputDir, fmt.Sprintf("CUSTOM_%s", fileinfo.Name()))
		fmt.Printf("Copying %s...\n\tto %s\n", srcPath, destPath)
		err = exec.Command("cp", srcPath, destPath).Run()
		if err != nil {
			panic(err)
		}
	}
}

func (g *GenerateProcess) HadErrors() bool {
	return len(g.errors) > 0
}

package build

import (
	"bytes"
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/ast"
	"go/token"
	"regexp"
	"sort"
	"strings"
)

func (g *GenerateProcess) GenerateEventFiles(target string) (
	message string,
	err error,
	sourceFiles map[string]string,
	moreTargets TargetsCollection) {

	// read the events.go file as an ast.File
	srcFileName := fmt.Sprintf("%s/events.go", g.engineDir)
	eventsAst, _, err := g.ReadSourceFile(srcFileName)
	if err != nil {
		msg := fmt.Sprintf("failed to generate ast.File for %s", srcFileName)
		return msg, err, nil, nil
	}
	// traverse the declarations in the ast.File to get the event names
	eventNames := getEventNames(srcFileName, eventsAst)
	sort.Strings(eventNames)
	if len(eventNames) == 0 {
		msg := fmt.Sprintf("no structs with name matching .*Event found in %s\n",
			srcFileName)
		return msg, nil, nil, nil
	}
	// generate source files
	sourceFiles = make(map[string]string)
	// generate enum source file
	sourceFiles["events_enum.go"] = generateEventsEnumFile(eventNames)
	// return
	return "generated", nil, sourceFiles, nil
}

func getEventNames(srcFile string, astFile *ast.File) (
	eventNames []string) {
	// for each declaration in the source file
	for _, d := range astFile.Decls {
		// cast to generic declaration
		decl, ok := d.(*(ast.GenDecl))
		if !ok {
			continue
		}
		// if not a type declaration, continue
		if decl.Tok != token.TYPE {
			continue
		}
		// get the name of the type
		name := decl.Specs[0].(*ast.TypeSpec).Name.Name
		// if it's not a .+Event name, continue
		if validName, _ := regexp.MatchString(".+Event", name); !validName {
			fmt.Printf("type %s in %s does not match regexp for an event "+
				"type (\".+Event\"). Will not include in generated files.\n",
				name, srcFile)
			continue
		}
		eventNames = append(eventNames, name)
		fmt.Printf("found event: %+v\n", name)
	}
	return eventNames
}

func generateEventsEnumFile(eventNames []string) string {
	// for each event name, create an uppercase const name
	constNames := make(map[string]string)
	for _, eventName := range eventNames {
		eventNameStem := strings.Replace(eventName, "Event", "", 1)
		constNames[eventName] = strings.ToUpper(eventNameStem) + "_EVENT"
	}
	// generate the source file
	var buffer bytes.Buffer

	Type().Id("EventType").Int().Render(&buffer)
	buffer.WriteString("\n\n")

	Const().Id("N_EVENT_TYPES").Op("=").Lit(len(eventNames)).Render(&buffer)
	buffer.WriteString("\n\n")

	// write the enum
	constDefs := make([]Code, len(eventNames))
	for i, eventName := range eventNames {
		constDefs[i] = Id(constNames[eventName]).Op("=").Iota()
	}
	Const().Defs(constDefs...).Render(&buffer)
	buffer.WriteString("\n\n")

	// write the enum->string function
	Var().Id("EVENT_NAMES").Op("=").
		Map(Id("EventType")).String().
		Values(DictFunc(func(d Dict) {
			for _, eventName := range eventNames {
				constName := constNames[eventName]
				d[Id(constName)] = Lit(constName)
			}
		})).
		Render(&buffer)
	return buffer.String()
}

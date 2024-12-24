package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/r3labs/diff/v3"
)

type Item struct {
	Itemname  string `xml:"itemname" json:"ingredient_name" diff:"ingedient, identifier"`
	Itemcount string `xml:"itemcount" json:"ingredient_count" diff:"item_count"`
	Itemunit  string `xml:"itemunit" json:"ingredient_unit,omitempty" diff:"unit"`
}

type Cake struct {
	Name       string `xml:"name" json:"name" diff:"recipe, indentifier"`
	Stovetime  string `xml:"stovetime" json:"time" diff:"time"`
	Ingredient []Item `xml:"ingredients>item" json:"ingredients" diff:"ingredients, indentifier"`
}

type Recipes struct {
	Recipes []Cake `xml:"cake" json:"cake"`
}

type DBReader interface {
	Read(file []byte) (Recipes, error)
}

type Json Recipes

func (book *Json) Read(content []byte) (Recipes, error) {
	err := json.Unmarshal(content, book)
	if err != nil {
		fmt.Printf("Broken json file %v\n", err)
		os.Exit(3)
	}
	return Recipes(*book), nil
}

type Xml Recipes

func (book *Xml) Read(content []byte) (Recipes, error) {
	err := xml.Unmarshal(content, book)
	if err != nil {
		fmt.Printf("Broken xml file %v\n", err)
		os.Exit(3)
	}
	return Recipes(*book), nil
}

func fileFormat(path string) (format string) {
	if strings.HasSuffix(path, ".json") {
		format = "json"
	} else if strings.HasSuffix(path, ".xml") {
		format = "xml"
	} else {
		fmt.Println("Wrong file format")
		os.Exit(4)
		format = "0"
	}
	return
}

func outRecipes(db DBReader, content []byte) Recipes {
	outPut, err := db.Read(content)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(5)
	}
	return outPut
}

func oldJsonCase(oldfile, newfile []byte) {
	var changelog diff.Changelog
	var oldRecipes Recipes
	oldStruct := new(Json)
	oldRecipes = outRecipes(oldStruct, oldfile)
	newStruct := new(Xml)
	_ = outRecipes(newStruct, newfile)
	differ, err := diff.NewDiffer(diff.DisableStructValues())
	if err != nil {
		fmt.Println("No such file")
		os.Exit(6)
	}
	changelog, err = differ.Diff(oldStruct, newStruct)
	if err != nil {
		fmt.Println("No such file")
		os.Exit(7)
	}
	logPrinter(changelog, oldRecipes)
}

func oldXmlCase(oldfile, newfile []byte) {
	var changelog diff.Changelog
	var oldRecipes Recipes
	oldStruct := new(Xml)
	oldRecipes = outRecipes(oldStruct, oldfile)
	newStruct := new(Json)
	_ = outRecipes(newStruct, newfile)
	differ, err := diff.NewDiffer(diff.DisableStructValues())
	if err != nil {
		fmt.Println("No such file")
		os.Exit(8)
	}
	changelog, err = differ.Diff(oldStruct, newStruct)
	if err != nil {
		fmt.Println("No such file")
		os.Exit(9)
	}
	logPrinter(changelog, oldRecipes)
}

func logPrinter(changelog diff.Changelog, oldBook Recipes) {
	for _, change := range changelog {
		switch change.Type {
		case diff.CREATE:
			fmt.Printf("ADDED %s\n", change.Path)
		case diff.UPDATE:
			prettyStr := prettyPath(change.Path, oldBook)
			fmt.Printf("CHANGED %s - %s instead of %s\n", prettyStr, change.To, change.From)
		case diff.DELETE:
			prettyStr := prettyPath(change.Path, oldBook)
			fmt.Printf("REMOVED %s\n", prettyStr)
		}
	}

}

func prettyPath(path []string, book Recipes) (prettyStr string) {
	if len(path) >= 2 {
		numCake, err := strconv.Atoi(path[1])
		if err != nil {
			fmt.Printf("erorr: %v", err)
			os.Exit(10)
		}
		switch len(path) {
		case 5:
			prettyStr = path[4] + " for"
			fallthrough
		case 4:
			numIngr, err := strconv.Atoi(path[3])
			if err != nil {
				fmt.Printf("erorr: %v", err)
				os.Exit(10)
			}
			prettyStr = fmt.Sprintf("%s %s \"%s\"", prettyStr, path[2], book.Recipes[numCake].Ingredient[numIngr].Itemname)
			prettyStr = fmt.Sprintf("%s for cake \"%s\"", prettyStr, book.Recipes[numCake].Name)
		case 3:
			prettyStr = fmt.Sprintf("%s %s for cake \"%s\"", prettyStr, path[2], book.Recipes[numCake].Name)
		default:
			prettyStr = fmt.Sprintf("%s", path)
		}
	} else {
		prettyStr = fmt.Sprintf("%s", path)
	}
	return
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Unknown panic happend: ", err)
		}
	}()
	if len(os.Args) != 5 {
		fmt.Println("Use patern --old <old_file_name> --new <new_file_name>")
		os.Exit(8)
	}
	inputArgs := os.Args[1:]
	if inputArgs[0] == "--old" && inputArgs[2] == "--new" {
		oldFormat := fileFormat(inputArgs[1])
		oldfile, oldErr := os.ReadFile(inputArgs[1])
		if oldErr != nil {
			fmt.Println("No such file")
			os.Exit(11)
		}
		newFormat := fileFormat(inputArgs[3])
		newfile, newErr := os.ReadFile(inputArgs[3])
		if newErr != nil {
			fmt.Println("No such file")
			os.Exit(12)
		}
		switch oldFormat {
		case "json":
			if newFormat != "xml" {
				fmt.Println("Old and new files must have json or xml format")
				os.Exit(13)
			}
			oldJsonCase(oldfile, newfile)
		case "xml":
			if newFormat != "json" {
				fmt.Println("Old and new files must have xml or json format")
				os.Exit(13)
			}
			oldXmlCase(oldfile, newfile)
		default:
			break
		}
	} else {
		fmt.Println("Use --old --new flags.")
		os.Exit(14)
	}
}

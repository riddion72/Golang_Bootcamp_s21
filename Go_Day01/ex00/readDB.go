package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Item struct {
	Itemname  string `xml:"itemname" json:"ingredient_name"`
	Itemcount string `xml:"itemcount" json:"ingredient_count"`
	Itemunit  string `xml:"itemunit" json:"ingredient_unit,omitempty"`
}

type Cake struct {
	Name       string `xml:"name" json:"name"`
	Stovetime  string `xml:"stovetime" json:"time"`
	Ingredient []Item `xml:"ingredients>item" json:"ingredients"`
}

type Recipes struct {
	XMLName xml.Name `xml:"recipes" json:"-"`
	Recipes []Cake   `xml:"cake" json:"cake"`
}

type DBReader interface {
	Read(file []byte) (Recipes, error)
}

type Json Recipes

func (book *Json) Read(content []byte) (Recipes, error) {
	err := json.Unmarshal(content, book)
	if err != nil {
		fmt.Printf("Broken json file %v\n", err)
		os.Exit(6)
	}
	return Recipes(*book), nil
}

type Xml Recipes

func (book *Xml) Read(content []byte) (Recipes, error) {
	err := xml.Unmarshal(content, book)
	if err != nil {
		fmt.Printf("Broken xml file %v\n", err)
		os.Exit(6)
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

func outPrint(db DBReader, content []byte) []byte {
	var res []byte
	outPut, err := db.Read(content)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	switch db.(type) {
	case *Json:
		res, err = xml.MarshalIndent(outPut, "", "    ")
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return nil
		}
	case *Xml:
		res, err = json.MarshalIndent(outPut, "", "    ")
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return nil
		}
	default:
		break
	}
	return res
}

func main() {
	var res []byte
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("Unknown panic happend: ", err)
		}
	}()
	if len(os.Args) != 3 {
		fmt.Println("Use -f flag for passing path to Args.")
		os.Exit(8)
	}
	inputArgs := os.Args[1:]
	if inputArgs[0] == "-f" && len(inputArgs) > 1 {
		format := fileFormat(inputArgs[1])
		file, err := os.ReadFile(inputArgs[1])
		if err != nil {
			fmt.Println("No such file")
			os.Exit(3)
		}
		switch format {
		case "json":
			myStruct := new(Json)
			res = outPrint(myStruct, file)
		case "xml":
			myStruct := new(Xml)
			res = outPrint(myStruct, file)
		default:
			break
		}
	} else {
		fmt.Println("Use '-f' flag for passing path to Args.")
		os.Exit(5)
	}
	fmt.Printf("%s\n", res)
}

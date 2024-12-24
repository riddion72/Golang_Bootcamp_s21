package main

import (
	"fmt"
	"reflect"
)

type UnknownPlant struct {
	FlowerType string
	LeafType   string
	Color      int `color_scheme:"rgb"`
}

type AnotherUnknownPlant struct {
	FlowerColor int
	LeafType    string
	Height      int `unit:"inches"`
}

func main() {
	flower1 := UnknownPlant{FlowerType: "flores", LeafType: "septum", Color: 255}
	flower2 := AnotherUnknownPlant{FlowerColor: 10, LeafType: "lanceolate", Height: 15}
	describePlant(flower1)
	fmt.Println()
	describePlant1(flower1)
	fmt.Println()
	describePlant(flower2)
	fmt.Println()
	describePlant1(flower2)
	fmt.Println()
}

func describePlant(inPut interface{}) {
	bosi_type := reflect.TypeOf(inPut)
	bosi_value := reflect.ValueOf(inPut)
	switch inPut.(type) {
	case AnotherUnknownPlant:
		field := bosi_type.Field(2)
		fmt.Printf("FlowerColor:%v\n", bosi_value.Field(0).Interface())
		fmt.Printf("LeafType:%v\n", bosi_value.Field(1).Interface())
		if tag_value, ok := field.Tag.Lookup("unit"); ok {
			fmt.Printf("Height(unit)=%s:%v\n", tag_value, bosi_value.Field(2).Int())
		} else {
			fmt.Printf("Height:%v\n", bosi_value.Field(2).Interface())
		}
	case UnknownPlant:
		field := bosi_type.Field(2)
		fmt.Printf("FlowerType:%v\n", bosi_value.Field(0).Interface())
		fmt.Printf("LeafType:%v\n", bosi_value.Field(1).Interface())
		if tag_value, ok := field.Tag.Lookup("color_scheme"); ok {
			fmt.Printf("Color(color_scheme)=%s:%v\n", tag_value, bosi_value.Field(2).Int())
		} else {
			fmt.Printf("Color:%v\n", bosi_value.Field(2).Interface())
		}
	default:
	}
}

func describePlant1(inPut interface{}) {
	bosi_type := reflect.TypeOf(inPut)
	bosi_value := reflect.ValueOf(inPut)
	for i := 0; i < bosi_type.NumField(); i++ {
		field := bosi_type.Field(i)
		if tag_value, ok := field.Tag.Lookup("color_scheme"); ok {
			fmt.Printf("Color(color_scheme)=%s:%v\n", tag_value, bosi_value.Field(i).Int())
		} else if tag_value, ok := field.Tag.Lookup("unit"); ok {
			fmt.Printf("Height(unit)=%s:%v\n", tag_value, bosi_value.Field(i).Int())
		} else {
			fmt.Printf("%s:%v\n", field.Name, bosi_value.Field(i).Interface())
		}
	}
}

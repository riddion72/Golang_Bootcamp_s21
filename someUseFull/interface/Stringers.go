package main

import "fmt"

type IPAddr [4]byte

func (ip IPAddr) String() string {
	var s string
	for i, b := range ip {
		s = fmt.Sprintf("%v%v", s, b)
		if i != 3 {
			s += "."
		}
	}

	return s

}

func main() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	fmt.Printf("%v\n", hosts)
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}

// package main

// import "fmt"

// func main() {
// 	var i interface{} = "bkj"

// 	s := i.(string)
// 	fmt.Println(s)
// 	describe(s)
// 	a := "dsda"
// 	s, ok := i.(string)
// 	s += a
// 	describe(s)
// 	fmt.Println(s, ok)

// 	f, ok := i.(float64)
// 	fmt.Println(f, ok)
// describe(f)
// 	_, _ = i.(float64) // panic
// 	fmt.Println(ok)
// 	describe(f)
// }

// func describe(i interface{}) {
// 	//var z int
// 	//i = string(i)
// 	//fmt.Printf("(%v, %T)\n", z, z)
// 	fmt.Printf("(%v, %T)\n", i, i)
// }

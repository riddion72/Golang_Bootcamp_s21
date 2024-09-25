package main

import (
	"fmt"
	//"errors"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {

		return 0, ErrNegativeSqrt(x)
	}
	z := x / 2
	for i := 1; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}

// package main

// import (
// 	"fmt"
// 	"time"
// 	//"strconv"
// )

// type MyError struct {
// 	When time.Time
// 	What string
// }

// func (e *MyError) Error() string {
// 	return fmt.Sprintf("at %v, %s",
// 		e.When, e.What)
// }

// func run() error {
// 	fmt.Printf("(, %T)\n", &MyError{
// 		time.Now(),
// 		"it didn't work",
// 	})
// 	return nil
// }

// func main() {
// 	if err := run(); err != nil {
// 		fmt.Println(err)
// 	}
// }

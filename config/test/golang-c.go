package main

/*
int Sqrt(int a) {
	return a*a;
}
*/
import "C"

import "fmt"

func main() {
	void, err := C.Sqrt(3) 
	fmt.Printf("%#v, %#v\n", void, err)
}
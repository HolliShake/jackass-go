package main
import "fmt"

func main() {
	jackAss := JackAss()
	jackAss.execute("./tests/test.j4")
	fmt.Println("Done!")
}
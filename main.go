package main
import (
	"fmt"
	"os"
)


func main() {
	// Arguments
	args := os.Args[1:]
	fmt.Println(args)

	// Create a new lexer
	lexer := Lexer("main.js", "+=")
	fmt.Println(lexer.nextToken())
}
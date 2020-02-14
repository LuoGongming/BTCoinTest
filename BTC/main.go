package main

//import "fmt"

func main() {
	bc := NewBlockChain("14PxkwD8cTpzNAT1PYXRwK4qRNbkBVtgFP")
	cli := CLI{bc}
	cli.Run()
}

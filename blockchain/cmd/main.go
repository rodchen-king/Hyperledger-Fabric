package main

import "../core"

func main() {
	bc := core.CreateNewBlockChain()
	bc.SetData(".")
	bc.SetData("This is the second block.")
	bc.Print()
}

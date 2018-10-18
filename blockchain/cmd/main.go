package main

import "../core"

func main() {
	bc := core.CreateNewBlockChain()
	bc.SetData("A借B 100元，2018/09/01")
	bc.SetData("A借C 200元，2018/09/05")
	bc.SetData("A还B 50元， 2018/09/10")
	bc.SetData("C借B 150元，2018/09/15")
	bc.SetData("C还B 100元，2018/09/20")
	bc.Print()
}

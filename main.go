package main

import "islands/generator"

func main() {
	island := generator.NewIsland(5)

	island.GenIsland()
	island.SaveImage("myIsland125r")
}

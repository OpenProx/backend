package main

import "github.com/OpenProx/backend"

func main() {
	i := backend.Instance{}
	panic(i.InitAndRun(":8080"))
}

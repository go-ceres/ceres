package main

import (
	{{.ImportStr}}
)

func main() {
    if err := bootstrap.Run(); err != nil {
        panic(err)
    }
}

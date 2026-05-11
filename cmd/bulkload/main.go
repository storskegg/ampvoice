package main

import (
    "fmt"

    "github.com/storskegg/ampvoice/app/bulkload"
)

func main() {
    if err := bulkload.Run(); err != nil {
        fmt.Println(err)
    }
}

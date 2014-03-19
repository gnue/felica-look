package main

import (
	"./felica"
	"fmt"
	"os"
)

func main() {
	for _, v := range os.Args[1:] {
		cardinfo := felica.Read(v)

		for syscode, currsys := range *cardinfo {
			fmt.Println("SYSTEM CODE: ", syscode)
			fmt.Println("  IDm: ", currsys.IDm())
			fmt.Println("  PMm: ", currsys.PMm())

			for svccode, data := range currsys.Services() {
				fmt.Println("  SERVICE CODE: ", svccode)

				for _, v := range data {
					fmt.Printf("      %X\n", v)
				}
			}
		}
	}
}

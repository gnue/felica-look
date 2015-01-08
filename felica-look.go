package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gnue/felica-look/felica"
)

// コマンドの使い方
func usage() {
	cmd := os.Args[0]
	fmt.Fprintf(os.Stderr, "usage: %s [options] [file...]\n", filepath.Base(cmd))
	flag.PrintDefaults()
	os.Exit(0)
}

// コード・リストを文字列のリストにする
func codes_to_strings(codes []uint16) []string {
	var list []string

	for _, code := range codes {
		list = append(list, fmt.Sprintf("%04X", code))
	}

	return list
}

// カード情報を簡易出力する
func show_info(cardinfo felica.CardInfo) {
	for syscode, currsys := range cardinfo {
		fmt.Printf("SYSTEM CODE: %04X\n", syscode)
		fmt.Println("  IDm: ", currsys.IDm)
		fmt.Println("  PMm: ", currsys.PMm)
		fmt.Println("  SERVICE CODES: ", codes_to_strings(currsys.ServiceCodes()))
	}
}

// カード情報をダンプ出力する
func dump_info(cardinfo felica.CardInfo) {
	for syscode, currsys := range cardinfo {
		fmt.Printf("SYSTEM CODE: %04X\n", syscode)
		fmt.Println("  IDm: ", currsys.IDm)
		fmt.Println("  PMm: ", currsys.PMm)

		for svccode, data := range currsys.Services {
			fmt.Printf("  SERVICE CODE: %04X\n", svccode)

			for _, v := range data {
				fmt.Printf("      %X\n", v)
			}
		}
	}
}

// カード情報をJSON出力
func output_json(cardinfo felica.CardInfo) {
	m := make(map[string]interface{})

	for syscode, sysinfo := range cardinfo {
		value := make(map[string]interface{})
		services := make(map[string]([]string))

		value["IDm"] = sysinfo.IDm
		value["PMm"] = sysinfo.PMm
		value["Services"] = services

		for svccode, data := range sysinfo.Services {
			key := fmt.Sprintf("%04X", svccode)
			raws := make([]string, len(data))
			services[key] = raws

			for i, v := range data {
				raws[i] = fmt.Sprintf("%X", v)
			}
		}

		key := fmt.Sprintf("%04X", syscode)
		m[key] = value
	}

	buf, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println(string(buf))
}

func main() {
	opts := felica.Options{}

	flag.BoolVar(&opts.Extend, "e", false, "extend information")
	flag.BoolVar(&opts.Hex, "x", false, "with hex dump")
	dump := flag.Bool("d", false, "dump")
	help := flag.Bool("h", false, "help")
	json := flag.Bool("json", false, "output JSON format")
	flag.Parse()

	if *help {
		usage()
	}

	if *json {
		opts.Format = felica.OUTPUT_JSON
	}

	show := func(path string) {
		cardinfo := felica.Read(path)

		if *dump {
			if *json {
				output_json(cardinfo)
			} else {
				dump_info(cardinfo)
			}
		} else {
			m := felica.Find(cardinfo)
			if m != nil {
				engine := m.Bind(cardinfo)

				fmt.Printf("%s:\n", engine.Name())
				engine.ShowInfo(&opts)
			} else {
				show_info(cardinfo)
			}
		}
	}

	if len(flag.Args()) == 0 {
		show("")
	} else {
		for _, v := range flag.Args() {
			show(v)
		}
	}
}

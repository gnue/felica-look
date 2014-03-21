package felica

import (
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
)

// YAML を読込む
func load_yaml(path string) (map[interface{}]interface{}, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	m := make(map[interface{}]interface{})
	err = goyaml.Unmarshal(contents, &m)

	return m, err
}

// テーブルを検索して表示用の文字列を返す
func disp_name(tables map[interface{}]interface{}, name string, value int, base int, opt_values ...int) interface{} {
	var v interface{}

	t := tables[name]

	if t != nil {
		v = t.(map[interface{}]interface{})[value]
	}

	if v == nil {
		if 0 < len(opt_values) {
			value = opt_values[0]
		}

		if base != 0 {
			f := fmt.Sprintf("%s%dX", "0x%0", base)
			v = fmt.Sprintf(f, value)
		} else {
			v = value
		}
	}

	return v
}

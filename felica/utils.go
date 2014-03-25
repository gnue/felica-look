package felica

import (
	"errors"
	"fmt"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
)

// ファイルを検索ディレクトリから探す
func search_file(fname string, dirs []string) (string, error) {
	for _, dir := range dirs {
		path := filepath.Join(dir, fname)

		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", errors.New("file not found")
}

// YAML を読込む
func LoadYAML(fname string) (map[interface{}]interface{}, error) {

	dirs := []string{".", filepath.Dir(os.Args[0])}
	path, err := search_file(fname, dirs)

	if err != nil {
		// ファイルが見つからなかった
		return nil, err
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		// 読込みに失敗
		return nil, err
	}

	m := make(map[interface{}]interface{})
	err = goyaml.Unmarshal(contents, &m)

	return m, err
}

// テーブルを検索して表示用の文字列を返す
func DispName(tables map[interface{}]interface{}, name string, value int, base int, opt_values ...int) interface{} {
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

package felica

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"

	"github.com/gfx/go-visual_width"
	"launchpad.net/goyaml"
)

// C言語で使うためにデータにアクセスするポインタを取得する
func DataPtr(data *[]byte) unsafe.Pointer {
	raw := (*reflect.SliceHeader)(unsafe.Pointer(data)).Data

	return unsafe.Pointer(raw)
}

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
	cmd, _ := filepath.EvalSymlinks(os.Args[0])
	bindir := filepath.Dir(cmd)
	moddir := filepath.Join(bindir, strings.TrimSuffix(fname, ".yml"))

	dirs := []string{".", bindir, moddir}
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

// 指定された表示文字になるように調整する
func DispString(str string, width int) string {
	s := visual_width.Truncate(str, true, width, "")
	n := width - visual_width.Measure(s, true)

	if 0 < n {
		s = s + strings.Repeat(" ", n)
	}

	return s
}

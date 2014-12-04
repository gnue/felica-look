package felica

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode/utf8"
	"unsafe"

	"github.com/moznion/go-unicode-east-asian-width"
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
	s, rest := StringTruncate(str, width, "…")

	if 0 < rest {
		s = s + strings.Repeat(" ", rest)
	}

	return s
}

// 文字列の表示幅
func StringWidth(str string) (w int) {
	for _, r := range str {
		if eastasianwidth.IsFullwidth(r) {
			w += 2
		} else {
			w++
		}
	}

	return
}

// 文字列の表示幅で切り詰める
func StringTruncate(str string, width int, tail string) (s string, rest int) {
	tailLen := StringWidth(tail)
	tailRest := width
	tp := 0
	i := 0

	rest = width
	s = str

	for _, r := range str {
		n := 1
		if eastasianwidth.IsFullwidth(r) {
			n = 2
		}

		if tp == 0 && rest < n+tailLen {
			// tail を追加できる位置を覚えておく
			tp = i
			tailRest = rest
		}

		if rest < n {
			s = str[0:tp] + tail
			rest = tailRest - tailLen
			break
		}

		rest -= n
		i += utf8.RuneLen(r)
	}

	return
}

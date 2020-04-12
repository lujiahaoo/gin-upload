package uploader

import (
	"reflect"
	"testing"
)

var (
	TestFileName = []string{".jpg", ".png", ".gif", ".sh", "ddd"}
)

func TestIsAllowImage(t *testing.T) {
	res := []bool{}
	for _, v := range TestFileName {
		res = append(res, IsAllowImage(v))
	}

	if !reflect.DeepEqual(res, []bool{true, true, true, false, false}) {
		t.Fatal("出错")
	}
	t.Log("成功")
}

func TestGenRandomString(t *testing.T) {
	length := 10

	res := GenRandomString(length)

	if len(res) != length {
		t.Fatal("出错")
	}
	t.Log("成功")
}
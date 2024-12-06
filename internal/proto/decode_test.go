package proto

import (
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"os"
	"strings"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	wd, _ := os.Getwd()
	var apiFilePath = strings.SplitN(wd, "/internal", 2)[0] + "/example/api/service.api"
	api, err := parser.Parse(apiFilePath, "")
	if err != nil {
		t.Fatal(err)
	}
	f, err := Unmarshal(api)
	if err != nil {
		t.Fatal(err)
	}
	data, err := f.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
}

package config

import (
	"fmt"
	"strings"
	"testing"
)

const (
	yamlStr = `testing:
  ps_config:
    ja-JP:
      hmm: "aaa"
      dict: "bbb"
      lm: "lm"
      kws_threshold: "1e-20"
      debug: "2"
`
)

func TestUnmarshallYaml(t *testing.T) {
	r := strings.NewReader(yamlStr)
	cfg, err := Load(r, "testing")
	fmt.Println(cfg)
	if err != nil {
		t.Error(err)
	}
}

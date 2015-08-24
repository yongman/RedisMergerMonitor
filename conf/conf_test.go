package conf

import (
	"fmt"
	"testing"
)

func TestLoadConf(t *testing.T) {
	conf, err := LoadConf("../redis-monitor.yml")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(conf)
}

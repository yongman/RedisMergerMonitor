package conf

import (
	"gopkg.in/yaml.v1"
	"io/ioutil"
)

type MonitorConf struct {
	HttpListen      string   `yaml:"httplisten,omitempty"`
	WsListen        string   `yaml:"wslisten,omitempty"`
	ServerType      string   `yaml:"servertype,omitempty"`
	MonitorInterval string   `yaml:"monitorinterval,omitempty"`
	Servers         []string `yaml:"servers,omitempty"`
}

func LoadConf(filename string) (*MonitorConf, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	mc := &MonitorConf{}
	err = yaml.Unmarshal(content, mc)
	if err != nil {
		return nil, err
	}
	return mc, nil
}

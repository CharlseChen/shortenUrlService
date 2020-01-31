package config

import (
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"golang.org/x/net/context"
)

type Configure struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var conf = Configure{
	Host: "127.0.0.1",
	Port:"50053",
}

func NewConfigure() (*Configure, error) {
	loader := confita.NewLoader(file.NewBackend("../config/tsconfig.json"))
	err:=loader.Load(context.Background(), &conf)
	//读取
	if err==nil {
		return &conf, nil
	}
	return nil,err
}

package main

import (
	"projectlayout/app/cmd/album/internal/conf"
	"projectlayout/app/cmd/album/internal/data"

	"github.com/google/wire"
)

func InitConf() error {
	err := conf.Init(flagconf)
	return err
}

func InitDB(err error) interface{} {
	if err != nil {
		panic(err)
	}
	if err = data.InitDB(); err != nil {
		panic(err)
	}
	return nil
}

func InitSetting() interface{} {
	wire.Build(InitDB, InitConf)
	return nil
}

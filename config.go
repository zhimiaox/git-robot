package main

import (
	"bytes"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/configor"
)

type cfg struct {
	App struct {
		PageSize         int
		JwtSecret        string
		MaxAuditProgress int
	}
	Mysql struct {
		Host        string
		User        string
		Password    string
		Database    string
		TablePrefix string
	}
	Redis struct {
		Host        string
		Auth        string
		MaxIdle     int
		MaxActive   int
		IdleTimeOut int
	}
}

var Config = &cfg{}
var filePath = "config.toml"

// Init 初始化配置
func (c *cfg) Init() error {
	return configor.Load(Config, filePath)
}

// ENV 获取当前配置场景
func (c *cfg) ENV() string {
	return configor.ENV()
}

// Save 保存配置
func (c *cfg) Save() (err error) {
	var (
		file   *os.File
		buffer bytes.Buffer
	)
	if file, err = os.Create(filePath); err != nil {
		return
	}
	defer file.Close()
	err = toml.NewEncoder(&buffer).Encode(Config)
	if err != nil {
		return
	}
	_, err = file.Write(buffer.Bytes())
	return
}

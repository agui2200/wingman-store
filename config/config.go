package config

import "github.com/jinzhu/configor"

type config struct {
	SchemaPackage string `json:"schema_path" required:"true"`            // ent.schema的目录
	TargetPackage string `json:"target" required:"true" default:"store"` // 生成代码的目标package
	Feature       struct {
		Privacy      bool // ent 特性 隐私层，数据权限控制
		EntQL        bool // ent 特性 联结 GraphQL
		Snapshot     bool // ent 特性 结构快照
		SchemaConfig bool // ent 特性 多数据库多表
	}
}

var C = config{}

func LoadConfig(f string) error {
	return configor.New(&configor.Config{
		Silent: true,
		Debug:  false,
	}).Load(&C, f)
}

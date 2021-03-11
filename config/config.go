package config

import "github.com/jinzhu/configor"

type config struct {
	SchemaPackage       string `json:"schema_path" required:"true"`            // ent.schema的目录
	TargetPackage       string `json:"target" required:"true" default:"store"` // 生成代码的目标package
	FeaturePrivacy      bool   `json:"feature_privacy"`                        // ent 特性 隐私层，数据权限控制
	FeatureEntQL        bool   `json:"feature_ent_ql"`                         // ent 特性 联结 GraphQL
	FeatureSnapshot     bool   `json:"feature_snapshot"`                       // ent 特性 结构快照
	FeatureSchemaConfig bool   `json:"feature_schema_config"`                  // ent 特性 多数据库多表
}

var C = config{}

func LoadConfig(f string) error {
	return configor.New(&configor.Config{
		Silent: false,
		Debug:  false,
	}).Load(&C, f)
}

package config

import (
	"go.uber.org/zap"
)

type Server struct {

	// 跨域配置
	Cors CORS `json:"cors" yaml:"cors"`

	//本地化
	Language Language `json:"language" yaml:"language"`

	HashIds HashIds `json:"hashIds" yaml:"hashIds"`
	//默认日志程序
	ServerLog ServerLog `json:"server_log" yaml:"server_log"`

	JWT      JWT      `json:"jwt" yaml:"jwt"`
	Redis    Redis    `json:"redis" yaml:"redis"`
	Clusters Clusters `json:"clusters" yaml:"clusters"`

	System   System   `json:"system" yaml:"system"`
	Captcha  Captcha  `mjson:"captcha" yaml:"captcha"`
	AutoCode Autocode `json:"autocode" yaml:"autocode"`
	// Mysqlsqlx
	MysqlDsn string  `json:"mysql_dsn"  yaml:"mysql_dsn"`
	Mysql    MysqlDB `json:"mysql" yaml:"mysql"`

	DBList []SpecializedDB `json:"db-list" yaml:"db-list"`
	// 存储
	Local Local `json:"local" yaml:"local"`
	AwsS3 AwsS3 `json:"aws-s3" yaml:"aws-s3"`

	Excel Excel `json:"excel" yaml:"excel"`
	Timer Timer `json:"timer" yaml:"timer"`

	//定时任务 解析配置
	Anysis Anysis `json:"anysis" yaml:"anysis"`
}

type ServerLog struct {
	*zap.SugaredLogger
}

func NewServerLog(log *zap.SugaredLogger) *ServerLog {
	return &ServerLog{log}
}

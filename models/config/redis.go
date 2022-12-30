package config

type Redis struct {
	Addr         string `mapstructure:"addr" json:"addr" yaml:"addr"`             // 服务器地址:端口
	Password     string `mapstructure:"password" json:"password" yaml:"password"` // 密码
	DB           int    `mapstructure:"db" json:"db" yaml:"db"`                   // redis的哪个数据库
	MaxRetries   int    `mapstructure:"max_retries" json:"max_retries" yaml:"max_retries"`
	PoolSize     int    `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns" json:"min_idle_conns" yaml:"min_idle_conns"`
}

type Clusters struct {
	Addrs    string `mapstructure:"addrs" json:"addrs" yaml:"addrs"`          // 服务器主节点地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                   // redis的哪个数据库

}

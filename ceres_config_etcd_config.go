package CeresConfigEtcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type Config struct {
	*clientv3.Config
	Prefix     string // etcd配置路径
	TrimPrefix string // 删除掉的头部字符串
	Encoding   string // 加解密
	Ctx        context.Context
}

// 默认的配置信息
func defaultConfig() *Config {
	conf := &Config{
		Config: &clientv3.Config{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5 * time.Second,
		},
		Ctx:      context.Background(),
		Prefix:   DefaultPrefix,
		Encoding: "json",
	}
	return conf
}

// 初始化参数
type Option func(o *Config)

// Addr 连接地址
func Addr(addrs ...string) Option {
	return func(o *Config) {
		o.Endpoints = addrs
	}
}

// Prefix 前缀
func Prefix(prefix string) Option {
	return func(o *Config) {
		o.Prefix = prefix
	}
}

// StripPrefix 前缀
func StripPrefix(stripPrefix string) Option {
	return func(o *Config) {
		o.TrimPrefix = stripPrefix
	}
}

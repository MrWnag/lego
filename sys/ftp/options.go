package ftp

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	IP       string //ftp 地址
	Port     int32  //ftp 访问端口
	User     string //ftp 登录用户名
	Password string //ftp 登录密码
}

func SetIP(v string) Option {
	return func(o *Options) {
		o.IP = v
	}
}
func SetPort(v int32) Option {
	return func(o *Options) {
		o.Port = v
	}
}
func SetUser(v string) Option {
	return func(o *Options) {
		o.User = v
	}
}
func SetPassword(v string) Option {
	return func(o *Options) {
		o.Password = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		IP:   "127.0.0.1",
		Port: 21,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		IP:   "127.0.0.1",
		Port: 21,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

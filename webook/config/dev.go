// !go:build k8s
// 没有k8s这个编译标签
package config

var Config = config{
	DB: DBConfig{
		// 本地连接
		//DSN: "localhost:13316",
		DSN: "root:root@tcp(localhost:13316)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:6379",
	},
}
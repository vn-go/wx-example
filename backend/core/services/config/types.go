package config

type configInfo struct {
	AppName  string `mapstructure:"app"`
	Database struct {
		Driver string `mapstructure:"driver"`
		DSN    string `mapstructure:"dsn"`
	} `mapstructure:"database"`

	Cache struct {
		PrefixKey string `mapstructure:"prefix"`
		Type      string `mapstructure:"type"`
		Memcache  struct {
			Servers string `mapstructure:"servers"`
		} `mapstructure:"memcache"`
		Redis struct {
			Address  string `mapstructure:"address"`
			Password string `mapstructure:"password"`
			Db       int    `mapstructure:"db"`
		} `mapstructure:"redis"`
		Badger struct {
			Directory string `mapstructure:"directory"`
		} `mapstructure:"badger"`
	} `mapstructure:"cache"`
	Jwt struct {
		Sercret          string `mapstructure:"default-sercret-key"`
		HashPasswordType string `mapstructure:"hash-type"`
		SecretLen        int    `mapstructure:"share-secret-len"`
	} `mapstructure:"jwt"`
	Log struct {
		Path     string `mapstructure:"path"`
		Size     int    `mapstructure:"size"`
		Age      int    `mapstructure:"age"`
		Backup   int    `mapstructure:"backup"`
		Compress bool   `mapstructure:"compress"`
	} `mapstructure:"log"`
	Cryptor struct {
		Key string `mapstructure:"key"`
	} `mapstructure:"cryptor"`
	Broker struct {
		Bus    string `mapstructure:"bus"`
		Topic  string `mapstructure:"topic"`
		Rabbit struct {
			Url      string `mapstructure:"url"`
			Exchange string `mapstructure:"exchange"`
			Queue    string `mapstructure:"queue"`
		} `mapstructure:"rabbit"`
		Redis struct {
			Addr     string `mapstructure:"addr"`
			Consumer string `mapstructure:"consumer"`
		} `mapstructure:"redis"`
		Kafka struct {
			Brokers string `mapstructure:"brokers"`
		} `mapstructure:"kafka"`
		Nats struct {
			Server string `mapstructure:"server"`
			Group  string `mapstructure:"group"`
		} `mapstructure:"nats"`
	}
	Tenant struct {
		IsMulti              bool   `mapstructure:"IsMulti"`
		DefaultAdminUser     string `mapstructure:"DefaultAdminUser"`
		DefaultAdminPassword string `mapstructure:"DefaultAdminPassword"`
	} `mapstructure:"tenant"`
	DefaultAuth struct {
		Username string `mapstructure:"user"`
		Password string `mapstructure:"pwd"`
	} `mapstructure:"default-auth"`
	Debug bool `mapstructure:"debug"`
	Bind  struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	} `mapstructure:"bind"`
}

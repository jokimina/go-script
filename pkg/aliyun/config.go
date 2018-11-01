package aliyun

import "os"

type BaseConfig interface {
	GetConfigFromEnv() *Config
}

type Config struct {
	AccessKeyId     string
	AccessKeySecret string
	RegionId        string
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) GetConfigFromEnv() *Config {
	c = &Config{
		AccessKeyId:     os.Getenv("ACCESS_KEY_ID"),
		AccessKeySecret: os.Getenv("ACCESS_KEY_SECRET"),
		RegionId:        os.Getenv("REGION_ID"),
	}
	if c.AccessKeyId == "" || c.AccessKeySecret == "" || c.RegionId == "" {
		panic("Get REGION_ID ACCESS_KEY_ID ACCESS_KEY_SECRET from environment variables failed")
	} else {
		return c
	}
}

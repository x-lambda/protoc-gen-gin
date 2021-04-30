package conf

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// Hostname 主机名
	Hostname = "localhost"
	// AppID 获取 APP_ID
	AppID = "localapp"
	// IsDevEnv 开发环境标志
	IsDevEnv = false
	// IsUatEnv 集成环境标志
	IsUatEnv = false
	// IsProdEnv 生产环境标志
	IsProdEnv = false
	// Env 运行环境
	Env = "dev"
	// Zone 服务区域
	Zone = "sh001"
)

var path string
var v *viper.Viper

func init() {
	Hostname, _ = os.Hostname()

	if appID := os.Getenv("APP_ID"); appID != "" {
		AppID = appID
	} else {
		logger().Warn("env APP_ID is empty")
	}

	if env := os.Getenv("DEPLOY_ENV"); env != "" {
		Env = env
	} else {
		logger().Warn("env DEPLOY_ENV is empty")
	}

	switch Env {
	case "prod", "pre":
		IsProdEnv = true
	case "uat":
		IsUatEnv = true
	default:
		IsDevEnv = true
	}

	path = os.Getenv("CONF_PATH")
	if path == "" {
		logger().Warn("env CONF_PATH is empty")

		var err error
		if path, err = os.Getwd(); err != nil {
			panic(err)
		}

		logger().WithField("path", path).Info("use default conf path")
	}

	v = viper.New()
	viper.SetConfigName("example")
	viper.AddConfigPath(path)
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.AutomaticEnv()
}

// OnConfigChange 注册文件变更回调，需要在WatchConfig()之前调用
// Warning: 业务代码不要调用
func OnConfigChange(run func()) {
	v.OnConfigChange(func(in fsnotify.Event) {
		run()
	})
}

// WatchConfig 启动配置变更监听
// Warning: 业务代码不要调用
func WatchConfig() {
	v.WatchConfig()
}

var levels = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
}

func logger() *logrus.Entry {
	if level, ok := levels[os.Getenv("LOG_LEVEL")]; ok {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logrus.DebugLevel)
	}

	return logrus.WithFields(logrus.Fields{
		"app_id":      AppID,
		"instance_id": Hostname,
		"env":         Env,
	})
}

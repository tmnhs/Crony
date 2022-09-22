package models

type (
	Mysql struct {
		Path         string `mapstructure:"path" json:"path" yaml:"path"`                             // 服务器地址
		Port         string `mapstructure:"port" json:"port" yaml:"port"`                             // 端口
		Config       string `mapstructure:"config" json:"config" yaml:"config"`                       // 高级配置
		Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`                     // 数据库名
		Username     string `mapstructure:"username" json:"username" yaml:"username"`                 // 数据库用户名
		Password     string `mapstructure:"password" json:"password" yaml:"password"`                 // 数据库密码
		MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"` // 空闲中的最大连接数
		MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
		LogMode      string `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`                  // 是否开启Gorm全局日志
		LogZap       bool   `mapstructure:"log-zap" json:"logZap" yaml:"log-zap"`                     // 是否通过zap写入日志文件
	}
	JWT struct {
		SigningKey  string `mapstructure:"signing-key" json:"signingKey" yaml:"signing-key"`    // jwt签名
		ExpiresTime int64  `mapstructure:"expires-time" json:"expiresTime" yaml:"expires-time"` // 过期时间
		BufferTime  int64  `mapstructure:"buffer-time" json:"bufferTime" yaml:"buffer-time"`    // 缓冲时间
		Issuer      string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`                  // 签发者
	}
	Email struct {
		Port     int    `mapstructure:"port" json:"port" yaml:"port"`             // 端口
		From     string `mapstructure:"from" json:"from" yaml:"from"`             // 收件人
		Host     string `mapstructure:"host" json:"host" yaml:"host"`             // 服务器地址
		IsSSL    bool   `mapstructure:"is-ssl" json:"isSSL" yaml:"is-ssl"`        // 是否SSL
		Secret   string `mapstructure:"secret" json:"secret" yaml:"secret"`       // 密钥
		Nickname string `mapstructure:"nickname" json:"nickname" yaml:"nickname"` // 昵称
	}
	System struct {
		Env          string `mapstructure:"env" json:"env" yaml:"env"`    // 环境值
		Addr         int    `mapstructure:"addr" json:"addr" yaml:"addr"` // 端口值
		LimitCountIP int    `mapstructure:"iplimit-count" json:"iplimitCount" yaml:"iplimit-count"`
		LimitTimeIP  int    `mapstructure:"iplimit-time" json:"iplimitTime" yaml:"iplimit-time"`
	}
	Zap struct {
		Level         string `mapstructure:"level" json:"level" yaml:"level"`                           // 级别
		Format        string `mapstructure:"format" json:"format" yaml:"format"`                        // 输出
		Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                        // 日志前缀
		Director      string `mapstructure:"director" json:"director"  yaml:"director"`                 // 日志文件夹
		ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"showLine"`                 // 显示行
		EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`       // 编码级
		StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"` // 栈名
		LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`  // 输出控制台
	}
	Config struct {
		JWT    JWT    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
		Zap    Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
		Email  Email  `mapstructure:"email" json:"email" yaml:"email"`
		System System `mapstructure:"system" json:"system" yaml:"system"`
		Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	}
)

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

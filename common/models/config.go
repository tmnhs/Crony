package models

const (
	CronyNodeTableName   = "node"
	CronyJobTableName    = "job"
	CronyJobLogTableName = "job_log"
	CronyUserTableName   = "user"
)

type (
	Mysql struct {
		Path         string `mapstructure:"path" json:"path" yaml:"path" ini:"path"`                                       // 服务器地址
		Port         string `mapstructure:"port" json:"port" yaml:"port" ini:"port"`                                       // 端口
		Config       string `mapstructure:"config" json:"config" yaml:"config" ini:"config"`                               // 高级配置
		Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name" ini:"db-name"`                            // 数据库名
		Username     string `mapstructure:"username" json:"username" yaml:"username" ini:"username"`                       // 数据库用户名
		Password     string `mapstructure:"password" json:"password" yaml:"password" ini:"password"`                       // 数据库密码
		MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns" ini:"max-idle-conns"` // 空闲中的最大连接数
		MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns" ini:"max-open-conns"` // 打开到数据库的最大连接数
		LogMode      string `mapstructure:"log-mode" json:"logMode" yaml:"log-mode" ini:"log-mode"`                        // 是否开启Gorm全局日志
		LogZap       bool   `mapstructure:"log-zap" json:"logZap" yaml:"log-zap" ini:"log-zap"`                            // 是否通过zap写入日志文件
	}
	Email struct {
		Port     int      `mapstructure:"port" json:"port" yaml:"port" ini:"port"`                 // 端口
		From     string   `mapstructure:"from" json:"from" yaml:"from" ini:"from"`                 // 收件人
		Host     string   `mapstructure:"host" json:"host" yaml:"host" ini:"host"`                 // 服务器地址
		IsSSL    bool     `mapstructure:"is-ssl" json:"isSSL" yaml:"is-ssl" ini:"is-ssl"`          // 是否SSL
		Secret   string   `mapstructure:"secret" json:"secret" yaml:"secret" ini:"secret"`         // 密钥
		Nickname string   `mapstructure:"nickname" json:"nickname" yaml:"nickname" ini:"nickname"` // 昵称
		To       []string `mapstructure:"to" json:"to" yaml:"to" ini:"to"`                         // 默认邮件发送对象
	}
	WebHook struct {
		Kind string `mapstructure:"kind" json:"kind" yaml:"kind" ini:"kind"`
		Url  string `mapstructure:"url" json:"url" yaml:"url" ini:"kind"`
	}
	Etcd struct {
		Endpoints   []string `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints" ini:"endpoints"`
		Username    string   `mapstructure:"username" json:"username" yaml:"username" ini:"username"` // 库用户名
		Password    string   `mapstructure:"password" json:"password" yaml:"password" ini:"password"` // 密码
		DialTimeout int64    `mapstructure:"dial-timeout" json:"dial-timeout" yaml:"dial-timeout" ini:"dial-timeout"`
		ReqTimeout  int64    `mapstructure:"req-timeout" json:"req-timeout" yaml:"req-timeout" ini:"req-timeout"`
	}
	System struct {
		Env                string `mapstructure:"env" json:"env" yaml:"env" ini:"env"`                                     // 环境值
		Addr               int    `mapstructure:"addr" json:"addr" yaml:"addr" ini:"addr"`                                 // 端口值
		NodeTtl            int64  `mapstructure:"node-ttl" json:"node-ttl" yaml:"node-ttl" ini:"node-ttl"`                 //
		JobProcTtl         int64  `mapstructure:"job-proc-ttl" json:"job-proc-ttl" yaml:"job-proc-ttl" ini:"job-proc-ttl"` //
		Version            string `mapstructure:"version" json:"version" yaml:"version" ini:"version"`                     //
		LogCleanPeriod     int64  `mapstructure:"log-clean-period" json:"log-clean-period" yaml:"log-clean-period" ini:"log-clean-period"`
		LogCleanExpiration int64  `mapstructure:"log-clean-expiration" json:"log-clean-expiration" yaml:"log-clean-expiration" ini:"log-clean-expiration"`
		CmdAutoAllocation  bool   `mapstructure:"cmd-auto-allocation" json:"cmd-auto-allocation" yaml:"cmd-auto-allocation" ini:"cmd-auto-allocation"`
	}
	Log struct {
		Level         string `mapstructure:"level" json:"level" yaml:"level" ini:"level"`                                    // 级别
		Format        string `mapstructure:"format" json:"format" yaml:"format" ini:"format"`                                // 输出
		Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix" ini:"prefix"`                                // 日志前缀
		Director      string `mapstructure:"director" json:"director"  yaml:"director" ini:"director"`                       // 日志文件夹
		ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"showLine" ini:"showLine"`                       // 显示行
		EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level" ini:"encode-level"`         // 编码级
		StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key" ini:"stacktrace-key"` // 栈名
		LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console" ini:"log-in-console"`  // 输出控制台
	}
	Config struct {
		WebHook WebHook `mapstructure:"webhook" json:"webhook" yaml:"webhook" ini:"webhook"`
		Log     Log     `mapstructure:"log" json:"log" yaml:"log" ini:"log"`
		Email   Email   `mapstructure:"email" json:"email" yaml:"email" ini:"email"`
		System  System  `mapstructure:"system" json:"system" yaml:"system" ini:"system"`
		Mysql   Mysql   `mapstructure:"mysql" json:"mysql" yaml:"mysql" ini:"mysql"`
		Etcd    Etcd    `mapstructure:"etcd" json:"etcd" yaml:"etcd" ini:"etcd"`
	}
)

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

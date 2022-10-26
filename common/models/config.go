package models

import "fmt"

const (
	CronyNodeTableName   = "node"
	CronyJobTableName    = "job"
	CronyJobLogTableName = "job_log"
	CronyUserTableName   = "user"
	CronyScriptTableName = "script"
)

type (
	Mysql struct {
		Path         string `mapstructure:"path" json:"path" yaml:"path" ini:"path"`
		Port         string `mapstructure:"port" json:"port" yaml:"port" ini:"port"`
		Config       string `mapstructure:"config" json:"config" yaml:"config" ini:"config"`
		Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name" ini:"db-name"`
		Username     string `mapstructure:"username" json:"username" yaml:"username" ini:"username"`
		Password     string `mapstructure:"password" json:"password" yaml:"password" ini:"password"`
		MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns" ini:"max-idle-conns"`
		MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns" ini:"max-open-conns"`
		LogMode      string `mapstructure:"log-mode" json:"logMode" yaml:"log-mode" ini:"log-mode"`
		LogZap       bool   `mapstructure:"log-zap" json:"logZap" yaml:"log-zap" ini:"log-zap"`
	}
	Email struct {
		Port     int      `mapstructure:"port" json:"port" yaml:"port" ini:"port"`
		From     string   `mapstructure:"from" json:"from" yaml:"from" ini:"from"`
		Host     string   `mapstructure:"host" json:"host" yaml:"host" ini:"host"`
		IsSSL    bool     `mapstructure:"is-ssl" json:"isSSL" yaml:"is-ssl" ini:"is-ssl"`
		Secret   string   `mapstructure:"secret" json:"secret" yaml:"secret" ini:"secret"`
		Nickname string   `mapstructure:"nickname" json:"nickname" yaml:"nickname" ini:"nickname"`
		To       []string `mapstructure:"to" json:"to" yaml:"to" ini:"to"`
	}
	WebHook struct {
		Kind string `mapstructure:"kind" json:"kind" yaml:"kind" ini:"kind"`
		Url  string `mapstructure:"url" json:"url" yaml:"url" ini:"kind"`
	}
	Etcd struct {
		Endpoints   []string `mapstructure:"endpoints" json:"endpoints" yaml:"endpoints" ini:"endpoints"`
		Username    string   `mapstructure:"username" json:"username" yaml:"username" ini:"username"`
		Password    string   `mapstructure:"password" json:"password" yaml:"password" ini:"password"`
		DialTimeout int64    `mapstructure:"dial-timeout" json:"dial-timeout" yaml:"dial-timeout" ini:"dial-timeout"`
		ReqTimeout  int64    `mapstructure:"req-timeout" json:"req-timeout" yaml:"req-timeout" ini:"req-timeout"`
	}
	System struct {
		Env                string `mapstructure:"env" json:"env" yaml:"env" ini:"env"`
		Addr               int    `mapstructure:"addr" json:"addr" yaml:"addr" ini:"addr"`
		NodeTtl            int64  `mapstructure:"node-ttl" json:"node-ttl" yaml:"node-ttl" ini:"node-ttl"`
		JobProcTtl         int64  `mapstructure:"job-proc-ttl" json:"job-proc-ttl" yaml:"job-proc-ttl" ini:"job-proc-ttl"`
		Version            string `mapstructure:"version" json:"version" yaml:"version" ini:"version"`
		LogCleanPeriod     int64  `mapstructure:"log-clean-period" json:"log-clean-period" yaml:"log-clean-period" ini:"log-clean-period"`
		LogCleanExpiration int64  `mapstructure:"log-clean-expiration" json:"log-clean-expiration" yaml:"log-clean-expiration" ini:"log-clean-expiration"`
		CmdAutoAllocation  bool   `mapstructure:"cmd-auto-allocation" json:"cmd-auto-allocation" yaml:"cmd-auto-allocation" ini:"cmd-auto-allocation"`
	}
	Log struct {
		Level         string `mapstructure:"level" json:"level" yaml:"level" ini:"level"`
		Format        string `mapstructure:"format" json:"format" yaml:"format" ini:"format"`
		Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix" ini:"prefix"`
		Director      string `mapstructure:"director" json:"director"  yaml:"director" ini:"director"`
		ShowLine      bool   `mapstructure:"show-line" json:"showLine" yaml:"showLine" ini:"showLine"`
		EncodeLevel   string `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level" ini:"encode-level"`
		StacktraceKey string `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key" ini:"stacktrace-key"`
		LogInConsole  bool   `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console" ini:"log-in-console"`
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

func (m *Mysql) EmptyDsn() string {
	if m.Path == "" {
		m.Path = "127.0.0.1"
	}
	if m.Port == "" {
		m.Port = "3306"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/", m.Username, m.Password, m.Path, m.Port)
}

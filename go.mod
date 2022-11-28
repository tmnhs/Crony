module github.com/tmnhs/crony

go 1.16

require (
	github.com/coreos/bbolt v1.3.0 // indirect
	github.com/coreos/etcd v3.3.9+incompatible
	github.com/coreos/go-systemd v0.0.0-20180828140353-eee3db372b31 // indirect
	github.com/coreos/pkg v0.0.0-20180108230652-97fdf19511ea // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.5.4
	github.com/gin-gonic/gin v1.8.1
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/jakecoffman/cron v0.0.0-20190106200828-7e2009c226a5
	github.com/jessevdk/go-flags v1.5.0
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/ouqiang/goutil v1.4.1
	github.com/pkg/errors v0.9.1
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/spf13/viper v1.13.0
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20171017195756-830351dc03c6 // indirect
	github.com/xiang90/probing v0.0.0-20160813154853-07dd2e8dfe18 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
	golang.org/x/net v0.0.0-20220722155237-a158d28d115b // indirect
	golang.org/x/sync v0.0.0-20220513210516-0976fa681c29
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gorm.io/driver/mysql v1.3.6
	gorm.io/gorm v1.23.10
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

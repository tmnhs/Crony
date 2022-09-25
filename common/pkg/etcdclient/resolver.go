package etcdclient

//服务发现要实现的接口
type Watcher interface {
	Watch() error

	Close() error
}

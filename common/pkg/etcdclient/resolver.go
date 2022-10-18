package etcdclient

type Watcher interface {
	Watch() error

	Close() error
}

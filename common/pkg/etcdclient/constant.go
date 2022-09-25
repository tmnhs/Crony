package etcdclient

const (
	KeyEtcdProfile = "/crony/"
	//node节点
	KeyEtcdNode    = KeyEtcdProfile + "node/"
	KeyEtcdProc    = KeyEtcdProfile + "proc/"
	KeyEtcdCmd     = KeyEtcdProfile + "cmd/"
	KeyEtcdOnce    = KeyEtcdProfile + "once/"
	KeyEtcdLock    = KeyEtcdProfile + "lock/"
	KeyEtcdNoticer = KeyEtcdProfile + "noticer/"
)

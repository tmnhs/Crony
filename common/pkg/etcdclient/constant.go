package etcdclient

const (
	KeyEtcdProfile = "/crony/"
	//node节点
	KeyEtcdNode    = KeyEtcdProfile + "node/"
	KeyEtcdProc    = KeyEtcdProfile + "proc/"
	KeyEtcdJob     = KeyEtcdProfile + "job/"
	KeyEtcdGroup   = KeyEtcdProfile + "group/"
	KeyEtcdOnce    = KeyEtcdProfile + "once/"
	KeyEtcdLock    = KeyEtcdProfile + "lock/"
	KeyEtcdNoticer = KeyEtcdProfile + "noticer/"
)

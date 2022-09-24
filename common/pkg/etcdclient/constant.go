package etcdclient

const (
	KeyEtcdProfile = "/crony/"
	//node节点
	KeyEtcdNode    = KeyEtcdProfile + "node/%s"
	KeyEtcdProc    = KeyEtcdProfile + "proc/%s"
	KeyEtcdCmd     = KeyEtcdProfile + "cmd/%s"
	KeyEtcdOnce    = KeyEtcdProfile + "once/%s"
	KeyEtcdLock    = KeyEtcdProfile + "lock/%s"
	KeyEtcdNoticer = KeyEtcdProfile + "noticer/%s"
)

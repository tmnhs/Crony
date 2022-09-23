package etcdclient

/*
   "node": "/crony/node/",
   "proc": "/crony/proc/",
   "cmd": "/crony/cmd/",
   "once": "/crony/once/",
   "lock": "/crony/lock/",
   "group": "/crony/group/",
   "noticer": "/crony/noticer/"
*/
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

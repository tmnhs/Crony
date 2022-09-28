package etcdclient

const (
	keyEtcdProfile = "/crony/"

	//node节点
	//key /crony/node/<node_uuid>
	KeyEtcdNodeProfile = keyEtcdProfile + "node/"
	KeyEtcdNode        = KeyEtcdNodeProfile + "%s"

	//key  /crony/proc/<node_uuid>/<group_id>/<job_id>/<pid>
	KeyEtcdProcProfile = keyEtcdProfile + "proc/%s/"
	KeyEtcdProc        = KeyEtcdProcProfile + "%d/%d/%d"

	//key /cronsun/job/<node_uuid>/<group_id>/<job_id>
	KeyEtcdJobProfile = keyEtcdProfile + "job/%s/"
	KeyEtcdJob        = KeyEtcdJobProfile + "%d/%d"

	// key /cronsun/once/group/<jobID>
	KeyEtcdOnceProfile = keyEtcdProfile + "once/"
	KeyEtcdOnce        = KeyEtcdOnceProfile + "%d/%d"

	KeyEtcdLock    = keyEtcdProfile + "lock/"
	KeyEtcdNoticer = keyEtcdProfile + "noticer/"
)

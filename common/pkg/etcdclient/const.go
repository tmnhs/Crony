package etcdclient

const (
	keyEtcdProfile = "/crony/"

	//node节点
	//key /crony/node/<node_uuid>
	KeyEtcdNodeProfile = keyEtcdProfile + "node/"
	KeyEtcdNode        = KeyEtcdNodeProfile + "%s"

	//key  /crony/proc/<node_uuid>/<job_id>/<pid>
	KeyEtcdProcProfile = keyEtcdProfile + "proc/"
	KeyEtcdProc        = KeyEtcdProcProfile + "%s/%d/%d"

	//key /crony/job/<node_uuid>/<job_id>
	KeyEtcdJobProfile = keyEtcdProfile + "job/%s/"
	KeyEtcdJob        = KeyEtcdJobProfile + "%d"

	// key /crony/once/<jobID>
	KeyEtcdOnceProfile = keyEtcdProfile + "once/"
	KeyEtcdOnce        = KeyEtcdOnceProfile + "%d"
)

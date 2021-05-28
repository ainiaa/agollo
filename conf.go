package apollo

type Conf struct {
	Endpoint           string // apollo 服务地址
	DefaultNamespace   string // 默认命名空间
	AppId              string // app_id
	Cluster            string // 默认的集群名称，默认：default
	LongPollerInterval int64  // 轮训间隔时间，默认：1s
	BackupFile         string // 备份文件存放地址，默认：.go-apollo
}

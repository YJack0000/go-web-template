package ports

type (
	TwccManager interface {
		// 任務容器
		RunTwccJob(string) error
		GetTwccJobStatus(string) (string, error)
		// 開發容器
		CreateTwccCCS() (string, error)
		TwccCCSAssociateIP(string) error
		GetTwccCCSEntryPoint(string) (string, error)
		DeleteTwccCCS(string) error
	}
)

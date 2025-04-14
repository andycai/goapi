package reposync

import "time"

// RepoConfig 仓库配置
type RepoConfig struct {
	RepoType1  string `json:"repo_type1"`  // 第一个仓库类型 (svn/git)
	RepoURL1   string `json:"repo_url1"`   // 第一个仓库URL
	LocalPath1 string `json:"local_path1"` // 第一个仓库本地路径
	Username1  string `json:"username1"`   // 第一个仓库用户名
	Password1  string `json:"password1"`   // 第一个仓库密码
	RepoType2  string `json:"repo_type2"`  // 第二个仓库类型 (svn/git)
	RepoURL2   string `json:"repo_url2"`   // 第二个仓库URL
	LocalPath2 string `json:"local_path2"` // 第二个仓库本地路径
	Username2  string `json:"username2"`   // 第二个仓库用户名
	Password2  string `json:"password2"`   // 第二个仓库密码
	ConfigPath string `json:"config_path"` // 配置文件路径
}

// CommitRecord 提交记录
type CommitRecord struct {
	Revision       string       `json:"revision"`        // 版本号
	Comment        string       `json:"comment"`         // 提交内容
	Author         string       `json:"author"`          // 提交人
	Time           time.Time    `json:"time"`            // 提交时间
	Synced         bool         `json:"synced"`          // 是否已同步
	AffectedIssues string       `json:"affected_issues"` // 受影响的问题列表
	ChangedFiles   []FileChange `json:"changed_files"`   // 变更的文件列表
}

// FileChange 文件变更
type FileChange struct {
	Path       string `json:"path"`        // 文件路径
	ChangeType string `json:"change_type"` // 变更类型 (A:新增, M:修改, D:删除)
}

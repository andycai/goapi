package shell

type ScriptConfig struct {
	Path string            // 脚本文件路径
	Args []string          // 命令行参数
	Env  map[string]string // 环境变量
}

type ScriptForm struct {
	Name        string              `json:"name"`
	Repository  string              `json:"repository"`
	Platform    string              `json:"platform"`
	PublishType string              `json:"publishType"`
	Params      string              `json:"params"`
	Ext         []map[string]string `json:"ext"`
}

package browse

import "time"

// FileEntry 存储文件信息的结构体
type FileEntry struct {
	Name         string    // 文件名
	Size         int64     // 文件大小
	FormatedSize string    // 格式化后的文件大小
	ModTime      time.Time // 修改时间
	IsDir        bool      // 是否是目录
	FileType     string    // 文件类型
}

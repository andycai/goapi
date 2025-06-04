package reposync

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/andycai/goapi/models"
	"github.com/andycai/goapi/modules/interface/git"
	"github.com/andycai/goapi/modules/interface/svn"
	"github.com/andycai/goapi/pkg/utility/path"
)

var config *RepoConfig

// initService 初始化服务
func initService() {
	config = &RepoConfig{ConfigPath: "./data/reposync_config.json"}

	// 尝试加载配置
	loadConfig()
}

// saveConfig 保存配置
func saveConfig() error {
	if config == nil {
		return errors.New("配置为空")
	}

	// 确保目录存在
	dir := filepath.Dir(config.ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 序列化配置
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// 记录操作日志
	// adminlog.WriteLog(c, "write", "config", permission.ID, fmt.Sprintf("创建权限：%s", permission.Name))

	// 写入文件
	return os.WriteFile(config.ConfigPath, data, 0644)
}

// loadConfig 加载配置
func loadConfig() error {
	// 检查文件是否存在
	if _, err := os.Stat(config.ConfigPath); os.IsNotExist(err) {
		// 配置文件不存在，使用默认配置
		return nil
	}

	// 读取文件
	data, err := os.ReadFile(config.ConfigPath)
	if err != nil {
		return err
	}

	// 反序列化配置
	return json.Unmarshal(data, config)
}

// updateConfig 更新配置
func updateConfig(conf *RepoConfig) error {
	// 验证路径
	if !path.IsValid(config.LocalPath1) || !path.IsValid(config.LocalPath2) {
		return errors.New("无效的本地路径")
	}

	// 更新配置
	config.ConfigPath = conf.ConfigPath
	config = conf

	// 保存配置
	return saveConfig()
}

// getConfig 获取配置
func getConfig() *RepoConfig {
	return config
}

// checkoutRepos 初始化检出仓库
func checkoutRepos() error {
	if config == nil {
		return errors.New("配置为空")
	}

	// 检出第一个仓库
	err := checkoutRepo(
		config.RepoType1,
		config.RepoURL1,
		config.LocalPath1,
		config.Username1,
		config.Password1,
	)
	if err != nil {
		return fmt.Errorf("检出第一个仓库失败: %v", err)
	}

	// 检出第二个仓库
	err = checkoutRepo(
		config.RepoType2,
		config.RepoURL2,
		config.LocalPath2,
		config.Username2,
		config.Password2,
	)
	if err != nil {
		return fmt.Errorf("检出第二个仓库失败: %v", err)
	}

	return nil
}

// checkoutRepo 检出单个仓库
func checkoutRepo(repoType, url, repopath, username, password string) error {
	if !path.IsValid(repopath) {
		return errors.New("无效的本地路径")
	}

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(repopath), 0755); err != nil {
		return err
	}

	// 根据仓库类型执行检出操作
	switch strings.ToLower(repoType) {
	case "svn":
		return svn.Checkout(url, repopath, username, password)
	case "git":
		return git.Clone(url, repopath, "", username, password)
	default:
		return fmt.Errorf("不支持的仓库类型: %s", repoType)
	}
}

// getCommits 获取第一个仓库的提交记录
func getCommits(limit, page int) ([]CommitRecord, int, error) {
	if config == nil {
		return nil, 0, errors.New("配置为空")
	}

	// 直接从数据库获取同步记录
	var records []models.RepoSyncRecord
	var total int64

	// 获取总记录数
	if err := app.DB.Model(&models.RepoSyncRecord{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取提交记录总数失败: %v", err)
	}

	// 分页查询，按照创建时间倒序排列
	offset := (page - 1) * limit
	if err := app.DB.Order("id DESC").Offset(offset).Limit(limit).Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("查询提交记录失败: %v", err)
	}

	// 转换为CommitRecord格式
	commits := make([]CommitRecord, 0, len(records))
	for _, record := range records {
		// 解析受影响的文件列表
		var changedFiles []FileChange
		if record.AffectedFiles != "" {
			// 格式为 "A:file1.txt,file2.txt;M:file3.txt,file4.txt;D:file5.txt"
			sections := strings.Split(record.AffectedFiles, ";")
			for _, section := range sections {
				parts := strings.SplitN(section, ":", 2)
				if len(parts) != 2 {
					continue
				}

				changeType := parts[0]
				fileList := strings.Split(parts[1], ",")

				// 处理可能带有省略号和计数的情况
				lastFile := fileList[len(fileList)-1]
				if strings.Contains(lastFile, "...（共") {
					fileList = fileList[:len(fileList)-1]
				}

				for _, file := range fileList {
					changedFiles = append(changedFiles, FileChange{
						Path:       file,
						ChangeType: changeType,
					})
				}
			}
		}

		commit := CommitRecord{
			Revision:       record.Revision,
			Comment:        record.Comment,
			Author:         record.Author,
			Time:           record.SyncTime,
			Synced:         record.Status == 1, // 1表示同步成功
			AffectedIssues: record.AffectedIssues,
			ChangedFiles:   changedFiles,
		}
		commits = append(commits, commit)
	}

	return commits, int(total), nil
}

// updateRepo 更新仓库
func updateRepo(repoType, repopath, username, password string) error {
	if !path.IsValid(repopath) {
		return errors.New("无效的本地路径")
	}

	// 根据仓库类型执行更新操作
	switch strings.ToLower(repoType) {
	case "svn":
		return svn.Update(repopath)
	case "git":
		return git.Pull(repopath)
	default:
		return fmt.Errorf("不支持的仓库类型: %s", repoType)
	}
}

// getRepoCommits 获取仓库提交记录
func getRepoCommits(repoType, repopath string, limit int) ([]CommitRecord, error) {
	if !path.IsValid(repopath) {
		return nil, errors.New("无效的本地路径")
	}

	// 根据仓库类型获取提交记录
	switch strings.ToLower(repoType) {
	case "svn":
		return getSvnCommits(repopath, limit)
	case "git":
		return getGitCommits(repopath, limit)
	default:
		return nil, fmt.Errorf("不支持的仓库类型: %s", repoType)
	}
}

// getSvnCommits 获取Svn提交记录
func getSvnCommits(path string, limit int) ([]CommitRecord, error) {
	// 获取Svn日志
	logOutput, err := svn.Log(path, limit)
	if err != nil {
		return nil, err
	}

	// 解析Svn日志
	var commits []CommitRecord
	logEntries := strings.Split(logOutput, "------------------------------------------------------------------------")

	for _, entry := range logEntries {
		if strings.TrimSpace(entry) == "" {
			continue
		}

		lines := strings.Split(strings.TrimSpace(entry), "\n")
		if len(lines) < 2 {
			continue
		}

		// 解析第一行: r123 | user | 2023-01-01 12:00:00 +0800
		headerParts := strings.Split(lines[0], " | ")
		if len(headerParts) < 3 {
			continue
		}

		// 提取版本号
		revision := strings.TrimPrefix(headerParts[0], "r")

		// 提取作者
		author := headerParts[1]

		// 解析时间
		timeStr := headerParts[2]
		// 移除括号中的额外信息
		if idx := strings.Index(timeStr, " ("); idx != -1 {
			timeStr = strings.TrimSpace(timeStr[:idx])
		}
		commitTime, err := time.Parse("2006-01-02 15:04:05 -0700", timeStr)
		if err != nil {
			// 如果解析失败，使用当前时间
			commitTime = time.Now()
		}

		// 提取注释
		comment := ""
		if len(lines) > 1 {
			comment = strings.Join(lines[1:], "\n")
		}

		// 获取变更文件列表
		changedFiles, err := getSvnFileChanges(path, revision)
		if err != nil {
			changedFiles = []FileChange{} // 如果获取失败，使用空列表
		}

		// 创建提交记录
		commits = append(commits, CommitRecord{
			Revision:     revision,
			Comment:      comment,
			Author:       author,
			Time:         commitTime,
			Synced:       false,
			ChangedFiles: changedFiles,
		})
	}

	return commits, nil
}

// getSvnFileChanges 获取Svn文件变更列表
func getSvnFileChanges(path, revision string) ([]FileChange, error) {
	// 执行svn log命令获取上一个版本号
	cmd := exec.Command("svn", "log", "-r", fmt.Sprintf("%s:1", revision), "--limit", "2", path)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("获取SVN日志失败: %v", err)
	}

	// 解析日志获取版本号
	lines := strings.Split(string(output), "\n")
	var prevRevision string
	for _, line := range lines {
		// 检查是否符合SVN日志格式：r123 | user | 2023-01-01 12:00:00 +0800
		parts := strings.Split(line, " | ")
		if len(parts) >= 3 {
			// 检查版本号格式：r123
			if strings.HasPrefix(parts[0], "r") && len(parts[0]) > 1 {
				rev := strings.TrimPrefix(parts[0], "r")
				// 验证版本号是否为数字
				if _, err := strconv.Atoi(rev); err == nil {
					if rev != revision {
						prevRevision = rev
						break
					}
				}
			}
		}
	}

	if prevRevision == "" {
		// 如果是第一个版本，获取该版本的所有文件作为新增
		cmd = exec.Command("svn", "list", "-r", revision, path)
		output, err = cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("获取SVN文件列表失败: %v", err)
		}

		// 解析输出
		var changes []FileChange
		lines = strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 所有文件都标记为新增
			changes = append(changes, FileChange{
				Path:       line,
				ChangeType: "A",
			})

			// 限制最大文件数为50
			if len(changes) >= 50 {
				break
			}
		}

		return changes, nil
	}

	// 执行svn diff命令获取变更文件
	cmd = exec.Command("svn", "diff", "--summarize", "-r", fmt.Sprintf("%s:%s", prevRevision, revision), path)
	output, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("获取SVN差异失败: %v", err)
	}

	// 解析输出
	var changes []FileChange
	lines = strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 格式: A/M/D path/to/file.txt
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 2 {
			changeType := strings.TrimSpace(parts[0])
			filePath := strings.TrimSpace(parts[1])

			// 计算相对路径
			relPath, err := filepath.Rel(path, filePath)
			if err != nil {
				continue
			}

			changes = append(changes, FileChange{
				Path:       relPath,
				ChangeType: changeType,
			})

			// 限制最大文件数为50
			if len(changes) >= 50 {
				break
			}
		}
	}

	return changes, nil
}

// getGitCommits 获取Git提交记录
func getGitCommits(path string, limit int) ([]CommitRecord, error) {
	// 获取Git日志
	logOutput, err := git.Log(path, limit)
	if err != nil {
		return nil, err
	}

	// 解析Git日志
	var commits []CommitRecord
	lines := strings.Split(logOutput, "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// 格式: abcd123 Commit message
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		revision := parts[0]
		comment := parts[1]

		// 获取提交详情
		details, err := getGitCommitDetails(path, revision)
		if err != nil {
			continue
		}

		// 获取变更文件列表
		changedFiles, err := getGitFileChanges(path, revision)
		if err != nil {
			changedFiles = []FileChange{} // 如果获取失败，使用空列表
		}

		// 创建提交记录
		commits = append(commits, CommitRecord{
			Revision:     revision,
			Comment:      comment,
			Author:       details.Author,
			Time:         details.Time,
			Synced:       false,
			ChangedFiles: changedFiles,
		})
	}

	return commits, nil
}

// GitCommitDetails Git提交详情
type GitCommitDetails struct {
	Author string
	Time   time.Time
}

// getGitCommitDetails 获取Git提交详情
func getGitCommitDetails(path, revision string) (*GitCommitDetails, error) {
	// 执行git show命令获取提交详情
	cmd := exec.Command("git", "-C", path, "show", "--pretty=format:%an|%at", "-s", revision)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析输出
	parts := strings.Split(strings.TrimSpace(string(output)), "|")
	if len(parts) != 2 {
		return nil, errors.New("无法解析提交详情")
	}

	// 解析作者
	author := parts[0]

	// 解析时间
	timestamp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}
	commitTime := time.Unix(timestamp, 0)

	return &GitCommitDetails{
		Author: author,
		Time:   commitTime,
	}, nil
}

// getGitFileChanges 获取Git文件变更列表
func getGitFileChanges(path, revision string) ([]FileChange, error) {
	// 执行git show命令获取变更文件
	cmd := exec.Command("git", "-C", path, "show", "--name-status", "--pretty=format:", revision)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析输出
	var changes []FileChange
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 格式: A path/to/file.txt
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) == 2 {
			changeType := parts[0]
			filePath := parts[1]

			changes = append(changes, FileChange{
				Path:       filePath,
				ChangeType: changeType,
			})
		}
	}

	return changes, nil
}

// SyncCommits 同步选中的提交记录
func SyncCommits(revisions []string) error {
	if config == nil {
		return errors.New("配置为空")
	}

	// 对版本号进行排序（从小到大）
	sort.Slice(revisions, func(i, j int) bool {
		// 尝试将版本号转换为数字进行比较
		numI, errI := strconv.Atoi(revisions[i])
		numJ, errJ := strconv.Atoi(revisions[j])

		// 如果两个都是数字，按数字大小排序
		if errI == nil && errJ == nil {
			return numI < numJ
		}

		// 如果不是数字，按字符串排序
		return revisions[i] < revisions[j]
	})

	// 更新第一个仓库
	err := updateRepo(
		config.RepoType1,
		config.LocalPath1,
		config.Username1,
		config.Password1,
	)
	if err != nil {
		return fmt.Errorf("更新第一个仓库失败: %v", err)
	}

	// 获取所有提交记录
	commits, err := getRepoCommits(
		config.RepoType1,
		config.LocalPath1,
		100, // 获取足够多的记录
	)
	if err != nil {
		return fmt.Errorf("获取提交记录失败: %v", err)
	}

	// 筛选需要同步的提交记录
	var selectedCommits []CommitRecord
	for _, rev := range revisions {
		for _, commit := range commits {
			if commit.Revision == rev {
				selectedCommits = append(selectedCommits, commit)
				break
			}
		}
	}

	// 同步每个提交记录
	for i := range selectedCommits {
		err := syncCommit(selectedCommits[i])
		if err != nil {
			return fmt.Errorf("同步提交记录 %s 失败: %v", selectedCommits[i].Revision, err)
		}
	}

	return nil
}

// syncCommit 同步单个提交记录
func syncCommit(commit CommitRecord) error {
	// 获取变更文件列表
	changes, err := getFileChanges(
		config.RepoType1,
		config.LocalPath1,
		commit.Revision,
	)
	if err != nil {
		// 记录同步失败状态
		app.DB.Create(&models.RepoSyncRecord{
			Revision:       commit.Revision,
			Comment:        commit.Comment,
			Author:         commit.Author,
			AffectedIssues: extractIssueNumbers(commit.Comment),
			AffectedFiles:  formatAffectedFiles(commit.ChangedFiles),
			SyncTime:       time.Now(),
			Status:         2, // 同步失败
		})
		return fmt.Errorf("获取变更文件列表失败: %v", err)
	}

	// 更新第一个仓库
	err = updateRepo(
		config.RepoType2,
		config.LocalPath2,
		config.Username2,
		config.Password2,
	)
	if err != nil {
		return fmt.Errorf("更新第二个仓库失败: %v", err)
	}

	// 处理每个变更文件
	for _, change := range changes {
		// 计算源路径和目标路径
		sourcePath := filepath.Join(config.LocalPath1, change.Path)
		targetPath := filepath.Join(config.LocalPath2, change.Path)

		// 根据变更类型执行操作
		switch change.ChangeType {
		case "A", "M": // 新增或修改
			// 检查源文件是否是目录
			sourceInfo, err := os.Stat(sourcePath)
			if err != nil {
				return fmt.Errorf("获取源文件信息失败: %v", err)
			}

			if sourceInfo.IsDir() {
				// 如果是目录，只需要创建目标目录
				if err := os.MkdirAll(targetPath, 0755); err != nil {
					return err
				}
			} else {
				// 如果是文件，确保目标目录存在并复制文件
				targetDir := filepath.Dir(targetPath)
				if err := os.MkdirAll(targetDir, 0755); err != nil {
					return err
				}

				// 复制文件
				if err := copyFile(sourcePath, targetPath); err != nil {
					return err
				}
			}
		case "D": // 删除
			// 检查目标路径是否是目录
			targetInfo, err := os.Stat(targetPath)
			if err == nil && targetInfo.IsDir() {
				// 如果是目录，使用RemoveAll删除
				if err := os.RemoveAll(targetPath); err != nil {
					return err
				}
			} else {
				// 如果是文件或不存在，尝试删除
				if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
					return err
				}

				// 检查并删除空目录
				targetDir := filepath.Dir(targetPath)
				if err := removeEmptyDirs(targetDir, config.LocalPath2); err != nil {
					return err
				}
			}
		}
	}

	// 提交到第二个仓库
	err = commitToRepo(
		config.RepoType2,
		config.LocalPath2,
		fmt.Sprintf("Sync from %s: %s", commit.Revision, commit.Comment),
	)
	if err != nil {
		return err
	}

	// 解析问题编号
	affectedIssues := extractIssueNumbers(commit.Comment)

	// 格式化受影响的文件
	affectedFiles := formatAffectedFiles(changes)

	// 记录同步成功状态
	err = app.DB.Create(&models.RepoSyncRecord{
		Revision:       commit.Revision,
		Comment:        commit.Comment,
		Author:         commit.Author,
		AffectedIssues: affectedIssues,
		AffectedFiles:  affectedFiles,
		SyncTime:       time.Now(),
		Status:         1, // 同步成功
	}).Error

	if err != nil {
		return fmt.Errorf("记录同步成功状态失败: %v", err)
	}

	// 更新第一个仓库
	err = updateRepo(
		config.RepoType2,
		config.LocalPath2,
		config.Username2,
		config.Password2,
	)
	if err != nil {
		return fmt.Errorf("更新第二个仓库失败2: %v", err)
	}

	return nil
}

// getFileChanges 获取文件变更列表
func getFileChanges(repoType, path, revision string) ([]FileChange, error) {
	// 根据仓库类型获取文件变更列表
	switch strings.ToLower(repoType) {
	case "svn":
		return getSvnFileChanges(path, revision)
	case "git":
		return getGitFileChanges(path, revision)
	default:
		return nil, fmt.Errorf("不支持的仓库类型: %s", repoType)
	}
}

// commitToRepo 提交到仓库
func commitToRepo(repoType, repopath, message string) error {
	if !path.IsValid(repopath) {
		return errors.New("无效的本地路径")
	}

	// 根据仓库类型执行提交操作
	switch strings.ToLower(repoType) {
	case "svn":
		// 获取文件状态
		cmd := exec.Command("svn", "status", repopath)
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("获取文件状态失败: %v", err)
		}

		// 处理每个文件的状态
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// 格式: ?/A/M/D/! path/to/file
			parts := strings.SplitN(line, " ", 2)
			if len(parts) != 2 {
				continue
			}

			status := parts[0]
			filePath := strings.TrimSpace(parts[1])

			switch status {
			case "?":
				// 未版本控制的文件，需要添加
				if err := svn.Add(filePath); err != nil {
					return fmt.Errorf("添加文件失败 %s: %v", filePath, err)
				}
			case "!", "D":
				// 丢失或删除的文件，需要从版本库中删除
				cmd := exec.Command("svn", "delete", filePath)
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("删除文件失败 %s: %v", filePath, err)
				}
			}
		}

		// 提交所有变更
		return svn.Commit(repopath, message)
	case "git":
		// 先添加所有变更
		_, err := git.ExecGitCommand(repopath, "add", "-A")
		if err != nil {
			return err
		}
		// 提交变更
		return git.Commit(repopath, message)
	default:
		return fmt.Errorf("不支持的仓库类型: %s", repoType)
	}
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	// 打开源文件
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	// 创建目标文件
	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	// 复制内容
	_, err = io.Copy(destination, source)
	return err
}

// removeEmptyDirs 递归删除空目录
func removeEmptyDirs(dir, rootDir string) error {
	// 检查路径是否安全
	if !path.IsValid(dir) {
		return errors.New("无效的目录路径")
	}

	// 如果已经到达根目录，停止递归
	if dir == rootDir {
		return nil
	}

	// 读取目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// 如果目录不为空，保留目录
	if len(entries) > 0 {
		return nil
	}

	// 删除空目录
	if err := os.Remove(dir); err != nil {
		return err
	}

	// 递归检查父目录
	return removeEmptyDirs(filepath.Dir(dir), rootDir)
}

// SyncChangesBetweenRevisions 同步两个版本之间的变更
func SyncChangesBetweenRevisions(fromRev, toRev string) (int, error) {
	if config == nil {
		return 0, errors.New("配置为空")
	}

	// 更新第一个仓库
	err := updateRepo(
		config.RepoType1,
		config.LocalPath1,
		config.Username1,
		config.Password1,
	)
	if err != nil {
		return 0, fmt.Errorf("更新第一个仓库失败: %v", err)
	}

	// 获取两个版本之间的变更
	var cmd *exec.Cmd
	if config.RepoType1 == "svn" {
		cmd = exec.Command("svn", "diff", "--summarize", "-r",
			fmt.Sprintf("%s:%s", fromRev, toRev), config.LocalPath1)
	} else if config.RepoType1 == "git" {
		cmd = exec.Command("git", "-C", config.LocalPath1, "diff",
			"--name-status", fromRev, toRev)
	} else {
		return 0, fmt.Errorf("不支持的仓库类型: %s", config.RepoType1)
	}

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("获取变更失败: %v", err)
	}

	// 解析变更
	var changes []FileChange
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var changeType, filePath string
		if config.RepoType1 == "svn" {
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				changeType = strings.TrimSpace(parts[0])
				filePath = strings.TrimSpace(parts[1])
				relPath, err := filepath.Rel(config.LocalPath1, filePath)
				if err == nil {
					filePath = relPath
				}
			}
		} else { // git
			parts := strings.SplitN(line, "\t", 2)
			if len(parts) == 2 {
				changeType = parts[0]
				filePath = parts[1]
			}
		}

		if changeType != "" && filePath != "" {
			changes = append(changes, FileChange{
				Path:       filePath,
				ChangeType: changeType,
			})
		}
	}

	// 同步变更
	// 更新目标仓库
	if err := updateRepo(
		config.RepoType2,
		config.LocalPath2,
		config.Username2,
		config.Password2,
	); err != nil {
		return 0, fmt.Errorf("更新目标仓库失败: %v", err)
	}

	// 处理每个变更
	for _, change := range changes {
		sourcePath := filepath.Join(config.LocalPath1, change.Path)
		targetPath := filepath.Join(config.LocalPath2, change.Path)

		switch change.ChangeType {
		case "A", "M": // 新增或修改
			sourceInfo, err := os.Stat(sourcePath)
			if err != nil {
				continue // 跳过不存在的文件
			}

			if sourceInfo.IsDir() {
				os.MkdirAll(targetPath, 0755)
			} else {
				os.MkdirAll(filepath.Dir(targetPath), 0755)
				copyFile(sourcePath, targetPath)
			}
		case "D": // 删除
			os.RemoveAll(targetPath)
			removeEmptyDirs(filepath.Dir(targetPath), config.LocalPath2)
		}
	}

	// 提交变更
	message := fmt.Sprintf("Sync changes between %s and %s", fromRev, toRev)
	if err := commitToRepo(config.RepoType2, config.LocalPath2, message); err != nil {
		return 0, fmt.Errorf("提交变更失败: %v", err)
	}

	return len(changes), nil
}

// FindUnsyncedRevisionRange 查找未同步的版本范围
func FindUnsyncedRevisionRange() (string, string, error) {
	if config == nil {
		return "", "", errors.New("配置为空")
	}

	// 更新第一个仓库
	err := updateRepo(
		config.RepoType1,
		config.LocalPath1,
		config.Username1,
		config.Password1,
	)
	if err != nil {
		return "", "", fmt.Errorf("更新第一个仓库失败: %v", err)
	}

	// 获取最近100条提交记录
	commits, err := getRepoCommits(
		config.RepoType1,
		config.LocalPath1,
		100,
	)
	if err != nil {
		return "", "", fmt.Errorf("获取提交记录失败: %v", err)
	}

	if len(commits) == 0 {
		return "", "", nil
	}

	// 找到未同步的最小和最大revision
	var minRev, maxRev string
	for _, commit := range commits {
		// 检查是否已同步
		var record models.RepoSyncRecord
		result := app.DB.Where("revision = ?", commit.Revision).First(&record)

		if result.Error != nil { // 未找到记录，说明未同步
			if minRev == "" {
				minRev = commit.Revision
			}
			maxRev = commit.Revision
		}
	}

	return minRev, maxRev, nil
}

// formatAffectedFiles 格式化受影响的文件列表
func formatAffectedFiles(changes []FileChange) string {
	if len(changes) == 0 {
		return ""
	}

	// 按类型分组文件
	filesByType := make(map[string][]string)
	for _, change := range changes {
		filesByType[change.ChangeType] = append(filesByType[change.ChangeType], change.Path)
	}

	// 构建格式化的字符串
	var parts []string

	// 按照A（新增）、M（修改）、D（删除）的顺序处理
	changeTypes := []string{"A", "M", "D"}
	for _, changeType := range changeTypes {
		files := filesByType[changeType]
		if len(files) > 0 {
			// 最多显示10个文件，避免字符串过长
			displayCount := len(files)
			if displayCount > 10 {
				displayCount = 10
			}

			formattedPart := fmt.Sprintf("%s:%s", changeType, strings.Join(files[:displayCount], ","))

			// 如果有更多文件，添加省略标记
			if len(files) > 10 {
				formattedPart += fmt.Sprintf("...（共%d个）", len(files))
			}

			parts = append(parts, formattedPart)
		}
	}

	return strings.Join(parts, ";")
}

// RefreshCommits 刷新提交记录
func RefreshCommits(limit int) error {
	if config == nil {
		return errors.New("配置为空")
	}

	// 更新第一个仓库
	err := updateRepo(
		config.RepoType1,
		config.LocalPath1,
		config.Username1,
		config.Password1,
	)
	if err != nil {
		return fmt.Errorf("更新第一个仓库失败: %v", err)
	}

	// 获取最近的提交记录
	commits, err := getRepoCommits(
		config.RepoType1,
		config.LocalPath1,
		limit,
	)
	if err != nil {
		return fmt.Errorf("获取提交记录失败: %v", err)
	}

	// 对提交记录按照 revision 从小到大排序
	sort.Slice(commits, func(i, j int) bool {
		// 尝试将版本号转换为数字进行比较
		numI, errI := strconv.Atoi(commits[i].Revision)
		numJ, errJ := strconv.Atoi(commits[j].Revision)

		// 如果两个都是数字，按数字大小排序
		if errI == nil && errJ == nil {
			return numI < numJ
		}

		// 如果不是数字，按字符串排序
		return commits[i].Revision < commits[j].Revision
	})

	// 遍历提交记录，将未同步的记录添加到数据库
	for _, commit := range commits {
		// 检查是否已存在
		var record models.RepoSyncRecord
		result := app.DB.Where("revision = ?", commit.Revision).First(&record)

		// 解析提交消息中的问题编号
		affectedIssues := extractIssueNumbers(commit.Comment)

		// 获取变更文件列表并格式化
		var affectedFiles string
		if len(commit.ChangedFiles) > 0 {
			affectedFiles = formatAffectedFiles(commit.ChangedFiles)
		} else {
			// 如果提交记录中没有变更文件信息，尝试从仓库获取
			changes, err := getFileChanges(
				config.RepoType1,
				config.LocalPath1,
				commit.Revision,
			)
			if err == nil && len(changes) > 0 {
				affectedFiles = formatAffectedFiles(changes)
			}
		}

		if result.Error != nil { // 未找到记录，说明未同步
			// 创建新记录
			err = app.DB.Create(&models.RepoSyncRecord{
				Revision:       commit.Revision,
				Comment:        commit.Comment,
				Author:         commit.Author,
				SyncTime:       time.Now(),
				Status:         0, // 未同步状态
				AffectedIssues: affectedIssues,
				AffectedFiles:  affectedFiles,
			}).Error

			if err != nil {
				return fmt.Errorf("创建同步记录失败: %v", err)
			}
		} else {
			// 需要更新的字段
			updates := make(map[string]interface{})

			// 更新问题列表（如果为空且新解析出的不为空）
			if record.AffectedIssues == "" && affectedIssues != "" {
				updates["affected_issues"] = affectedIssues
			}

			// 更新文件列表（如果为空且新解析出的不为空）
			if record.AffectedFiles == "" && affectedFiles != "" {
				updates["affected_files"] = affectedFiles
			}

			// 如果有需要更新的字段
			if len(updates) > 0 {
				err = app.DB.Model(&record).Updates(updates).Error
				if err != nil {
					return fmt.Errorf("更新同步记录失败: %v", err)
				}
			}
		}
	}

	return nil
}

// extractIssueNumbers 从提交消息中提取问题编号
// 支持格式: #123, fix #123, close #123, fixed #123, closes #123, JIRA-123, BUG-123, ISSUE-123
func extractIssueNumbers(comment string) string {
	var issues []string

	// 匹配 #123 格式
	hashTagRegex := regexp.MustCompile(`#(\d+)`)
	hashMatches := hashTagRegex.FindAllStringSubmatch(comment, -1)
	for _, match := range hashMatches {
		if len(match) > 1 {
			issues = append(issues, "#"+match[1])
		}
	}

	// 匹配 JIRA-123, BUG-123, ISSUE-123 等格式
	issueKeyRegex := regexp.MustCompile(`(JIRA|BUG|ISSUE|TASK)-(\d+)`)
	keyMatches := issueKeyRegex.FindAllStringSubmatch(comment, -1)
	for _, match := range keyMatches {
		if len(match) > 2 {
			issues = append(issues, match[1]+"-"+match[2])
		}
	}

	// 去重
	uniqueIssues := make(map[string]bool)
	var result []string
	for _, issue := range issues {
		if !uniqueIssues[issue] {
			uniqueIssues[issue] = true
			result = append(result, issue)
		}
	}

	// 最多保存前5个问题编号
	if len(result) > 5 {
		result = result[:5]
	}

	return strings.Join(result, ", ")
}

// ClearSyncData 清空同步数据
func ClearSyncData() error {
	// 使用 model 方式删除所有同步记录
	if err := app.DB.Unscoped().Where("1 = 1").Delete(&models.RepoSyncRecord{}).Error; err != nil {
		return fmt.Errorf("清空同步数据失败: %v", err)
	}

	return nil
}

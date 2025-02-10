package core

import (
	"flag"

	"github.com/BurntSushi/toml"
)

type Config struct {
	App       AppConfig      `toml:"app"`
	Server    ServerConfig   `toml:"server"`
	Database  DatabaseConfig `toml:"database"`
	JSONPaths JSONPathConfig `toml:"json_paths"`
	FTP       FTPConfig      `toml:"ftp"`
	Auth      AuthConfig     `toml:"auth"`
	Cors      CorsConfig     `toml:"cors"`
}

type ServerConfig struct {
	Host         string             `toml:"host"`
	Port         int                `toml:"port"`
	Output       string             `toml:"output"`
	ScriptPath   string             `toml:"script_path"`
	StaticPaths  []StaticPathConfig `toml:"static_paths"`
	UserDataPath string             `toml:"user_data_path"`
	CDNPath      string             `toml:"cdn_path"`
	CDN2Path     string             `toml:"cdn2_path"`
}

type DatabaseConfig struct {
	Driver          string `toml:"driver"`
	DSN             string `toml:"dsn"`
	MaxOpenConns    int    `toml:"max_open_conns"`
	MaxIdleConns    int    `toml:"max_idle_conns"`
	ConnMaxLifetime int64  `toml:"conn_max_lifetime"`
}

type FTPConfig struct {
	Host       string `toml:"host"`
	Port       string `toml:"port"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	APKPath    string `toml:"apk_path"`
	ZIPPath    string `toml:"zip_path"`
	LogDir     string `toml:"log_dir"`
	MaxLogSize int64  `toml:"max_log_size"`
}

type JSONPathConfig struct {
	ServerList string `toml:"server_list"`
	LastServer string `toml:"last_server"`
	ServerInfo string `toml:"server_info"`
	NoticeList string `toml:"notice_list"`
	NoticeNum  string `toml:"notice_num"`
}

type StaticPathConfig struct {
	Route string `toml:"route"`
	Path  string `toml:"path"`
}

type AuthConfig struct {
	JWTSecret   string `toml:"jwt_secret"`
	TokenExpire int    `toml:"token_expire"`
}

type AppConfig struct {
	IsDev    bool `toml:"is_dev"`    // 是否为开发环境
	IsSecure bool `toml:"is_secure"` // 是否启用安全模式
}

type CorsConfig struct {
	Enabled          bool     `toml:"enabled"`           // 是否启用跨域
	AllowedOrigins   []string `toml:"allowed_origins"`   // 允许的源
	AllowedMethods   []string `toml:"allowed_methods"`   // 允许的 HTTP 方法
	AllowedHeaders   []string `toml:"allowed_headers"`   // 允许的请求头
	AllowCredentials bool     `toml:"allow_credentials"` // 是否允许携带认证信息
	MaxAge           int      `toml:"max_age"`           // 预检请求结果缓存时间（小时）
}

var config Config

func LoadConfig() error {
	// 定义配置文件路径参数
	configPath := flag.String("config", "conf.toml", "配置文件路径")
	host := flag.String("host", "", "主机地址")
	port := flag.Int("port", 0, "端口号")
	output := flag.String("output", "", "输出目录")
	scriptPath := flag.String("script_path", "", "脚本路径")
	userDataPath := flag.String("user_data_path", "", "用户数据路径")
	dbPath := flag.String("db", "", "数据库路径")
	ftpHost := flag.String("ftp_host", "", "FTP主机地址")
	ftpPort := flag.String("ftp_port", "", "FTP端口")
	ftpUser := flag.String("ftp_user", "", "FTP用户名")
	ftpPass := flag.String("ftp_pass", "", "FTP密码")
	ftpApkPath := flag.String("ftp_apk_path", "", "FTP APK上传路径")
	ftpZipPath := flag.String("ftp_zip_path", "", "FTP ZIP上传路径")
	isDev := flag.Bool("dev", false, "是否为开发环境")

	flag.Parse()

	if _, err := toml.DecodeFile(*configPath, &config); err != nil {
		return err
	}

	// 设置默认值
	if config.Database.MaxOpenConns == 0 {
		config.Database.MaxOpenConns = 100 // 默认最大连接数
	}
	if config.Database.MaxIdleConns == 0 {
		config.Database.MaxIdleConns = 10 // 默认最大空闲连接数
	}
	if config.Database.ConnMaxLifetime == 0 {
		config.Database.ConnMaxLifetime = 3600 // 默认连接生命周期为1小时
	}

	// 命令行参数覆盖配置文件
	if *host != "" {
		config.Server.Host = *host
	}
	if *port != 0 {
		config.Server.Port = *port
	}
	if *output != "" {
		config.Server.Output = *output
	}
	if *scriptPath != "" {
		config.Server.ScriptPath = *scriptPath
	}
	if *userDataPath != "" {
		config.Server.UserDataPath = *userDataPath
	}
	if *dbPath != "" {
		config.Database.DSN = *dbPath
	}
	if *ftpHost != "" {
		config.FTP.Host = *ftpHost
	}
	if *ftpPort != "" {
		config.FTP.Port = *ftpPort
	}
	if *ftpUser != "" {
		config.FTP.User = *ftpUser
	}
	if *ftpPass != "" {
		config.FTP.Password = *ftpPass
	}
	if *ftpApkPath != "" {
		config.FTP.APKPath = *ftpApkPath
	}
	if *ftpZipPath != "" {
		config.FTP.ZIPPath = *ftpZipPath
	}
	if *isDev {
		config.App.IsDev = true
	}

	return nil
}

func GetConfig() Config {
	return config
}

func GetServerConfig() ServerConfig {
	return config.Server
}

func GetDatabaseConfig() DatabaseConfig {
	return config.Database
}

func GetFTPConfig() FTPConfig {
	return config.FTP
}

func GetJSONPathConfig() JSONPathConfig {
	return config.JSONPaths
}

func GetCorsConfig() CorsConfig {
	return config.Cors
}

func UpdateServerConfig(newConfig ServerConfig) {
	config.Server = newConfig
}

func UpdateDatabaseConfig(newConfig DatabaseConfig) {
	config.Database = newConfig
}

func UpdateFTPConfig(newConfig FTPConfig) {
	config.FTP = newConfig
}

// IsDevelopment 返回是否为开发环境
func IsDevelopment() bool {
	return config.App.IsDev
}

// IsSecureMode 返回是否为安全模式
func IsSecureMode() bool {
	return config.App.IsSecure
}

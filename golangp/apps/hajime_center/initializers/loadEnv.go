package initializers

import (
	"hajime/golangp/common/logging"
	"time"

	"github.com/spf13/viper"
)

type MinioConfig struct {
	MinioAccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	MinioBucketUrl string `mapstructure:"MINIO_BUCKET_URL"`
	MinioBucket    string `mapstructure:"MINIO_BUCKET"`
}

type LocalStorageConfig struct {
	LocalStoragePath string `mapstructure:"LOCAL_STORAGE_PATH"`
}

type Config struct {
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	ServerPort     string `mapstructure:"PORT"`

	DBHostDify         string `mapstructure:"POSTGRES_HOST_DIFY"`
	DBUserNameDify     string `mapstructure:"POSTGRES_USER_DIFY"`
	DBUserPasswordDify string `mapstructure:"POSTGRES_PASSWORD_DIFY"`
	DBNameDify         string `mapstructure:"POSTGRES_DB_DIFY"`
	DBPortDify         string `mapstructure:"POSTGRES_PORT_DIFY"`

	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`
	Domain       string `mapstructure:"DOMAIN"`

	TokenSecret    string        `mapstructure:"TOKEN_SECRET"`
	TokenExpiresIn time.Duration `mapstructure:"TOKEN_EXPIRED_IN"`
	TokenMaxAge    int           `mapstructure:"TOKEN_MAXAGE"`

	EmailFrom string `mapstructure:"EMAIL_FROM"`
	SMTPHost  string `mapstructure:"SMTP_HOST"`
	SMTPPass  string `mapstructure:"SMTP_PASS"`
	SMTPPort  int    `mapstructure:"SMTP_PORT"`
	SMTPUser  string `mapstructure:"SMTP_USER"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	AppTokenExpiresIn      time.Duration `mapstructure:"APP_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	DifyConsoleEmail       string `mapstructure:"DIFY_CONSOLE_UMAIL"`
	DifyConsolePassword    string `mapstructure:"DIFY_CONSOLE_PASSWORD"`
	DifyEditorEmail        string `mapstructure:"DIFY_EDITOR_UMAIL"`
	DifyEditorPassword     string `mapstructure:"DIFY_EDITOR_PASSWORD"`
	DifyUserEmail          string `mapstructure:"DIFY_USER_UMAIL"`
	DifyUserPassword       string `mapstructure:"DIFY_USER_PASSWORD"`
	DifyConsoleStoragePath string `mapstructure:"DIFY_CONSOLE_STORAGE_PATH"`
	DifyHost               string `mapstructure:"DIFY_HOST"`
	DifyApiKey             string `mapstructure:"DIFY_API_KEY"`

	AiServerHost string `mapstructure:"AI_SERVER_HOST"`
	AiServerPort string `mapstructure:"AI_SERVER_PORT"`

	StorageType     string `mapstructure:"STORAGE_TYPE"`
	MaxUploadSize   int64  `mapstructure:"MAX_UPLOAD_SIZE"`
	MaxUploadNumber int    `mapstructure:"MAX_UPLOAD_NUMBER"`

	// xminio storage
	Minio MinioConfig

	// local storage
	LocalStorage LocalStorageConfig

	// thread number
	ThreadNumber int `mapstructure:"THREAD_NUMBER"`

	// chat config
	ApiKey string `mapstructure:"API_KEY"`
	// openai提供的接口 空字符串使用默认接口
	ApiURL string `mapstructure:"API_URL"`
	// 监听接口
	Listen string `mapstructure:"LISTEN"`
	// 代理
	Proxy         string   `mapstructure:"PROXY"`
	AdminEmail    []string `mapstructure:"ADMIN_EMAIL"`
	AdminPassword string   `mapstructure:"ADMIN_PASSWORD"`
}

func LoadEnv(path string) (config Config, err error) {

	viper.AddConfigPath(path)
	viper.AddConfigPath("golangp/apps/hajime_center") // for bazel run

	//Config file names search in order: app.env, app.dev.env
	configNames := []string{"app", "app.dev"}

	viper.SetConfigType("env")
	viper.AutomaticEnv()

	for _, configName := range configNames {
		viper.SetConfigName(configName)
		err = viper.ReadInConfig()
		if err == nil {
			logging.Info("Using config file: %s", viper.ConfigFileUsed())
			break
		}
	}

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	if config.StorageType == "s3" {
		err = viper.Unmarshal(&config.Minio)
	} else {
		err = viper.Unmarshal(&config.LocalStorage)
	}
	if err != nil {
		return
	}

	return
}

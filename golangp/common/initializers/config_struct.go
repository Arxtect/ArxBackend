package initializers

import "time"

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
	DBSslMode      string `mapstructure:"POSTGRES_SSL_MODE"`

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

	DifyHost               string   `mapstructure:"DIFY_HOST"`
	DifyApiKey             string   `mapstructure:"DIFY_API_KEY"`
	DifyConsoleEmail       string   `mapstructure:"DIFY_CONSOLE_UMAIL"`
	DifyConsolePassword    string   `mapstructure:"DIFY_CONSOLE_PASSWORD"`
	DifyConsoleStoragePath string   `mapstructure:"DIFY_CONSOLE_STORAGE_PATH"`
	DifyShowAppIDList      []string `mapstructure:"DIFY_SHOW_APP_ID_LIST"`

	StorageType     string `mapstructure:"STORAGE_TYPE"`
	MaxUploadSize   int64  `mapstructure:"MAX_UPLOAD_SIZE"`
	MaxUploadNumber int    `mapstructure:"MAX_UPLOAD_NUMBER"`

	// xminio storage
	Minio MinioConfig

	// local storage
	LocalStorage LocalStorageConfig

	// thread number
	ThreadNumber int `mapstructure:"THREAD_NUMBER"`

	// wechat conf
	OfficialAccountAppid         string `mapstructure:"OFFICIAL_ACCOUNT_APPID"`
	OfficialAccountAppSecret     string `mapstructure:"OFFICIAL_ACCOUNT_APP_SECRET"`
	OfficialAccountRedisAddr     string `mapstructure:"OFFICIAL_ACCOUNT_REDIS_ADDR"`
	OfficialAccountMessageToken  string `mapstructure:"OFFICIAL_ACCOUNT_MESSAGE_TOKEN"`
	OfficialAccountMessageAesKey string `mapstructure:"OFFICIAL_ACCOUNT_MESSAGE_AES_KEY"`

	PaymentRedisAddr          string `mapstructure:"PAYMENT_REDIS_ADDR"`
	PaymentAppID              string `mapstructure:"PAYMENT_APP_ID"`
	PaymentMchID              string `mapstructure:"PAYMENT_MCH_ID"`
	PaymentMchApiV3Key        string `mapstructure:"PAYMENT_MCH_API_V3_KEY"`
	PaymentKey                string `mapstructure:"PAYMENT_KEY"`
	PaymentCertPath           string `mapstructure:"PAYMENT_CERT_PATH"`
	PaymentKeyPath            string `mapstructure:"PAYMENT_KEY_PATH"`
	PaymentSerialNo           string `mapstructure:"PAYMENT_SERIAL_NO"`
	PaymentCertificateKeyPath string `mapstructure:"PAYMENT_CERTIFICATE_KEY_PATH"`
	PaymentWechatPaySerial    string `mapstructure:"PAYMENT_WECHAT_PAY_SERIAL"`
	PaymentRSAPublicKeyPath   string `mapstructure:"PAYMENT_RSA_PUBLIC_KEY_PATH"`
	PaymentNotifyURL          string `mapstructure:"PAYMENT_NOTIFY_URL"`
	PaymentSubMchID           string `mapstructure:"PAYMENT_SUB_MCH_ID"`
	PaymentSubAppID           string `mapstructure:"PAYMENT_SUB_APP_ID"`

	AdminPassword string   `mapstructure:"ADMIN_PASSWORD"`
	AdminEmail    []string `mapstructure:"ADMIN_EMAIL"`

	// Grpc
	GrpcServer string `mapstructure:"GRPC_SERVER"`

	// Data Service
	DataServiceHost string `mapstructure:"DATA_SERVICE_HOST"`
}

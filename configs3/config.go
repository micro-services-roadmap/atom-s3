package configs3

type Local struct {
	Path      string `mapstructure:"path" json:"path" yaml:"path"`                   // 本地文件访问路径
	StorePath string `mapstructure:"store-path" json:"store-path" yaml:"store-path"` // 本地文件存储路径
}

type S3 struct {
	Endpoint  string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	Bucket    string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Region    string `mapstructure:"region" json:"region,optional" yaml:"region"`
	SecretID  string `mapstructure:"secret-id" json:"secret-id" yaml:"secret-id"`
	SecretKey string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	BaseURL   string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
}

type AwsS3 struct {
	S3 // `yaml:",inline" mapstructure:",squash"`

	PathPrefix       string `mapstructure:"path-prefix" json:"path-prefix,optional" yaml:"path-prefix"`
	S3ForcePathStyle bool   `mapstructure:"s3-force-path-style" json:"s3-force-path-style,optional" yaml:"s3-force-path-style"`
	DisableSSL       bool   `mapstructure:"disable-ssl" json:"disable-ssl,optional" yaml:"disable-ssl"`
}

type TencentCOS struct {
	S3

	PathPrefix string `mapstructure:"path-prefix" json:"path-prefix" yaml:"path-prefix"`
}

type AliyunOSS struct {
	S3

	BasePath string `mapstructure:"base-path" json:"base-path" yaml:"base-path"`
}

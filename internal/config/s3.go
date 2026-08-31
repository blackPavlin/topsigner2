package config

type S3Config struct {
	Endpoint    string `envconfig:"ENDPOINT"     reqired:"true"`
	Region      string `envconfig:"REGION"       required:"true"`
	AccessKey   string `envconfig:"ACCESS_KEY"   required:"true"`
	SecretKey   string `envconfig:"SECRET_KEY"   required:"true"`
	ImageBucket string `envconfig:"IMAGE_BUCKET" required:"true"`
	FontBucket  string `envconfig:"FONT_BUCKET"  required:"true"`
	Secure      bool   `envconfig:"SECURE"       default:"true"`
}

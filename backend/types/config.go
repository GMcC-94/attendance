package types

type Config struct {
	JWTSecret          string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPass             string
	DBName             string
	AWSRegion          string
	AWSBucketName      string
	AWSSecretAccessKey string
	AWSAccessKeyID     string
}

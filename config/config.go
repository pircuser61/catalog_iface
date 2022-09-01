package config

const (
	UseKafka      = false
	HttpAddr      = ":8080"
	GrpcAddr      = ":8081"
	GrpcStoreAddr = ":8082"

	Topic_create = "good_create"
	Topic_update = "good_update"
	Topic_delete = "good_delete"
	Topic_error  = "errors"

	JaegerHostPort = "localhost:6831"

	RedisAddr       = "localhost:6379"
	RedisResponseDb = 1
	RedisPassword   = ""
)

var (
	KafkaBrokers = []string{"localhost:19091", "localhost:29091", "localhost:39091"}
)

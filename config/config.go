package config

const (
	UseKafka      = true
	HttpAddr      = ":8080"
	GrpcAddr      = ":8081"
	GrpcStoreAddr = ":8082"
)

var (
	KafkaBrokers = []string{"localhost:19091", "localhost:29091", "localhost:39091"}
)

const (
	Topic_create = "good_create"
	Topic_update = "good_update"
	Topic_delete = "good_delete"
	Topic_error  = "errors"
)

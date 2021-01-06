module github.com/go-ceres/ceres-config-etcd

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/go-ceres/ceres-config v1.0.2
	github.com/go-ceres/ceres-error v1.0.2
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/uuid v1.1.3 // indirect
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0 // indirect
	google.golang.org/grpc v1.34.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

go 1.15

package CeresConfigEtcd

import (
	CeresConfig "github.com/go-ceres/ceres-config"
	CeresError "github.com/go-ceres/ceres-error"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var DefaultPrefix = "/ceres/config/"

type etcdSource struct {
	client  *clientv3.Client
	config  *Config
	changed chan struct{}
	err     error
}

func (e *etcdSource) Read() (*CeresConfig.DataSet, error) {
	if e.err != nil {
		return nil, e.err
	}
	kv := clientv3.NewKV(e.client)
	res, err := kv.Get(e.config.Ctx, e.config.Prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	data := makeMapData(res.Kvs, Unmarshals[e.config.Encoding], e.config.TrimPrefix)
	b, err := Marshals[e.config.Encoding](data)
	if err != nil {
		return nil, CeresError.New("error reading source: " + err.Error())
	}

	cs := &CeresConfig.DataSet{
		Format:    e.getUnmarshal(),
		Source:    e.String(),
		Timestamp: time.Now(),
		Data:      b,
	}
	return cs, nil
}

func (e *etcdSource) Write(set *CeresConfig.DataSet) error {
	panic("implement me")
}

func (e *etcdSource) IsChanged() <-chan struct{} {
	return e.changed
}

// 开启监听
func (e *etcdSource) Watch() {
	if e.changed == nil {
		e.changed = make(chan struct{})
	}
	go func() {
		ch := e.client.Watch(e.config.Ctx, e.config.Prefix, clientv3.WithPrefix())
		for {
			select {
			case _, ok := <-ch:
				if !ok {
					return
				}
				select {
				case e.changed <- struct{}{}:
				default:
				}
			}
		}
	}()
}

func (e *etcdSource) String() string {
	return "etcd"
}

func (e *etcdSource) UnWatch() {
	close(e.changed)
	e.changed = nil
	e.client.Watcher.Close()
}

func (e *etcdSource) getUnmarshal() string {
	if e.config.Encoding != "" {
		return e.config.Encoding
	}
	return "json"
}

// 创建一个资源
func NewSource(opts ...Option) CeresConfig.Source {
	conf := defaultConfig()
	for _, opt := range opts {
		opt(conf)
	}
	// 删除掉key前缀
	if conf.TrimPrefix == "" {
		conf.TrimPrefix = conf.Prefix
	}

	cli, err := clientv3.New(*conf.Config)
	return &etcdSource{
		client: cli,
		config: conf,
		err:    err,
	}
}

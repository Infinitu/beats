package length

import (
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/elastic/beats/metricbeat/mb/parse"
	"github.com/elastic/beats/metricbeat/module/redis"

	rd "github.com/garyburd/redigo/redis"
)

var (
	debugf = logp.MakeDebug("redis-length")
)
var targetKeys []string = nil

func init() {
	if err := mb.Registry.AddMetricSet("redis", "length", New, parse.PassThruHostParser); err != nil {
		panic(err)
	}
}

// MetricSet for fetching Redis server information and statistics.
type MetricSet struct {
	mb.BaseMetricSet
	pool *rd.Pool
}

// New creates new instance of MetricSet
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	// Unpack additional configuration options.
	config := struct {
		IdleTimeout time.Duration `config:"idle_timeout"`
		Network     string        `config:"network"`
		MaxConn     int           `config:"maxconn" validate:"min=1"`
		Password    string        `config:"password"`
		LengthTargetKeys    []string        `config:"length_target_keys"`
	}{
		Network:  "tcp",
		MaxConn:  10,
		Password: "",
		LengthTargetKeys: []string{},
	}
	err := base.Module().UnpackConfig(&config)
	if err != nil {
		return nil, err
	}

	targetKeys = config.LengthTargetKeys
	return &MetricSet{
		BaseMetricSet: base,
		pool: redis.CreatePool(base.Host(), config.Password, config.Network,
			config.MaxConn, config.IdleTimeout, base.Module().Config().Timeout),
	}, nil
}

// Fetch fetches metrics from Redis by issuing the INFO command.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	// Fetch default INFO.
	

	values := map[string]int{}
	for _, key := range targetKeys {
		length, err := redis.FetchLength(key, m.pool.Get())
		if err != nil {
			return nil, err
		}
		debugf("Redis LLEN from %s: %d", m.Host(), length)

		values[key] = length
	}

	
	return eventsMapping(values), nil
}

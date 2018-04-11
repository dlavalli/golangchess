package beater

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
    "bytes"
	// "fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"
	"github.com/garyburd/redigo/redis"
	"github.com/tata/smshubbeat/config"
)

type Smshubbeat struct {
	period        time.Duration
	host          string
	port          int
	dbid          int
	luascript     string
	network       string
	maxconn       int
	auth          bool
	pass          string
	settings      config.ConfigSettings
	client        publisher.Client
	redisPool     *redis.Pool
	done          chan struct{}
}

// Load the configuration file information
func (bt *Smshubbeat) Config(b *beat.Beat, cfg *common.Config) error {
	err := cfg.Unpack(&bt.settings)
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return err
	}

	if bt.settings.BeatSettings.Period != nil {
		bt.period = time.Duration(*bt.settings.BeatSettings.Period) * time.Second
	} else {
		bt.period = config.DEFAULT_PERIOD
	}

	if bt.settings.RedisSettings.Host != nil {
		bt.host = *bt.settings.RedisSettings.Host
	} else {
		bt.host = config.DEFAULT_HOST
	}

	if bt.settings.RedisSettings.Port != nil {
		bt.port = *bt.settings.RedisSettings.Port
	} else {
		bt.port = config.DEFAULT_PORT
	}

	if bt.settings.RedisSettings.Dbid != nil {
		bt.dbid = *bt.settings.RedisSettings.Dbid
	} else {
		bt.dbid = config.DEFAULT_DBID
	}

	var script []byte
	if bt.settings.RedisSettings.Luascript != nil {
		script, err = ioutil.ReadFile(*bt.settings.RedisSettings.Luascript) // just pass the file name
	} else {
		script, err = ioutil.ReadFile(config.DEFAULT_LUA_SCRIPT) // just pass the file name
	}
	if err != nil {
		bt.luascript = ""
	} else {
		bt.luascript = string(script)
	}

	if bt.settings.RedisSettings.Network != nil {
		bt.network = *bt.settings.RedisSettings.Network
	} else {
		bt.network = config.DEFAULT_NETWORK
	}

	if bt.settings.RedisSettings.Maxconn != nil {
		bt.maxconn = *bt.settings.RedisSettings.Maxconn
	} else {
		bt.maxconn = config.DEFAULT_MAX_CONN
	}

	if bt.settings.RedisSettings.Auth.Required != nil {
		bt.auth = *bt.settings.RedisSettings.Auth.Required
	} else {
		bt.auth = config.DEFAULT_AUTH_REQUIRED
	}

	if bt.settings.RedisSettings.Auth.Requiredpass != nil {
		bt.pass = *bt.settings.RedisSettings.Auth.Requiredpass
	} else {
		bt.pass = config.DEFAULT_AUTH_REQUIRED_PASS
	}

	logp.Debug("smshubbeat", "Redis Init smshubbeat")
	logp.Debug("smshubbeat", "Redis Period %v", bt.period)
	logp.Debug("smshubbeat", "Redis Host %v", bt.host)
	logp.Debug("smshubbeat", "Redis Port %v", bt.port)
	logp.Debug("smshubbeat", "Redis DBId %v", bt.dbid)
	logp.Debug("smshubbeat", "Redis Network %v", bt.network)
	logp.Debug("smshubbeat", "Redis Max Connections %v", bt.maxconn)
	logp.Debug("smshubbeat", "Redis Auth %t", bt.auth)

	return nil
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (be beat.Beater, err error) {
	bt := &Smshubbeat{}
	if err = bt.Config(b, cfg); err != nil {
		logp.Err("Config error")
	}
	return bt, err
}

// Setup beater resources
func (bt *Smshubbeat) Setup(b *beat.Beat) error {
	bt.client = b.Publisher.Connect()
	bt.done = make(chan struct{})

	// Set up redis pool (default inmplementation without dbid selection
	redisPool := redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial(bt.network, bt.host+":"+strconv.Itoa(bt.port))

		if err != nil {
			return nil, err
		}

		return c, err
	}, bt.maxconn)
	bt.redisPool = redisPool

	if bt.auth {
		c := bt.redisPool.Get()
		defer c.Close()

		authed, err := c.Do("AUTH", bt.pass)
		if err != nil {
			return err
		} else {
			logp.Debug("smshubbeat", "AUTH %v", authed)
		}
	}

	return nil
}

// Run the beater's loop
func (bt *Smshubbeat) Run(b *beat.Beat) error {
	var err error

	bt.Setup(b)
	ticker := time.NewTicker(bt.period)
	defer ticker.Stop()

	if len(bt.luascript) > 0 {
		for {
			select {
			case <-bt.done:
				return nil
			case <-ticker.C:
			}

			timerStart := time.Now()

			// Retrieve and publish latest kpi
			err = bt.retrieveLatestKpi()
			if err != nil {
				logp.Err("Error retrieving latest kpi: %v", err)
			}

			timerEnd := time.Now()
			duration := timerEnd.Sub(timerStart)
			if duration.Nanoseconds() > bt.period.Nanoseconds() {
				logp.Warn("Ignoring tick(s) due to processing taking longer than one period")
			}
		}
	}

	return err
}

// Beater is stopping
func (bt *Smshubbeat) Stop() {
	bt.client.Close()
	close(bt.done)
	bt.redisPool.Close()
}

// Retrieve and publish to logstash the latest kpi values
func (bt *Smshubbeat) retrieveLatestKpi() error {

	// Get a connection from pool
	c := bt.redisPool.Get()

	// Defer connection close on exit
	defer c.Close()

	if bt.auth {
		authed, err := c.Do("AUTH", bt.pass)
		if err != nil {
			logp.Err("auth error: %v", bt.pass)
			return err
		} else {
			logp.Debug("smshubbeat", "AUTH %v", authed)
		}
	}

    beforets := time.Now().UnixNano()
	reply, err := c.Do("EVAL", bt.luascript, 0, bt.dbid)
    afterts := time.Now().UnixNano()
    ellapsedmicro := (afterts - beforets) / int64(time.Millisecond)

	if err != nil {
		return err

	} else {
		switch reply := reply.(type) {
		case []interface{}:
			if len(reply) > 0 {
				result := make([][]string, len(reply))
				for i := range reply {
					if reply[i] == nil {
						continue
					}
					strings, err := redis.Strings(reply[i], err)
					if err != nil {
						logp.Err("Failed to convert to Strings array")
					} else {
						result[i] = strings
					}
				}

                now := common.Time(time.Now())
                var event common.MapStr
				var events []common.MapStr
				for _, data := range result {

                    // Expecting: [kpi:cnt:ss7box:group1.key1 0]
                    // Array index: 0 = key, 1 = value
					if len(data) >= 2 {
                    
	                    // Expected key format: 
                        // keyfields[0] = kpi tag
                        // keyfields[1] = data typ (int, cnt,str)
                        // keyfields[2] = box name (smppbox, ss7box, httpbox, router)
                        // keyfields[3] = kpi name (group.kpi) 
                        keyfields := strings.SplitN(data[0], ":", 4)
                        if len(keyfields) == 4 {
                            var family = "default"
                            kpifields := strings.SplitN(keyfields[3], ".", 2)
                            if len(kpifields) >= 2 {
                                family = kpifields[0]
                            }
                            
                            // data[0] --> key
                            // data[1] --> value
                            // keyfields[1] --> data type
                            // keyfields[2] --> box name
                            // keyfields[3] --> kpi name
                            // family --> kpi family
                            var buffer bytes.Buffer
                            buffer.WriteString(keyfields[2])
                            buffer.WriteByte('.')
                            buffer.WriteString(keyfields[3])
                            kpiname := buffer.String()

                            if event != nil {
                                if event["metricset.module"] != keyfields[2] || 
                                    event["metricset.name"] != family {
                                    events = append(events, event) 
                                    event = nil
                                } 
                            }

                            if event == nil {
                                var metricset map[string]string
                                metricset = make(map[string]string)
                                metricset["module"] = keyfields[2] 
                                metricset["name"] = family
                                metricset["rtt"] = strconv.FormatInt(ellapsedmicro, 10)
   
                                event = common.MapStr{
								    "@timestamp": now,
                                    "@version": "1.0",
                                    "host": bt.host,
                                    "type": "metricsets",
                                    "metricset": metricset,
                                } 
                            }

                            // TODO - Need to structure the  key.key.key  below   as { "module": { "name": { "kpi": "value" } } }  
                            
                            if  keyfields[1] == "str" {
                                event[kpiname] = data[1]
                            } else {
                                event[kpiname], err = strconv.Atoi(data[1])
					            if err != nil {
                                    event[kpiname] = 0
                                }
                            }
                        }
					}
				}

                // Add the last event generated
                if event != nil {
                    events = append(events, event) 
                    event = nil
                }

				if len(events) > 0 {
					// fmt.Println(events)
					// logp.Debug("smshubbeat", "Would have %v events to publish", len(events))
			        bt.client.PublishEvents(events)
				}
			}
		}
		return nil
	}
}

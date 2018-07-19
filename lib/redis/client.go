package redis

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/garyburd/redigo/redis"
)

// Redigo
func newPool() *redis.Pool {
	timeout  := redis.DialConnectTimeout(2 * time.Second)
	password := redis.DialPassword(beego.AppConfig.String("redis_pass"))
	tcp 	 := beego.AppConfig.String("redis_tcp")
	port 	 := beego.AppConfig.String("redis_port")

	return &redis.Pool{
		MaxIdle		: 500,
		MaxActive	: 20000, // max number of connections
		IdleTimeout	: 5 * time.Second,
		Wait		: true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(tcp, port, timeout, password)
			if err != nil {
				beego.Error(err.Error())
			}
			return c, err
		},
	}
}
var pool = newPool()

// Set Redis with Timeout
func SetEx(key interface{}, val interface{}, timeout time.Duration) (err error) {
	c := pool.Get()

	if beego.BConfig.RunMode == "dev" {
		beego.Debug("[REDIS] SET")
		beego.Debug(key)
		beego.Debug(val)
	}

	switch c.Err() {
		case nil:
			_, err = c.Do("SETEX", key, int64(timeout/time.Second), val)
			c.Close()
			return err
		default:
			c.Close()
			return c.Err()
	}
}

// Set redis no expiry
func Set(key interface{}, val interface{}) (err error) {
	c := pool.Get()

	if beego.BConfig.RunMode == "dev" {
		beego.Debug("[REDIS] SET NO TIMEOUT")
		beego.Debug(key)
		beego.Debug(val)
	}

	switch c.Err() {
	case nil:
		_, err = c.Do("SET", key, val)
		c.Close()
		return err
	default:
		c.Close()
		return c.Err()
	}
}

// Del Redis
func Del(val ...interface{}) (err error) {
	c := pool.Get()

	if beego.BConfig.RunMode == "dev" {
		beego.Debug("[REDIS] DELETE")
		beego.Debug(val)
	}

	switch c.Err() {
		case nil:
			_, err = c.Do("DEL", val...)
			c.Close()
			return err
		default:
			c.Close()
			return c.Err()
	}
}

// Get Result as a string
func GetString(command string, key ...interface{}) (rs string, err error) {
	c := pool.Get()

	if beego.BConfig.RunMode == "dev" {
		beego.Debug("[REDIS] GETSTRING "+ command)
		beego.Debug(key)
	}

	switch c.Err() {
		case nil:
			rs, err = redis.String(c.Do(command, key...))
			c.Close()
			return rs, err
		default:
			c.Close()
			return rs, c.Err()
	}
}
// Get Result as a []string
func GetStrings(command string, key ...interface{}) (rs []string, err error) {
	c := pool.Get()

	if beego.BConfig.RunMode == "dev" {
		beego.Debug("[REDIS] GETSTRINGS "+ command)
		beego.Debug(key)
	}

	switch c.Err() {
		case nil:
			rs, err = redis.Strings(c.Do(command, key...))
			c.Close()
			return rs, err
		default:
			c.Close()
			return rs, c.Err()
	}
}

// Get Result as an int
func GetInt(command string, key ...interface{}) (rs int, err error) {
	c := pool.Get()

	if beego.BConfig.RunMode == "dev" {
		beego.Debug("[REDIS] GETINT " + command)
		beego.Debug(key)
	}

	switch c.Err() {
		case nil:
			rs, err = redis.Int(c.Do(command, key...))
			c.Close()
			return rs, err
		default:
			c.Close()
			return rs, c.Err()
	}
}
// Get Result as an []int
func GetInts(command string, key ...interface{}) (rs []int, err error) {
	c := pool.Get()
	switch c.Err() {
	case nil:
		rs, err = redis.Ints(c.Do(command, key...))
		c.Close()
		return rs, err
	default:
		c.Close()
		return rs, c.Err()
	}
}

// Get Result as an int64
func GetInt64(command string, key ...interface{}) (rs int64, err error) {
	c := pool.Get()
	switch c.Err() {
	case nil:
		rs, err = redis.Int64(c.Do(command, key...))
		c.Close()
		return rs, err
	default:
		c.Close()
		return rs, c.Err()
	}
}
// Get Result as an []int64
func GetInt64s(command string, key ...interface{}) (rs []int64, err error) {
	c := pool.Get()
	switch c.Err() {
	case nil:
		rs, err = redis.Int64s(c.Do(command, key...))
		c.Close()
		return rs, err
	default:
		c.Close()
		return rs, c.Err()
	}
}

// Get Result as a float64
func GetFloat64(command string, key ...interface{}) (rs float64, err error) {
	c := pool.Get()
	switch c.Err() {
	case nil:
		rs, err = redis.Float64(c.Do(command, key...))
		c.Close()
		return rs, err
	default:
		c.Close()
		return rs, c.Err()
	}
}
// Get Result as a []float64
func GetFloat64s(command string, key ...interface{}) (rs []float64, err error) {
	c := pool.Get()
	switch c.Err() {
	case nil:
		rs, err = redis.Float64s(c.Do(command, key...))
		c.Close()
		return rs, err
	default:
		c.Close()
		return rs, c.Err()
	}
}

// Get Result as a []byte
func GetBytes(command string, key ...interface{}) (rs []byte, err error) {
	c := pool.Get()

	if beego.BConfig.RunMode == "dev" {
		beego.Debug("[REDIS] GETBYTES "+ command)
		beego.Debug(key)
	}

	switch c.Err() {
		case nil:
			rs, err = redis.Bytes(c.Do(command, key...))
			c.Close()
			return rs, err
		default:
			c.Close()
			return rs, c.Err()
	}
}

package distlock

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/redigo"
	redigolib "github.com/gomodule/redigo/redis"
	"log"
	"os"
	"time"
)

func GetLock(key string) (*redsync.Mutex, error) {
	certBase64 := os.Getenv("REDIS_CERT_BASE64")
	connectionUrl := os.Getenv("REDIS_CONNECTION_URL")

	cert, err := base64.StdEncoding.DecodeString(certBase64)
	if err != nil {
		log.Fatalln("could not encode cert")
	}

	config := &tls.Config{InsecureSkipVerify: true}
	config.RootCAs = x509.NewCertPool()
	config.ClientAuth = tls.RequireAndVerifyClientCert
	if ok := config.RootCAs.AppendCertsFromPEM(cert); !ok {
		log.Fatalln("failed to append cert")
		return nil, errors.New("failed to append cert")
	}

	pool := redigo.NewPool(&redigolib.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redigolib.Conn, error) {
			return redigolib.DialURL(connectionUrl, redigolib.DialTLSConfig(config))
		},
		TestOnBorrow: func(c redigolib.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	})

	rs := redsync.New(pool)
	mutex := rs.NewMutex(key)
	return mutex, nil
}

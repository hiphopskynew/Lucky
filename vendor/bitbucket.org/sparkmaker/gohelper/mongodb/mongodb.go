package mongodb

import (
	"errors"
	"regexp"
	"strings"
	"time"

	mgo "github.com/globalsign/mgo"
)

type metadata struct {
	url      string
	database string
	option   Option
}

type Mongo struct {
	Database  *mgo.Database
	session   *mgo.Session
	meta      metadata
	closeTime int64
}

// Option is a configuration of MongoDB connection
// - ConnectionLimitPerHost is a maximum of connections has a default value are 100 connections
// - Timeout is timeout as the amount of time to wait for a server to respond
// - MaxLifeTime will close the connection if it exceeds the MaxLifeTime. Since new a connection (0 is forever)
// - MaxIdleTime will close the connection if idle since the last sync (0 is forever)
type Option struct {
	ConnectionLimitPerHost int           //Default 100 connections
	Timeout                time.Duration //Default 60 seconds
	MaxLifeTime            time.Duration //Default 60 seconds
	Username               string
	Password               string
}

var flock []*Mongo = []*Mongo{}

func (opt *Option) autoValue() {
	opt.ConnectionLimitPerHost = setIntOrElse(opt.ConnectionLimitPerHost, 100)
	opt.Timeout = setTimeOrElse(opt.Timeout, time.Second*60)
	opt.MaxLifeTime = setTimeOrElse(opt.MaxLifeTime, time.Second*60)
}

func (mdb *Mongo) life() {
	mdb.setCloseTime()
	for now := range time.Tick(time.Second) {
		if now.UnixNano() > mdb.closeTime {
			mdb.Close()
			break
		}
	}
}

// Get using current database name
func (mdb *Mongo) GetDatabase() string {
	return mdb.meta.database
}

// Change pointer to new database name
func (mdb *Mongo) Use(database string) *Mongo {
	return NewWithOption(mdb.meta.url, database, mdb.meta.option)
}

// Clone configuration from established connection and new connection
func (mdb *Mongo) Clone() *Mongo {
	return NewWithOption(mdb.meta.url, mdb.meta.database, mdb.meta.option)
}

// Close established current connection
func (mdb *Mongo) Close() {
	if mdb.meta.option.MaxLifeTime == 0 {
		return
	}
	mdb.session.Close()
	mdb.leaveFlock()
}

func (mdb *Mongo) joinFlock() {
	flock = append(flock, mdb)
}

func (mdb *Mongo) leaveFlock() {
	for i, db := range flock {
		if mdb == db {
			flock = append(flock[:i], flock[i+1:]...)
			break
		}
	}
}

func (mdb *Mongo) setCloseTime() {
	mdb.closeTime = time.Now().Add(mdb.meta.option.MaxLifeTime).UnixNano()
}

func hasAuth(url string) bool {
	r, _ := regexp.Compile(`^mongodb://.+:.+@.+`)
	return r.MatchString(url)
}

func auth(meta metadata) (string, string) {
	if !hasAuth(meta.url) {
		return meta.option.Username, meta.option.Password
	}

	url := meta.url
	r, _ := regexp.Compile(`^mongodb://`)
	url = r.ReplaceAllString(url, "")
	url = strings.Split(url, "@")[0]
	username := strings.Split(url, ":")[0]
	password := strings.Split(url, ":")[1]

	return username, password
}

func setTimeOrElse(v time.Duration, or time.Duration) time.Duration {
	if v <= 0 {
		return or
	}
	return v
}

func setIntOrElse(v int, or int) int {
	if v <= 0 {
		return or
	}
	return v
}

func useInFlock(meta metadata) (*Mongo, error) {
	if len(flock) >= meta.option.ConnectionLimitPerHost {
		for _, db := range flock {
			if db.meta == meta {
				db.setCloseTime()
				return db, nil
			}
		}
	}

	return &Mongo{}, errors.New("did not match connection in the flock")
}

func login(session *mgo.Session, meta metadata) {
	username, password := auth(meta)
	if len(strings.TrimSpace(username)) == 0 && len(strings.TrimSpace(password)) == 0 {
		return
	}

	if err := session.Login(&mgo.Credential{Username: username, Password: password}); err != nil {
		panic(err)
	}
}

func establishConnection(meta metadata) *Mongo {
	if db, err := useInFlock(meta); err == nil {
		return db
	}
	session, err := mgo.DialWithTimeout(meta.url, meta.option.Timeout)
	if err != nil {
		panic(err)
	}
	login(session, meta)
	established := &Mongo{Database: session.DB(meta.database), session: session, meta: meta}
	established.joinFlock()
	go established.life()
	return established
}

// New a database connection
func New(url string, database string) *Mongo {
	return NewWithOption(url, database, Option{})
}

// NewWithOption for new a database connection with custom configurations
func NewWithOption(url string, database string, opt Option) *Mongo {
	opt.autoValue() // Auto set values default if empty or mistake value
	return establishConnection(metadata{url: url, database: database, option: opt})
}

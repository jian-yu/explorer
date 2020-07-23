package db

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

type mongoStore struct {
	URL  string
	Name string
	*viper.Viper
	*mgo.Session
}

//NewMongoStore new a mongodb session
func NewMongoStore() MgoOperator {
	mgoURL := viper.GetString(`MongoDB.URL`)
	mgoDBName := viper.GetString(`MongoDB.DBName`)
	if mgoURL == "" || mgoDBName == "" {
		log.Error().Str(`params`, "mongo url or name bot be empty").Msg(`NewMongoStore`)
		return nil
	}
	info, _ := mgo.ParseURL(mgoURL)
	info.PoolLimit = viper.GetInt(`MongoDB.PoolLimit`)
	info.Timeout = viper.GetDuration(`MongoDB.Timeout`)
	sess, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}

	if err = sess.Ping(); err != nil {
		panic(err)
	}

	return &mongoStore{
		URL:     mgoURL,
		Name:    mgoDBName,
		Session: sess,
	}
}

func (m *mongoStore) GetSession() *mgo.Session {
	return m.Session.Clone()
}

func (m *mongoStore) GetDBConn() *mgo.Database {
	sess := m.Session.Clone()
	return sess.DB(m.Name)
}

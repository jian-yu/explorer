package db

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)


type mongoStore struct {
	URL  string
	Name string
	*mgo.Session
}

//NewMongoStore new a mongodb session
func NewMongoStore() MgoOperator {
	mgoURL := viper.GetString(`MongoDB.URL`)
	mgoDBName := viper.GetString(`MongoDB.DBName`)
	if mgoURL == "" || mgoDBName == "" {
		log.Panic().Str(`params`, "mongo url or name bot be empty").Msg(`NewMongoStore`)
	}
	log.Debug().Interface(`mgoURL`, mgoURL).Msg(`NewMongoStore`)
	log.Debug().Interface(`mgoDBName`, mgoDBName).Msg(`NewMongoStore`)

	info, err := mgo.ParseURL(mgoURL)
	if err != nil {
		panic(err)
	}

	info.PoolLimit = viper.GetInt(`MongoDB.PoolLimit`)
	info.Timeout = time.Duration(viper.GetInt(`MongoDB.Timeout`)) * time.Second
	//info.Mechanism = viper.GetString(`MongoDB.Mechanism`)
	info.Username = viper.GetString(`MongoDB.Username`)
	info.Password = viper.GetString(`MongoDB.Password`)

	sess, err := mgo.DialWithInfo(info)
	if err != nil {
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

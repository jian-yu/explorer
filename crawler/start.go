package crawler

import "explorer/db"

func OnStart() {
	mgoStore := db.NewMongoStore()
	crawler := New(mgoStore)
	crawler.Run()
}

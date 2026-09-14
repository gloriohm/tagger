package main

import (
	"log"
	"tagger/internal/fetcher"
	"tagger/internal/storage"
	"tagger/internal/tagger"
	"tagger/internal/ui"

	"go.uploadedlobster.com/musicbrainzws2"
)

func main() {
	database, err := storage.Open("tagger.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	repo := tagger.NewSQLiteRepo(database)
	appInfo := musicbrainzws2.AppInfo{Name: "mangum0pus@proton.me", Version: "0.0.1"}
	mbclient := musicbrainzws2.NewClient(appInfo)
	fetchSvc := fetcher.NewBrainzClient(mbclient)
	tagSvc := tagger.NewService(repo, fetchSvc)

	if err = tagSvc.InitializeDirectory("test"); err != nil {
		log.Fatal(err)
	}

	if err := ui.Run(tagSvc); err != nil {
		log.Fatal(err)
	}
}

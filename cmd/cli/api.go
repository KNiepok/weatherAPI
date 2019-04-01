package cli

import (
	"database/sql"
	"github.com/kniepok/weatherAPI/api"
	"github.com/kniepok/weatherAPI/conf"
	"github.com/kniepok/weatherAPI/fetcher/open_weather"
	"github.com/kniepok/weatherAPI/service"
	"github.com/kniepok/weatherAPI/store/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ortuman/jackal/log"
	"github.com/spf13/cobra"
)

var (
	dbPath = "./sqlite.db"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var (
	apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Start API",
		Long:  `Start API`,
		Run: func(cmd *cobra.Command, args []string) {

			db, err := sql.Open("sqlite3", dbPath)

			if err != nil {
				log.Fatal(err)
			}

			weatherService, err := sqlite.NewWeatherStorage(db)

			if err != nil {
				panic(err)
			}

			bookmarksService, err := sqlite.NewBookmarkStorage(db)

			if err != nil {
				panic(err)
			}

			//apiKey := os.Getenv("OWM_API_KEY")
			apiKey := "4c8382b8189f3b8ceb421d363a84c4f2"
			weatherFetcher := open_weather.New(apiKey)

			service := service.New(weatherService, bookmarksService, weatherFetcher)

			s := api.New(service)

			if err != nil {
				log.Fatal(err)
			}
			err = s.ListenAndServe()

			if err != nil {
				log.Fatal(err)
			}

			<-conf.Stop.Chan()
			conf.Stop.Wait()
		},
	}
)

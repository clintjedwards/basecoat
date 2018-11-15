package api

import (
	"fmt"

	"github.com/clintjedwards/basecoat/api/storage"
	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/models"
	"github.com/clintjedwards/basecoat/utils"
	"github.com/go-pg/pg"
)

//API represents an instance of the basecoat restful api service
type API struct {
	config *config.Config
	db     *pg.DB
}

//NewAPI creates a new instance of the basecoat restful api service
func NewAPI(config *config.Config) *API {

	db := storage.NewPostgresDB(config.Database.User, config.Database.Password, config.Database.URL, config.Database.Name)
	err := storage.InitDB(db, []interface{}{&models.Formula{}, &models.Job{}})
	if err != nil {
		utils.StructuredLog("fatal", "cannot connect to database", err)
	}

	utils.StructuredLog("info", fmt.Sprintf("connected to database: %s@%s/%s", config.Database.User, config.Database.URL, config.Database.Name), err)

	return &API{
		config: config,
		db:     db,
	}
}

package feature_tests

import (
	"gcstatus/config"
	testutils "gcstatus/tests/utils"

	"gorm.io/gorm"
)

var (
	dbConn *gorm.DB
	env    *config.Config
)

func init() {
	dbConn, env = testutils.SetupTestDB(nil)
}

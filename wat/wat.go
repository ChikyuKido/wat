package wat

import (
	"github.com/ChikyuKido/wat/wat/server/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitWat(engine *gin.Engine, db *gorm.DB) {
	wat.InitDatabase(db)

}

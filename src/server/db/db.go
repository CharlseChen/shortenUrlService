package db
import (
	"config"
	"fmt"
	"os"

	// 给 GORM 使用绑定mysql数据库
	_ "github.com/go-sql-driver/mysql"

	"github.com/jinzhu/gorm"
)

//// 常量
const (
	RecordNotFound         = "record not found"
	DataOpreationErrorNo   = 9003
	DataOpreationErrorInfo = "数据操作失败"
	LibTypeTerm            = 1
	LibTypeSent            = 2
)

//MyDBKv 数据库映射
type MyDBKv map[string]interface{}

var (
	//DbConn 数据库连接
	DbConn *gorm.DB
)

//Init 初始化
func Init() {
	var err error
	dbConfig,_:= config.NewConfigure()
	connStr := fmt.Sprintf("Connecting to database [%s@%s:%s/%s]...",
		dbConfig.DB_Config.username,
		dbConfig.DB_Config.host,
		dbConfig.DB_Config.port,
		dbConfig.DB_Config.db_name)
	fmt.Println(connStr)

	DbConn, err = connDB(dbConfig.DB_Config.host,
		dbConfig.DB_Config.port,
		dbConfig.DB_Config.db_name,
		dbConfig.DB_Config.username,
		dbConfig.DB_Config.password)

	if err != nil {
		fmt.Println("Connect to database failed! Error=\"", err.Error(), "\"")
		os.Exit(1)
	}
	fmt.Println("初始化数据库成功", err.Error(), "\"")
	// initAppSchema()

}

func connDB(host string, port string, dbname string,
	username string, password string) (*gorm.DB, error) {
	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		username, password, host, port, dbname)

	return gorm.Open("mysql", connStr)
}
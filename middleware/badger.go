package middleware

import (
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const dbKey = "Db"

var db *badger.DB

// ConnectDB 初始化BadgerDB连接并保存到全局变量
func ConnectDB(dataPath string, options ...badger.Options) error {
	// 合并配置选项
	opt := badger.DefaultOptions(dataPath)
	if len(options) > 0 {
		opt = options[0]
	}

	// 确保数据目录存在
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		log.Errorf("创建数据库目录失败: %v", err)
		return err
	}

	// 打开数据库
	dbInstance, err := badger.Open(opt)
	if err != nil {
		log.Errorf("连接数据库失败: %v", err)
		return err
	}

	db = dbInstance
	log.Info("BadgerDB连接成功")
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Errorf("关闭数据库失败: %v", err)
			return err
		}
		log.Info("数据库连接已关闭")
	}
	return nil
}

// DbConnector 中间件：注入数据库连接到请求上下文
func DbConnector() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 将全局db实例克隆到请求上下文
		c.Set(dbKey, db)

		// 请求前操作（可选）
		// 例如开启事务等，根据业务需求实现

		c.Next()

		// 请求后操作（可选）
		// 例如提交事务或清理资源
	}
}

// GetDB 从上下文中获取数据库实例
func GetDB(c *gin.Context) *badger.DB {
	return c.MustGet(dbKey).(*badger.DB)
}

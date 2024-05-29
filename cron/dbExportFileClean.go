package cron

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func DbExportFileClean() {
	// 定时清理导出后数据库后的临时数据
	// 每30分钟从redis中查询key前缀为total_export_*的数据，并判断create_time是否超过2个小时，超过则先删除临时目录数据，再从redis中删除
	keys, err := G.RedisClient.Keys("total_export_*").Result()
	if err != nil {
		G.Logger.Errorf("获取key失败:[%s]", err.Error())
		return
	}
	for _, key := range keys {
		go func(key string) {
			resultMap, err := G.RedisClient.HGetAll(key).Result()
			if err != nil {
				G.Logger.Errorf("获取key失败:[%s]", err.Error())
				return
			}
			createTimeInt64, err := strconv.ParseInt(resultMap["create_time"], 0, 0)
			if err != nil {
				G.Logger.Errorf("createTimeInt64失败:[%s]", err.Error())
				return
			}

			diff := int64(1000 * 60 * 4)
			currTimeStamp := time.Now().UnixMilli()
			if (currTimeStamp - createTimeInt64) > diff {
				G.Logger.Infof("[定时任务]清理导出后数据库后的临时数据，key=[%s]", key)
				splits := strings.Split(resultMap["save_path"], "/")
				rmPath := strings.Join(splits[:2], "/")
				// fmt.Println("key:", key, resultMap["save_path"], rmPath)
				err = os.RemoveAll(rmPath)
				if err != nil {
					G.Logger.Errorf("移除目录失败:[%s]", err.Error())
					return
				}
				G.RedisClient.Del(key)
			}
		}(key)
	}
}

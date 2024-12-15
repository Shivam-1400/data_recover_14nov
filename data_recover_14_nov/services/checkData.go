package services

import (
	"context"
	"data_recover_14_nov/databases"
	"data_recover_14_nov/globals"
	"data_recover_14_nov/model"
	"fmt"
	"sync"
	"time"
)

func CheckData(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Data checking stopped due to shutdown signal.")
			return
		default:
			if globals.DataMap != nil {
				globals.DataMap.Range(func(key, value interface{}) bool {
					logData, ok := value.(model.BulksmsLogData)
					if ok {
						fmt.Printf("Processing Log Data\n")
						globals.DataMap.Delete(key)
						checkRedisLvl1(logData)
					} else {
						fmt.Println("Failed to assert value as BulksmsLogData")
					}
					return true
				})
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func checkRedisLvl1(data model.BulksmsLogData) {
	key := globals.ApplicationConfig.RedisQueue.QueueLvl1 + "-" + data.BReqID

	exists, _ := databases.RedisClient.Exists(context.Background(), key).Result()

	if exists > 0 {
		val, _ := databases.RedisClient.Get(context.Background(), key).Result()
		fmt.Printf("found breqid :: %s\n", key)
		checkRedisLvl3(val, data)
	} else {
		fmt.Printf("breqid %s not found\n", key)
	}
}

func checkRedisLvl3(deliveryid string, data model.BulksmsLogData) {
	key := globals.ApplicationConfig.RedisQueue.QueueLvl3 + "-" + deliveryid + "-" + data.To
	exists, _ := databases.RedisClient.Exists(context.Background(), key).Result()

	if exists == 1 {
		fmt.Printf("found redisid :: %v\n", key)
		putLvl3(key, data)
	} else {
		fmt.Printf("redisid %s not found\n", key)
	}
}

func putLvl3(redisid string, data model.BulksmsLogData) {
	result := make(map[string]string)

	result["appid"] = data.AppID
	result["tid"] = data.TID
	result["feedid"] = data.FeedID
	result["entid"] = data.EntID
	result["keyword"] = data.Keyword
	result["orgTempalte"] = data.OrgTemplate
	result["dotStarCnt"] = data.DotStarCnt
	result["spaceFlag"] = data.SpaceFlag
	result["specialCharFlag"] = data.SpecialCharFlag
	result["customdomain"] = data.CustomDomain
	result["token"] = data.Token
	result["dmCheckStatus"] = data.DMCheckStatus
	result["breqid"] = data.BReqID
	result["btid"] = data.BTID
	result["traicategoryid"] = data.TraiCategoryID
	result["traimessagetype"] = data.TraiMessageType
	result["traimessagemode"] = data.TraiMessageMode
	result["bsms_intime"] = data.BSMSInTime
	result["template_id"] = data.TemplateID
	result["bmsgtag"] = data.BMsgTag
	result["text"] = data.Text
	result["to"] = data.To
	result["from"] = data.From
	result["dltentityid"] = data.DLTEntityID
	result["bsms_outtime"] = data.BSMSOutTime

	var fields []interface{}
	for key, value := range result {
		fields = append(fields, key, value)
	}

	fmt.Printf("Putting data into hash redisid :: %s\n", redisid)
	_, err := databases.RedisClient.HSet(context.Background(), redisid, fields).Result()
	if err != nil {
		fmt.Printf("Error inserting data into redisd:: %s.\t %v\n", redisid, err)
	}

	fmt.Printf("Data inserted into hash successfully redisis: %v\n", redisid)
	_, err1 := databases.RedisClient.LPush(context.Background(), "crate_db_todo_list", redisid).Result()
	if err1 != nil {
		fmt.Printf("Error LPUSH to crate_db_todo_list. redisid: %s", redisid)
	}
	fmt.Printf("Successful LPUSH to crate_db_todo_list. redisid: %s", redisid)
}

package mixflow

import (
	"fmt"

	"sync"

	"weibo.com/hotweibo/model/mixflow"
	"weibo.com/hotweibo/service"
)

func Test() {
	gr_num := 3
	total := 0
	uve_ids := []string{}
	rec_ids := []string{}
	loc_info := map[string]string{}
	data_chans := make(chan map[string]interface{}, gr_num)
	defer close(data_chans)
	go mixflow.FetchUve(data_chans)
	go mixflow.FetchRec(data_chans)
	go mixflow.FetchLocation(data_chans)
	var waitgroup sync.WaitGroup
	waitgroup.Add(gr_num)
	for inter_data := range data_chans {
		switch inter_data["type"].(string) {
		case "uve":
			uve_ids = inter_data["data"].([]string)
		case "rec":
			rec_ids = inter_data["data"].([]string)
		case "loc":
			loc_info = inter_data["data"].(map[string]string)

		}
		waitgroup.Done()
		total++
		if total == gr_num {
			break
		}
	}
	waitgroup.Wait()
	fmt.Println(uve_ids, rec_ids, loc_info)
	mixflow.FetchUserTagNScale()
	return
	//MC

	//var element = &common.Element{
	//	Key:     "test-golang",
	//	Flags:   0,
	//	Exptime: 3600,
	//	Value:   []byte("test mc pool"),
	//}
	//err := basic.Mcpool.Add(element)
	//if err != nil {
	//	fmt.Println(err)
	//}
	item, err := basic.Mcpool.Get("test-golang")
	if err != nil {
		fmt.Println(err)
		return
	}
	key := item.Key()
	value := item.Value()
	//flags := item.Flags()
	//cas := item.Cas()

	fmt.Println(key, string(value))

	//DB
	rows, err := basic.Dbpool.Query("SELECT * FROM mlog limit 1")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]string)
	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
	}
	fmt.Println(record)

	//redis
	conn := basic.Redispool.Get()
	defer conn.Close()
	n, err := conn.Do("get", "test-redigo")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(n.([]byte)))
	}
}

package mixflow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/ningjh/memcached/common"

	"weibo.com/hotweibo/service"
)

func FetchUve(data_chans chan<- map[string]interface{}) {
	res := map[string]interface{}{"type": "uve", "data": []string{}}
	req, err := http.NewRequest("POST", "http://api.uve.mobile.sina.cn/uve/service/hot_tweets_feed", strings.NewReader("uid=2908068201&from=1056095010&ua=Xiaomi-MI+4LTE__weibo__5.6.0__android__android4.4.4&wm=20005_0002&ip=10.209.73.87&lang=zh_CN&source=99075054&gsid=_2A254dT0EDeTxGeRL6FAU8SbJyTuIHXVZIzfMrDV6PUJbrdANLUbHkWpv0dtVA5JVzfydSKEu5PBeuRzAqQ..&unread_status=20"))
	if err != nil {
		data_chans <- res
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := basic.Httppool.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		data_chans <- res
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		data_chans <- res
		return
	}
	_, ok := js.CheckGet("ad")
	if ok {
		ads, err := js.Get("ad").Array()
		if err != nil {
			data_chans <- res
			return
		}
		for _, v := range ads {
			line, _ := v.(map[string]interface{})
			mid, exists := line["id"]
			if exists {
				res["data"] = append(res["data"].([]string), mid.(string))
			}
		}
	}
	data_chans <- res
}

func FetchRec(data_chans chan<- map[string]interface{}) {
	res := map[string]interface{}{"type": "rec", "data": []string{}}
	item, err := basic.HotMcpool.Get("mixed_feeds_rec_test")
	if err == nil {
		js, err := simplejson.NewJson(item.Value())
		if err == nil {
			rec_data, err := js.Array()
			if err == nil {
				for _, v := range rec_data {
					line, _ := v.(map[string]interface{})
					mid, exists := line["mid"]
					if exists {
						res["data"] = append(res["data"].([]string), mid.(string))
					}
				}
				data_chans <- res
				return
			}
		}
	}

	rows, err := basic.BasicDbpool.Query("SELECT  `mid`, `page`, `end_time`, `start_time`, `status`  FROM hotmblog_mixed_feeds_rec WHERE  `end_time` >= 1450337684 AND `page` <= 10 AND `status` = 1 ORDER BY  page asc LIMIT  0, 100")
	defer rows.Close()
	if err != nil {
		data_chans <- res
		return
	}
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	res_data := []map[string]string{}
	record := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		res_data = append(res_data, record)
	}
	rec_str, err := json.Marshal(res_data)
	if err != nil {
		data_chans <- res
		return
	}
	var element = &common.Element{
		Key:     "mixed_feeds_rec_test",
		Flags:   0,
		Exptime: 3600,
		Value:   rec_str,
	}
	err = basic.HotMcpool.Add(element)
	for _, v := range res_data {
		mid, exists := v["mid"]
		if exists {
			res["data"] = append(res["data"].([]string), mid)
		}
	}
	data_chans <- res
	return
}

func FetchLocation(data_chans chan<- map[string]interface{}) {
	res := map[string]interface{}{"type": "loc", "data": map[string]string{}}
	item, err := basic.HotMcpool.Get("user_location_info_test_2908068201")
	if err == nil {
		js, err := simplejson.NewJson(item.Value())
		if err == nil {
			loc_data, err := js.Map()
			loc_info := map[string]string{"province_name": loc_data["province_name"].(string), "city_name": loc_data["city_name"].(string)}
			if err == nil {
				res["data"] = loc_info
				data_chans <- res
				return
			}
		}
	}

	req, err := http.NewRequest("GET", "http://i2.api.weibo.com/2/darwin/table/show.json?appkey=1428722706&table=user_city&key=1042:2908068201&source=99075054", nil)
	if err != nil {
		data_chans <- res
		return
	}
	resp, err := basic.Httppool.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		data_chans <- res
		return
	}
	js, err := simplejson.NewJson(body)
	if err != nil {
		data_chans <- res
		return
	}
	if _, ok := js.CheckGet("columns"); !ok {
		data_chans <- res
		return
	}
	city, city_ok := js.Get("columns").CheckGet("city_dm")
	province, province_ok := js.Get("columns").CheckGet("province_dm")
	if !city_ok || !province_ok {
		data_chans <- res
		return
	}
	province_data, err := province.String()
	if err != nil {
		data_chans <- res
		return
	}
	city_data, err := city.String()
	if err != nil {
		data_chans <- res
		return
	}
	loc_info := map[string]string{"province_name": province_data, "city_name": city_data}
	loc_str, err := json.Marshal(loc_info)
	if err != nil {
		data_chans <- res
		return
	}
	var element = &common.Element{
		Key:     "user_location_info_test_2908068201",
		Flags:   0,
		Exptime: 3600,
		Value:   loc_str,
	}
	err = basic.HotMcpool.Add(element)
	res["data"] = loc_info

	data_chans <- res
}

func FetchUserTagNScale() {
	conn := basic.Redis7474pool.Get()
	defer conn.Close()
	data, err := conn.Do("HMGET", "2908068201", "u_tag_li", "a_bd_score_li")
	if err != nil {
		fmt.Println(err)
		return
	}
	new_data := data.([]interface{})
	var tag_data string = ""
	var scale_data string = ""
	if new_data[0] != nil {
		tag_data = string(new_data[0].([]byte))
	}
	if new_data[1] != nil {
		scale_data = string(new_data[1].([]byte))
	}
	//res := map[string]interface{}{"scale": []float32{}, "tag": map[string]float32{}}
	if scale_data != "" {
		for v := range strings.Split(scale_data, ",") {
			fmt.Println(v)
			//res["scale"] = append(res["scale"], ParseFloat(v))
		}
	}
	if tag_data != "" {
		for v := range strings.Split(tag_data, ",") {
			fmt.Println(v)
			//lv := strings.Split(v, ";")
			//fmt.Println(lv)
			//res["tag"] = append(res["tag"], lv[0])
		}
	}
}

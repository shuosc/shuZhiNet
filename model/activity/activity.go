package activity

import (
	"encoding/json"
	"shuZhiNet/infrastructure"
	"time"
)

type Type struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var Types = []Type{
	{"1", "学业辅导"},
	{"2", "师生互动"},
	{"3", "志愿服务"},
	{"4", "社会实践"},
	{"5", "创新创业"},
	{"6", "文体活动"},
	{"7", "素质拓展"},
	{"8", "就业实习"},
}

func GetTypeByName(name string) Type {
	for _, typeObject := range Types {
		if typeObject.Name == name {
			return typeObject
		}
	}
	return Type{}
}

type Activity struct {
	Id         string    `json:"id"`
	TypeId     string    `json:"type_id"`
	Title      string    `json:"title"`
	Leader     string    `json:"leader"`
	Address    string    `json:"address"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	SignUpTime time.Time `json:"sign_up_time"`
}

func marshal(activity Activity) []byte {
	marshaled, _ := json.Marshal(activity)
	return marshaled
}

func unmarshal(binaryData []byte) Activity {
	result := Activity{}
	_ = json.Unmarshal(binaryData, &result)
	return result
}

func Save(activity Activity) {
	infrastructure.Redis.Set("Activity_"+activity.Id, marshal(activity), 0)
}

func Get(id string) (Activity, error) {
	binaryData, err := infrastructure.Redis.Get("Activity_" + id).Result()
	return unmarshal([]byte(binaryData)), err
}

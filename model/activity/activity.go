package activity

import (
	"encoding/json"
	"shuZhiNet/infrastructure"
	"time"
)

type Activity struct {
	TypeId     string    `json:"type_id"`
	Id         string    `json:"id"`
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
	json.Unmarshal(binaryData, &result)
	return result
}

func Save(activity Activity) {
	infrastructure.Redis.Set("Activity_"+activity.Id, marshal(activity), 0)
}

func Get(id string) (Activity, error) {
	binaryData, err := infrastructure.Redis.Get("Activity_" + id).Result()
	return unmarshal([]byte(binaryData)), err
}

package activity

import (
	"log"
	"shuZhiNet/model/student"
)

func TakePartInActivity(student student.Student, activityId string, phoneNumber string, mailAddress string) {
	client := student.Client()
	getURL := "http://www.sz.shu.edu.cn/api/HuoDong/HuoDBMXX/GetHuoDBM?hdid=" + activityId +
		"&shouJhm=" + phoneNumber + "&email=" + mailAddress
	_, err := client.Get(getURL)
	if err != nil {
		log.Println(err)
	}
}

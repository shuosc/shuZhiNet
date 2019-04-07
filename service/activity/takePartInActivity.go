package activity

import (
	"shuZhiNet/model/student"
)

func TakePartInActivity(student student.Student, activityId string, phoneNumber string, mailAddress string) {
	client := student.Client()
	getURl := "http://www.sz.shu.edu.cn/api/HuoDong/HuoDBMXX/GetHuoDBM?hdid=" + activityId +
		"&shouJhm=" + phoneNumber + "&email=" + mailAddress
	_, _ = client.Get(getURl)
}

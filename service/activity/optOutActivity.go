package activity

import (
	"shuZhiNet/model/student"
)

func OptOutActivity(student student.Student, id string) {
	client := student.Client()
	getURL := "http://www.sz.shu.edu.cn/api/HuoDong/HuoDBMXX/DeleteHuoDBM?hdbmid=" + id
	_, _ = client.Get(getURL)
}

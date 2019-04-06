# API reference
## 登陆
URL: $BACKEND_ADDRESS/login

method: POST

body: 
```json
{"username":"学生证号","password":"学生密码"}
```

response: 
```json
{"student_name": "学生姓名","token": "$JWT_TOKEN"}
```
## 获取所有活动
URL: $BACKEND_ADDRESS/activities

method: GET

response:
```json
[
  {
    "type_id": "活动类型",
    "id": "活动id",
    "title": "活动标题",
    "leader": "学院",
    "address": "活动地点", 
    "start_time": "开始时间",
    "end_time": "结束时间",
    "sign_up_time": "开始报名时间"
  },
  ...
]
```
## 参与活动

URL: $BACKEND_ADDRESS/engage

method: POST

head: Authorization: Bearer $JWT_TOKEN

body: 
```json
{"activity_id": "活动id","phone_number": "","mail_address": ""}
```

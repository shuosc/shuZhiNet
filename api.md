# API reference
## 登录
### Request
- method: `POST`
- URL: `/login`
- Format: `json`
- Body: 
```json
{
  "username":"学生证号",
  "password":"学生密码"
}
```

### Response:
```json
{
  "student_name": "学生姓名",
  "token": "$JWT_TOKEN"
}
```

## 所有活动
### Request
- method: `GET`
- URL: `/activities`

### Response
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
### Request
- method: `POST`
- URL: `/take-part`
- Authorization: Bearer Token

body: 
```json
{
  "activity_id": "活动id",
  "phone_number": "",
  "mail_address": ""
}
```

## 退出活动
### Request
- method: `POST`
- URL: `/opt-out`
- Authorization: Bearer Token

body: 
```json
{
  "activity_id": "活动id"
}
```
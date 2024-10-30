# check_in
無心言看一下

### 資料庫中有兩個table
1.使用者資料
2.打卡紀錄

### 註冊(POST) http://localhost:8080/register
BODY:真實姓名、使用者姓名、密碼、email
```
{
  "name": "葉致嘉",
  "username": "jason",
  "password": "qwe123",
  "email": "chichia.yeh@gmail.com"
}

username重複會Error:
{
  "error": "Username already exists",
  "error_code": 1003
}
```

### 登入(POST) http://localhost:8080/login
BODY:使用者姓名、密碼
按下登入後，查詢使用者資料的table，確認身分後，跳轉到http://localhost:8080/check_in?=username(上次上課教的 前端處理)

### 打卡(POST) http://localhost:8080/check_in/username
按下打卡按鈕後，會去抓網址後面的username

```
POST http://localhost:8080/check_in/jason

{
  "check_in_id": 6,
  "check_in_time": "2024-10-30T16:55:57+08:00",
  "message": "Check-in successful",
  "username": "jason"
}
```

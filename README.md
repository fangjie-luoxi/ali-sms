## 阿里云短信服务
使用示例
```go
client := ali_sms.NewSmsClient("accessKeyId", "accessKeySecret")
res, err := client.SendSms("阿里云短信测试", "SMS_154950909", `{"code":"123456"}`, "155xxxx6770")
```
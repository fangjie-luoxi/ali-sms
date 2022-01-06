## 阿里云短信服务
使用示例
```go
client := ali_sms.NewSmsClient("accessKeyId", "accessKeySecret")
param := map[string]string{
"PhoneNumbers":  "155xxxx6770",
"SignName":      "阿里云短信测试",
"TemplateCode":  "SMS_154950909",
"TemplateParam": `{"code":"123456"}`,
}
res, err := client.SendSms(param)
```
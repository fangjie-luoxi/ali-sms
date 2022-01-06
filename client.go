package ali_sms

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type SmsClient struct {
	Request *Request
	Client  *http.Client

	gatewayUrl      string
	accessKeyId     string
	accessKeySecret string
}

func NewSmsClient(accessKeyId, accessKeySecret string) *SmsClient {
	smsClient := new(SmsClient)
	smsClient.Request = &Request{}
	smsClient.Client = &http.Client{}
	smsClient.gatewayUrl = "https://dysmsapi.aliyuncs.com/"
	smsClient.accessKeyId = accessKeyId
	smsClient.accessKeySecret = accessKeySecret
	return smsClient
}

// SendSms PhoneNumbers:电话号码,SignName:短信签名名称,TemplateCode:短信模板ID,TemplateParam:短信模板变量对应的实际值
func (smsClient *SmsClient) SendSms(param map[string]string) (bool, error) {
	phoneNumbers := param["PhoneNumbers"]
	signName := param["SignName"]
	templateCode := param["TemplateCode"]
	templateParam := param["TemplateParam"]
	if phoneNumbers == "" || signName == "" || templateCode == "" || templateParam == "" {
		return false, errors.New("参数不足")
	}

	err := smsClient.Request.SetParamsValue(smsClient.accessKeyId, phoneNumbers, signName, templateCode, templateParam)
	if err != nil {
		return false, err
	}
	endpoint := smsClient.Request.BuildSmsRequestEndpoint(smsClient.accessKeySecret, smsClient.gatewayUrl)
	request, _ := http.NewRequest("GET", endpoint, nil)
	response, err := smsClient.Client.Do(request)
	if err != nil {
		return false, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	result := &Response{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return false, err
	}
	if result.Code != "OK" {
		return false, errors.New(result.Message)
	}

	return true, err
}

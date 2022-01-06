package ali_sms

import (
	"encoding/json"
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

func (smsClient *SmsClient) SendSms(signName, templateCode, templateParam, phoneNumbers string) (*Response, error) {
	err := smsClient.Request.SetParamsValue(smsClient.accessKeyId, phoneNumbers, signName, templateCode, templateParam)
	if err != nil {
		return nil, err
	}
	endpoint := smsClient.Request.BuildSmsRequestEndpoint(smsClient.accessKeySecret, smsClient.gatewayUrl)
	request, _ := http.NewRequest("GET", endpoint, nil)
	response, err := smsClient.Client.Do(request)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	result := &Response{}
	err = json.Unmarshal(body, result)
	return result, err
}

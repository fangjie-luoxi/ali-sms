package ali_sms

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"math/rand"
	"net/url"
	"sort"
	"strings"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var encoding = base32.NewEncoding("ybndrfg8ejkmcpqxot1uwisza345h897")
var src = rand.NewSource(time.Now().UnixNano())

type Request struct {
	//system parameters
	AccessKeyId      string
	Timestamp        string
	Format           string
	SignatureMethod  string
	SignatureVersion string
	SignatureNonce   string
	Signature        string

	//business parameters
	Action          string
	Version         string
	RegionId        string
	PhoneNumbers    string
	SignName        string
	TemplateCode    string
	TemplateParam   string
	SmsUpExtendCode string
	OutId           string
}

// Response @see https://help.aliyun.com/document_detail/55284.html#出参列表
type Response struct {
	RequestId string `json:"RequestId"`
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	BizId     string `json:"BizId"`
}

func (req *Request) SetParamsValue(accessKeyId, phoneNumbers, signName, templateCode, templateParam string) error {
	req.AccessKeyId = accessKeyId
	now := time.Now()
	local, err := time.LoadLocation("GMT")
	if err != nil {
		return err
	}
	req.Timestamp = now.In(local).Format("2006-01-02T15:04:05Z")
	req.Format = "json"
	req.SignatureMethod = "HMAC-SHA1"
	req.SignatureVersion = "1.0"
	req.SignatureNonce = encoding.EncodeToString(randString(27))
	req.Action = "SendSms"
	req.Version = "2017-05-25"
	req.RegionId = "cn-hangzhou"
	req.PhoneNumbers = phoneNumbers
	req.SignName = signName
	req.TemplateCode = templateCode
	req.TemplateParam = templateParam
	req.SmsUpExtendCode = "90999"
	req.OutId = "abcdefg"
	return nil
}

func (req *Request) BuildSmsRequestEndpoint(accessKeySecret, gatewayUrl string) string {
	// common params
	systemParams := make(map[string]string)
	systemParams["SignatureMethod"] = req.SignatureMethod
	systemParams["SignatureNonce"] = req.SignatureNonce
	systemParams["AccessKeyId"] = req.AccessKeyId
	systemParams["SignatureVersion"] = req.SignatureVersion
	systemParams["Timestamp"] = req.Timestamp
	systemParams["Format"] = req.Format

	// business params
	businessParams := make(map[string]string)
	businessParams["Action"] = req.Action
	businessParams["Version"] = req.Version
	businessParams["RegionId"] = req.RegionId
	businessParams["PhoneNumbers"] = req.PhoneNumbers
	businessParams["SignName"] = req.SignName
	businessParams["TemplateParam"] = req.TemplateParam
	businessParams["TemplateCode"] = req.TemplateCode
	businessParams["SmsUpExtendCode"] = req.SmsUpExtendCode
	businessParams["OutId"] = req.OutId
	sortQueryString, signature := generateQueryStringAndSignature(businessParams, systemParams, accessKeySecret)
	return gatewayUrl + "?Signature=" + signature + sortQueryString
}

func generateQueryStringAndSignature(businessParams map[string]string, systemParams map[string]string, accessKeySecret string) (string, string) {
	keys := make([]string, 0)
	allParams := make(map[string]string)
	for key, value := range businessParams {
		keys = append(keys, key)
		allParams[key] = value
	}

	for key, value := range systemParams {
		keys = append(keys, key)
		allParams[key] = value
	}

	sort.Strings(keys)

	sortQueryStringTmp := ""
	for _, key := range keys {
		rstkey := specialUrlEncode(key)
		rstval := specialUrlEncode(allParams[key])
		sortQueryStringTmp = sortQueryStringTmp + "&" + rstkey + "=" + rstval
	}

	sortQueryString := strings.Replace(sortQueryStringTmp, "&", "", 1)
	stringToSign := "GET" + "&" + specialUrlEncode("/") + "&" + specialUrlEncode(sortQueryString)

	sign := sign(accessKeySecret+"&", stringToSign)
	signature := specialUrlEncode(sign)
	return sortQueryStringTmp, signature
}

func specialUrlEncode(value string) string {
	rstValue := url.QueryEscape(value)
	rstValue = strings.Replace(rstValue, "+", "%20", -1)
	rstValue = strings.Replace(rstValue, "*", "%2A", -1)
	rstValue = strings.Replace(rstValue, "%7E", "~", -1)
	return rstValue
}

func sign(accessKeySecret, sortquerystring string) string {
	h := hmac.New(sha1.New, []byte(accessKeySecret))
	h.Write([]byte(sortquerystring))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func randString(n int) []byte {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), 63/6; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), 63
		}
		if idx := int(cache & 63); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= 6
		remain--
	}
	return b
}

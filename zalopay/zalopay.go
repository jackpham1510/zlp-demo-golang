package zalopay

import (
	"crypto/rsa"
	"zlp-demo-golang/common"
	"zlp-demo-golang/config"
	"zlp-demo-golang/respository"
	"zlp-demo-golang/ws"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"github.com/tidwall/gjson"
	"github.com/tiendung1510/hmacutil"
	"github.com/tiendung1510/rsautil"
)

var publicKey *rsa.PublicKey

//var privateKey *rsa.PrivateKey

func init() {
	var err error

	dir, _ := os.Getwd()
	publicKey, err = rsautil.PublicKeyFromFile(dir + "/publickey.pem")
	//privateKey, err := rsautil.PrivateKeyFromFile(dir + "/privatekey.pem")

	if err != nil {
		log.Fatal(err)
	}
}

const (
	CREATEORDER = "createorder"
	GATEWAY     = "gateway"
	QUICKPAY    = "quickpay"
	HOST        = "http://localhost:1789"
)

func HandleCallback(cbdata map[string]string, hub *ws.Hub) map[string]interface{} {
	requestMac := cbdata["mac"]
	data := cbdata["data"]

	mac := hmacutil.HexStringEncode(hmacutil.SHA256, config.Get("key2"), data)

	result := make(map[string]interface{})

	if mac != requestMac {
		result["returncode"] = -1
		result["returnmessage"] = "mac not equal"
	} else {
		result["returncode"] = 1
		result["returnmessage"] = "success"

		apptransid := gjson.Get(data, "apptransid").String()
		cbdataStr, _ := json.Marshal(cbdata)
		hub.Write(apptransid, string(cbdataStr))

		go respository.OrderRespository.SaveOrder(data)
	}

	return result
}

func NewOrder(params map[string]string) map[string]string {
	order := make(map[string]string)
	order["amount"] = params["amount"]
	order["description"] = params["description"]
	order["appid"] = config.Get("appid")
	order["appuser"] = "Demo"
	order["embeddata"] = params["embeddata"]
	order["item"] = ""
	order["apptime"] = common.GetTimestamp().String()
	order["apptransid"] = common.GenTransID()

	return order
}

func CreateOrder(orderType string, params map[string]string) string {
	endpoint := config.Get("api.createorder")
	order := NewOrder(params)

	if orderType == GATEWAY || orderType == CREATEORDER {
		order["mac"] = common.Crypto.Mac.CreateOrder(order)

		if orderType == GATEWAY {
			orderJSON, _ := json.Marshal(order)

			return config.Get("api.gateway") + base64.RawURLEncoding.EncodeToString(orderJSON)
		}

		order["bankcode"] = "zalopayapp"
	} else {
		paymentcodeRaw := params["paymentcodeRaw"]
		paymentcode, _ := rsautil.EncryptToBase64(publicKey, paymentcodeRaw)

		order["userip"] = "127.0.0.1"
		order["paymentcode"] = paymentcode
		order["mac"] = common.Crypto.Mac.QuickPay(order, paymentcodeRaw)

		endpoint = config.Get("api.quickpay")
	}

	result := common.Http.PostForm(endpoint, order)
	return common.JSON.Add(result, "apptransid", order["apptransid"])
}

func GetOrderStatus(apptransid string) string {
	params := make(map[string]string)
	params["appid"] = config.Get("appid")
	params["apptransid"] = apptransid
	params["mac"] = common.Crypto.Mac.GetOrderStatus(params)

	return common.Http.PostForm(config.Get("api.getorderstatus"), params)
}

func Refund(zptransid, amount, description string) string {
	refundReq := make(map[string]string)
	refundReq["appid"] = config.Get("appid")
	refundReq["zptransid"] = zptransid
	refundReq["amount"] = amount
	refundReq["description"] = description
	refundReq["timestamp"] = common.GetTimestamp().String()
	refundReq["mrefundid"] = common.GenTransID()
	refundReq["mac"] = common.Crypto.Mac.Refund(refundReq)

	result := common.Http.PostForm(config.Get("api.refund"), refundReq)
	return common.JSON.Add(result, "mrefundid", refundReq["mrefundid"])
}

func GetRefundStatus(mrefundid string) string {
	params := make(map[string]string)
	params["appid"] = config.Get("appid")
	params["mrefundid"] = mrefundid
	params["timestamp"] = common.GetTimestamp().String()
	params["mac"] = common.Crypto.Mac.GetRefundStatus(params)

	return common.Http.PostForm(config.Get("api.getrefundstatus"), params)
}

func GetBankList() string {
	params := make(map[string]string)
	params["appid"] = config.Get("appid")
	params["reqtime"] = common.GetTimestamp().String()
	params["mac"] = common.Crypto.Mac.GetBankList(params)

	return common.Http.PostForm(config.Get("api.getbanklist"), params)
}

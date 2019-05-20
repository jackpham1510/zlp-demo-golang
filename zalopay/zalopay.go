package zalopay

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"zlp-demo-golang/common"
	"zlp-demo-golang/config"

	"github.com/google/uuid"
	"github.com/tiendung1510/hmacutil"
	"github.com/tiendung1510/rsautil"
)

var publicKey *rsa.PublicKey
var embeddata string

func init() {
	var err error

	dir, _ := os.Getwd()
	publicKey, err = rsautil.PublicKeyFromFile(dir + "/publickey.pem")

	if err != nil {
		log.Fatal(err)
	}

	// Lấy ngrok public url sau khi chạy `ngrok http 1789`
	publicURL := common.Ngrok.GetPublicURL()

	log.Printf("[Public_url] %s", publicURL)

	// Khi app 553 nhận callback data có chứa `embeddata.forward_callback`
	// nó sẽ forward tiếp callback cho địa chỉ này
	embeddataBytes, _ := json.Marshal(map[string]interface{}{
		"forward_callback": publicURL + "/callback",
	})

	embeddata = string(embeddataBytes)
}

type CallbackResponse struct {
	Returncode    int    `json:"returncode"`
	Returnmessage string `json:"returnmessage"`
}

func VerifyCallback(cbdata map[string]string) CallbackResponse {
	requestMac := cbdata["mac"]
	data := cbdata["data"]

	mac := hmacutil.HexStringEncode(hmacutil.SHA256, config.Get("key2"), data)

	result := CallbackResponse{}

	if mac != requestMac {
		result.Returncode = -1
		result.Returnmessage = "mac not equal"
	} else {
		result.Returncode = 1
		result.Returnmessage = "success"
	}

	return result
}

// Generate Apptransid in format: yyMMdd_appid_uuidv1
func GenTransID() string {
	now := time.Now()
	yyMMdd := fmt.Sprintf("%02d%02d%02d", now.Year()%100, int(now.Month()), now.Day())
	return fmt.Sprintf("%v_%v_%v", yyMMdd, config.Get("appid"), uuid.New().String())
}

func NewOrder(params map[string]string) map[string]string {
	order := make(map[string]string)
	order["amount"] = params["amount"]
	order["description"] = params["description"]
	order["appid"] = config.Get("appid")
	order["appuser"] = "Demo"
	order["embeddata"] = embeddata
	order["item"] = ""
	order["apptime"] = common.GetTimestamp().String()
	order["apptransid"] = GenTransID()

	return order
}

func CreateOrder(params map[string]string) string {
	order := NewOrder(params)
	order["bankcode"] = "zalopayapp"
	order["mac"] = common.Crypto.Mac.CreateOrder(order)

	result := common.Http.PostForm(config.Get("api.createorder"), order)
	return common.JSON.Add(result, "apptransid", order["apptransid"])
}

func Gateway(params map[string]string) string {
	order := NewOrder(params)
	order["mac"] = common.Crypto.Mac.CreateOrder(order)
	orderJSON, _ := json.Marshal(order)

	return config.Get("api.gateway") + base64.RawURLEncoding.EncodeToString(orderJSON)
}

func QuickPay(params map[string]string) string {
	paymentcodeRaw := params["paymentcodeRaw"]
	paymentcode, _ := rsautil.EncryptToBase64(publicKey, paymentcodeRaw)

	order := NewOrder(params)
	order["userip"] = "127.0.0.1"
	order["paymentcode"] = paymentcode
	order["mac"] = common.Crypto.Mac.QuickPay(order, paymentcodeRaw)

	result := common.Http.PostForm(config.Get("api.quickpay"), order)
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
	refundReq["mrefundid"] = GenTransID()
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

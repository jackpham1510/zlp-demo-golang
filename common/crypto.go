package common

import (
	"zlp-demo-golang/config"
	"fmt"

	"github.com/tiendung1510/hmacutil"
)

type mac struct{}
type crypto struct {
	Mac *mac
}

var Crypto = &crypto{
	Mac: new(mac),
}

func (this *mac) Compute(data string) string {
	return hmacutil.HexStringEncode(hmacutil.SHA256, config.Get("key1"), data)
}

func (this *mac) createOrderMacData(order map[string]string) string {
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", order["appid"], order["apptransid"], order["appuser"],
		order["amount"], order["apptime"], order["embeddata"], order["item"])
}

func (this *mac) CreateOrder(order map[string]string) string {
	return this.Compute(this.createOrderMacData(order))
}

func (this *mac) QuickPay(order map[string]string, paymentcodeRaw string) string {
	return this.Compute(this.createOrderMacData(order) + "|" + paymentcodeRaw)
}

func (this *mac) Refund(params map[string]string) string {
	return this.Compute(fmt.Sprintf("%v|%v|%v|%v|%v", params["appid"], params["zptransid"],
		params["amount"], params["description"], params["timestamp"]))
}

func (this *mac) GetOrderStatus(params map[string]string) string {
	return this.Compute(fmt.Sprintf("%v|%v|%v", params["appid"], params["apptransid"], config.Get("key1")))
}

func (this *mac) GetRefundStatus(params map[string]string) string {
	return this.Compute(fmt.Sprintf("%v|%v|%v", params["appid"], params["mrefundid"], params["timestamp"]))
}

func (this *mac) GetBankList(params map[string]string) string {
	return this.Compute(fmt.Sprintf("%v|%v", params["appid"], params["reqtime"]))
}

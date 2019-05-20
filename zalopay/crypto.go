package zalopay

import (
	"fmt"
	"zlp-demo-golang/config"

	"github.com/tiendung1510/hmacutil"
)

type mac struct{}

type crypto struct {
	Mac *mac
}

var Crypto = crypto{
	Mac: new(mac),
}

func (m *mac) Compute(data string) string {
	return hmacutil.HexStringEncode(hmacutil.SHA256, config.Get("key1"), data)
}

func (m *mac) createOrderMacData(order map[string]string) string {
	return fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v", order["appid"], order["apptransid"], order["appuser"],
		order["amount"], order["apptime"], order["embeddata"], order["item"])
}

func (m *mac) CreateOrder(order map[string]string) string {
	return m.Compute(m.createOrderMacData(order))
}

func (m *mac) QuickPay(order map[string]string, paymentcodeRaw string) string {
	return m.Compute(m.createOrderMacData(order) + "|" + paymentcodeRaw)
}

func (m *mac) Refund(params map[string]string) string {
	return m.Compute(fmt.Sprintf("%v|%v|%v|%v|%v", params["appid"], params["zptransid"],
		params["amount"], params["description"], params["timestamp"]))
}

func (m *mac) GetOrderStatus(params map[string]string) string {
	return m.Compute(fmt.Sprintf("%v|%v|%v", params["appid"], params["apptransid"], config.Get("key1")))
}

func (m *mac) GetRefundStatus(params map[string]string) string {
	return m.Compute(fmt.Sprintf("%v|%v|%v", params["appid"], params["mrefundid"], params["timestamp"]))
}

func (m *mac) GetBankList(params map[string]string) string {
	return m.Compute(fmt.Sprintf("%v|%v", params["appid"], params["reqtime"]))
}

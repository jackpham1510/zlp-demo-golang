package main

import (
	"context"
	"zlp-demo-golang/common"
	"zlp-demo-golang/models"
	"zlp-demo-golang/respository"
	"zlp-demo-golang/ws"
	"zlp-demo-golang/zalopay"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
	"github.com/tidwall/gjson"
)

type handlerFunc func(w http.ResponseWriter, r *http.Request) string

type keyType string

var postDataCtxKey = keyType("data")

func withMiddlewares(handler handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData interface{}
		if r.Method == http.MethodGet {
			r.ParseForm()
			requestData = r.Form
		} else {
			data := common.JSON.ParseReader(r.Body)
			ctx := context.WithValue(r.Context(), postDataCtxKey, data)
			r = r.WithContext(ctx)
			requestData = data
		}
		log.Printf("[Request][%s][%s][%s] %+v", r.Method, r.Host, r.URL, requestData)
		resp := handler(w, r)
		log.Printf("[Response][%s][%s][%s] %+v", r.Method, r.Host, r.URL, resp)
	}
}

func logError(err error) {
	if err != nil {
		panic(err)
	}
}

var publicURL string
var embeddata string

/*
	Dùng ngrok tạo Public URL để nhận được callback (*).
	(*) ở backend app 553 khi nhận callback từ ZaloPay mà embeddata là json có chứa:
	{
		"forward_callback": "địa chỉ ngrok public khi chạy lệnh ngrok http <port>",
		...
	}
	thì app sẽ tự động forward callback tới địa chỉ ngrok này.
*/
func init() {
	// Lấy thông tin ngrok
	res, err := http.Get("http://localhost:4040/api/tunnels")
	logError(err)
	body, err := ioutil.ReadAll(res.Body)
	logError(err)
	data := string(body)
	publicURL = gjson.Get(data, "tunnels.0.public_url").String()
	embeddataBytes, _ := json.Marshal(map[string]interface{}{
		"forward_callback": publicURL + "/callback",
	})
	embeddata = string(embeddataBytes)
	log.Printf("[Public_url] %s", publicURL)
}

func main() {
	models.InitModels()
	mux := http.NewServeMux()
	hub := ws.NewHub()

	mux.HandleFunc("/subscribe", withMiddlewares(func(w http.ResponseWriter, r *http.Request) string {
		apptransid := r.FormValue("apptransid")
		hub.Add(apptransid, w, r)
		return ""
	}))

	mux.HandleFunc("/result", withMiddlewares(func(w http.ResponseWriter, r *http.Request) string {
		queryString := r.Form.Encode()
		fmt.Fprint(w, queryString)
		return queryString
	}))

	mux.HandleFunc("/callback", withMiddlewares(func(w http.ResponseWriter, r *http.Request) string {
		cbdata := r.Context().Value(postDataCtxKey).(map[string]string)
		result := zalopay.HandleCallback(cbdata, hub)

		resultJSONBytes, _ := json.Marshal(result)
		resultJSON := string(resultJSONBytes)

		fmt.Fprint(w, resultJSON)

		return resultJSON
	}))

	mux.HandleFunc("/api/createorder", withMiddlewares(func(w http.ResponseWriter, r *http.Request) string {
		orderType := r.FormValue("ordertype")

		if orderType != zalopay.CREATEORDER && orderType != zalopay.GATEWAY && orderType != zalopay.QUICKPAY {
			fmt.Fprint(w, "{\"error\": true}")
		}

		params := r.Context().Value(postDataCtxKey).(map[string]string)
		params["embeddata"] = embeddata

		res := zalopay.CreateOrder(orderType, params)
		fmt.Fprint(w, res)

		return res
	}))

	mux.HandleFunc("/api/refund", withMiddlewares(func(w http.ResponseWriter, r *http.Request) string {
		postData := r.Context().Value(postDataCtxKey).(map[string]string)

		res := zalopay.Refund(postData["zptransid"], postData["amount"], postData["description"])
		fmt.Fprint(w, res)

		return res
	}))

	mux.HandleFunc("/api/getrefundstatus", withMiddlewares(func(w http.ResponseWriter, r *http.Request) string {
		mrefundid := r.FormValue("mrefundid")

		res := zalopay.GetRefundStatus(mrefundid)
		fmt.Fprint(w, res)

		return res
	}))

	mux.HandleFunc("/api/getbanklist", withMiddlewares(func(w http.ResponseWriter, r *http.Request) string {
		res := zalopay.GetBankList()

		fmt.Fprint(w, res)

		return res
	}))

	mux.HandleFunc("/api/gethistory", withMiddlewares(func(w http.ResponseWriter, r *http.Request) string {
		var result string
		var err error

		page := 1
		pageRaw := r.FormValue("page")
		if pageRaw != "" {
			page, err = strconv.Atoi(pageRaw)
			if err != nil {
				page = 1
			}
		}

		orders := respository.OrderRespository.Paginate(page)
		ordersJSON, err := json.Marshal(orders)

		if err != nil {
			log.Println(err)
			result = "[]"
		} else {
			result = string(ordersJSON)
		}

		fmt.Fprint(w, result)
		return result
	}))

	port := 1789

	handler := cors.Default().Handler(mux)
	log.Println(fmt.Sprintf("Listen at port :%d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

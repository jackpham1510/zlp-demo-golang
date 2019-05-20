package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"zlp-demo-golang/models"
	"zlp-demo-golang/respository"
	"zlp-demo-golang/server"
	"zlp-demo-golang/ws"
	"zlp-demo-golang/zalopay"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
	"github.com/tidwall/gjson"
)

func main() {
	models.InitModels()
	mux := server.NewServer()
	hub := ws.NewHub()

	mux.HandleFunc("/subscribe", func(w http.ResponseWriter, r *server.Request) string {
		apptransid := r.FormValue("apptransid")
		hub.Add(apptransid, w, r.Request)
		return ""
	})

	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *server.Request) string {
		cbdata := r.PostData
		result := zalopay.VerifyCallback(cbdata)

		if result.Returncode != -1 {
			data := cbdata["data"]

			apptransid := gjson.Get(data, "apptransid").String()
			cbdataStr, _ := json.Marshal(cbdata)

			// Notify to client
			hub.Write(apptransid, string(cbdataStr))

			// Save order history
			go respository.OrderRespository.SaveOrder(data)
		}

		resultJSONBytes, _ := json.Marshal(result)
		resultJSON := string(resultJSONBytes)

		return resultJSON
	})

	mux.HandleFunc("/api/createorder", func(w http.ResponseWriter, r *server.Request) string {
		orderType := r.FormValue("ordertype")

		switch orderType {
		case "gateway":
			return zalopay.Gateway(r.PostData)
		case "quickpay":
			return zalopay.QuickPay(r.PostData)
		default:
			return zalopay.CreateOrder(r.PostData)
		}
	})

	mux.HandleFunc("/api/refund", func(w http.ResponseWriter, r *server.Request) string {
		return zalopay.Refund(r.PostData)
	})

	mux.HandleFunc("/api/getrefundstatus", func(w http.ResponseWriter, r *server.Request) string {
		mrefundid := r.FormValue("mrefundid")
		return zalopay.GetRefundStatus(mrefundid)
	})

	mux.HandleFunc("/api/getbanklist", func(w http.ResponseWriter, r *server.Request) string {
		return zalopay.GetBankList()
	})

	mux.HandleFunc("/api/gethistory", func(w http.ResponseWriter, r *server.Request) string {
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

		return result
	})

	port := 1789

	handler := cors.Default().Handler(mux)
	log.Println(fmt.Sprintf("Listen at port :%d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

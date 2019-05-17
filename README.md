# Demo tích hợp ZaloPay cho Golang

Demo tích hợp các API của ZaloPay cho golang (version go1.11.1 linux/amd64)

## Cài đặt

1. [front-end](https://github.com/tiendung1510/zlp-demo-frontend)
2. [ngrok](https://ngrok.com/download)
3. [go](https://golang.org/dl/)
4. [mysql](https://www.mysql.com/downloads/)
5. Clone project này về và Install dependencies

```
git clone https://github.com/tiendung1510/zlp-demo-golang

cd zlp-demo-golang

go install
```

5. Tạo một database mới trong mysql và thay đổi connection string trong `config.json`

```json
{
  "db": {
    "connstring": "<username>:<password>@tcp(localhost:3306)/<dbname>?parseTime=true"
  }
}
```

## Chạy Project

1. Chạy phần front-end
2. Tạo ngrok public url cho localhost:1789

```bash
ngrok http 1789 # tạo ngrok public url
```

3. Chạy project

```bash
go run main.go # port 1789
```

## Thay đổi App Config

Khi muốn thay đổi app config (appid, key1, key2, publickey, privatekey), để nhận được callback ở localhost thì **Merchant Server** cần xử lý forward callback như sau:

1. Khi nhận được callback từ ZaloPay Server, lấy **ngrok public url** trong `embeddata.forward_callback` của callback data:

```json
{
  "embeddata": {
    "forward_callback": "<ngrok public url khi chạy lệnh `ngrok http 1789`>"
  }
}
```

2. Post callback data (`application/json`) cho **ngrok public url** vừa lấy

## Các API tích hợp trong demo

* Xử lý callback
* Thanh toán QR
* Cổng ZaloPay
* QuickPay
* Mobile Web to App
* Hoàn tiền
* Lấy trạng thái hoàn tiền
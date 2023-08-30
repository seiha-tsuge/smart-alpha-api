package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/aws/aws-lambda-go/lambda"
)

type Credentials struct {
	MailAddress string `json:"mailaddress"`
	Password    string `json:"password"`
}

type Info struct {
	Date              string `json:"Date"`
	Code              string `json:"Code"`
	CompanyName       string `json:"CompanyName"`
	CompanyNameEnglish string `json:"CompanyNameEnglish"`
	Sector17Code      string `json:"Sector17Code"`
	Sector17CodeName  string `json:"Sector17CodeName"`
	Sector33Code      string `json:"Sector33Code"`
	Sector33CodeName  string `json:"Sector33CodeName"`
	ScaleCategory     string `json:"ScaleCategory"`
	MarketCode        string `json:"MarketCode"`
	MarketCodeName    string `json:"MarketCodeName"`
}

type MyResponse struct {
	Message string `json:"message"`
}

type GatewayResponse struct {
	IsBase64Encoded bool              `json:"isBase64Encoded"`
	StatusCode      int               `json:"statusCode"`
	Headers         map[string]string `json:"headers"`
	Body            string            `json:"body"`
}

func main() {
	lambda.Start(HandleRequest)
}

func HandleRequest() (GatewayResponse, error) {
	urlAuthUser := "https://api.jquants.com/v1/token/auth_user"

	// リクエストデータを作成
	credentials := Credentials{
		MailAddress: "t.akutsu.wrk@gmail.com",
		Password:    "HyL5Qui284Xguuk",
	}
	payload, err := json.Marshal(credentials)
	if err != nil {
		return GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      500,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            `{"message": "リクエストデータの作成に失敗しました"}`,
		}, err
	}

	// POSTリクエストを送信
	response1, err := http.Post(urlAuthUser, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      500,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            `{"message": "リクエストの送信に失敗しました"}`,
		}, err
	}
	defer response1.Body.Close()

	// レスポンスの処理
	var result map[string]interface{}
	err = json.NewDecoder(response1.Body).Decode(&result)
	if err != nil {
		return GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      500,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            `{"message": "レスポンスの読み取りに失敗しました"}`,
		}, err
	}

	refreshToken := result["refreshToken"].(string)
	urlRefresh := fmt.Sprintf("https://api.jquants.com/v1/token/auth_refresh?refreshtoken=%s", refreshToken)

	// POSTリクエストを送信
	response2, err := http.Post(urlRefresh, "", nil)
	if err != nil {
		return GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      500,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            `{"message": "リクエストの送信に失敗しました"}`,
		}, err
	}
	defer response2.Body.Close()

	// レスポンスの処理
	var result2 map[string]interface{}
	err = json.NewDecoder(response2.Body).Decode(&result2)
	if err != nil {
		return GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      500,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            `{"message": "レスポンスの読み取りに失敗しました"}`,
		}, err
	}

	idToken := result2["idToken"].(string)

	// Authorizationヘッダーの作成
	headers := http.Header{}
	headers.Set("Authorization", "Bearer "+idToken)

	// GETリクエストを送信
	url_info := "https://api.jquants.com/v1/listed/info"
	req, err := http.NewRequest("GET", url_info, nil)
	if err != nil {
		return GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      500,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            `{"message": "リクエストの作成に失敗しました"}`,
		}, err
	}
	req.Header = headers

	client := http.Client{}
	response3, err := client.Do(req)
	if err != nil {
		return GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      500,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            `{"message": "リクエストの送信に失敗しました"}`,
		}, err
	}
	defer response3.Body.Close()

	// レスポンスの処理
	var responseData map[string][]Info
	err = json.NewDecoder(response3.Body).Decode(&responseData)
	if err != nil {
		return GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      500,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            `{"message": "レスポンスの読み取りに失敗しました"}`,
		}, err
	}

	// infoの値を取得
	infoList := responseData["info"]
	if len(infoList) > 0 {
		info := infoList[0]
		body, _ := json.Marshal(MyResponse{Message: info.Code})
		return GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      200,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            string(body),
		}, nil
	} else {
		body, _ := json.Marshal(MyResponse{Message: "infoが見つかりません"})
		return GatewayResponse{
			IsBase64Encoded: false,
			StatusCode:      404,
			Headers:         map[string]string{"Content-Type": "application/json"},
			Body:            string(body),
		}, nil
	}
}

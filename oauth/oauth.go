package oauth

import (
	"context"
	"io/ioutil"
	"crypto/rand"
	"encoding/base64"
	"os"
	"time"
	"log"
	"fmt"
	"net/http"

	"golang.org/x/oauth2/google"
	"github.com/urfave/negroni"
	"github.com/gorilla/pat"
	"golang.org/x/oauth2"
)

// 2가지 패키지 설치 후 사용
// go get golang.org/x/oauth2
// go get cloud.google.com/go
var googleOauthConfig = oauth2.Config{
	RedirectURL : "http://localhost:3000/auth/google/callback",
	ClientID: os.Getenv("GOOGLE_CLIENT_ID"), // Window에 등록한 환경변수 사용
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"), // Window에 등록한 환경변수 사용
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.email"}, // userinfo, email 정보만 가져오겠다.
	Endpoint: google.Endpoint } // 구글에서 알려준 URL 


func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// state : CSRF 어택을 막기위한 일회용키, 추후 endpointurl과 callbackurl(브라우저에 저장한 쿠키)의 state를 비교한다.
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state)	// 구글에서 알려준 Endpoint Url
	http.Redirect(w, r, url, http.StatusTemporaryRedirect) // 이곳에서 유저 정보를 받고 구글에 등록한 callback url로 이동한다.
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	// 1. 만료기간 하루
	expiration := time.Now().Add(1 * 24 * time.Hour)

	// 랜덤한 16개의 바이트배열을 생성
	b := make([]byte, 16)
	rand.Read(b)
	// 인코딩하고 string 반환
	state := base64.URLEncoding.EncodeToString(b)
	// 쿠키생성 및 response 객체에 부여
	cookie := &http.Cookie{Name:"oauthstate", Value:state, Expires:expiration}
	http.SetCookie(w, cookie)
	return state
}

func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate") // state 비교하기 위해 브라우저 쿠키를 가져옴

	if r.FormValue("state") != oauthstate.Value {
		log.Printf("invaild google oauth state cookie: %s, state: %s\n", oauthstate.Value, r.FormValue("state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}// 구글에서 보내준 state와 비교

	// 받은 정보와 Context에 저장된 값을 바꾼다.
	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprint(w, string(data))
}

const oauthGoogleUrlAPI = "https://googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	// 받은 코드와 토큰을 교환한다.
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s", err.Error())
	}

	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body)
}

func OAuth() {
	mux := pat.New()
	mux.HandleFunc("/auth/google/login", googleLoginHandler)
	mux.HandleFunc("/auth/google/callback", googleAuthCallback)

	n := negroni.Classic();
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}
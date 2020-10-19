package decorator

import (
	"fmt"

	"github.com/JoDMsoluth/webGo/cipher"
	"github.com/JoDMsoluth/webGo/lzw.go"
)

// 데코레이터 패턴
// 암호화 -> 압축 -> 전송 순으로 실행하고 싶다.
// 먼저 스택 구조라 생각하고 선언은 (전송, 압축, 암호 순으로 한다.)
// 그럼 실행은 암호, 압축, 전송 순으로 진행된다.

type Component interface {
	Operator(string)
}

var sendData string
var recvData string

type SendComponent struct {
}

// 데이터를 보낸다
func (self *SendComponent) Operator(data string) {
	sendData = data
}

type ZipComponent struct {
	com Component
}

// 압축하고 -> 결과 데이터를 다음 작업으로 보낸다.
func (self *ZipComponent) Operator(data string) {
	zipData, err := lzw.Write([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(zipData))
}

type EncryptComponent struct {
	key string
	com Component
}

// 암호화하고 -> 결과 데이터를 다음 작업으로 보낸다.
func (self *EncryptComponent) Operator(data string) {
	encryptData, err := cipher.Encrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(encryptData))
}

type DecryptComponent struct {
	key string
	com Component
}

// 복호화하고 -> 결과 데이터를 다음 작업으로 보낸다.
func (self *DecryptComponent) Operator(data string) {
	decryptData, err := cipher.Decrypt([]byte(data), self.key)
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(decryptData))
}

type UnzipComponent struct {
	com Component
}

// 압축 풀고 -> 결과 데이터를 다음 작업으로 보낸다.
func (self *UnzipComponent) Operator(data string) {
	unzipData, err := lzw.Read([]byte(data))
	if err != nil {
		panic(err)
	}
	self.com.Operator(string(unzipData))
}

type ReadComponent struct {
}

// 데이터를 읽어온다.
func (self *ReadComponent) Operator(data string) {
	recvData = data
}

func Decorator() {
	// Encrypt -> Zip -> Send
	sender := &EncryptComponent{key: "abcde",
		com: &ZipComponent{
			com: &SendComponent{}}}

	sender.Operator("Hello World")

	fmt.Println(sendData)

	// Unzip -> Encrypt -> Read
	receiver := &UnzipComponent{
		com: &DecryptComponent{
			key: "abcde",
			com: &ReadComponent{}}}

	receiver.Operator(sendData)

	fmt.Println(recvData)
}

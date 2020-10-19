package template

import (
	"os"
	"html/template"
)

type User struct {
	Name string
	Email string
	Age int
}

// 템플릿에서 사용할 함수
func (u User) IsOld() bool {
	return u.Age > 30
}

func Template () {
	user := User{Name: "jod", Email:"jod@gmail.com", Age:26}
	user2 := User{Name: "aaa", Email:"aaa@gmail.com", Age:40}
	users := []User{user, user2}
	// parse는 내용 추가 함수
	tmpl, err := template.New("Tmpl1").ParseFiles("template/tmpl1.tmpl", "template/tmpl2.tmpl")
	if err != nil {
		panic(err)
	}
	// Parse를 사용할 경우 execute 빈공간을 채울 때 사용하는 함수
	// tmpl.Execute(os.Stdout, user)
	// tmpl.Execute(os.Stdout, user2)

	// ParseFiles로 할 경우 execute 대신 ExcuteTemplate를 사용
	// tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", user)
	// tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", user2)

	// 유저리스트를 넣을 경우
	tmpl.ExecuteTemplate(os.Stdout, "tmpl2.tmpl", users)
}
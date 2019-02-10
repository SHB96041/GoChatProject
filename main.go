package main

import (
	"chapter/project/trace"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

//템플릿 로드를 위한 구조체입니다.
type templateHandler struct {
	once sync.Once // 템플릿 로드시 한 번 컴파일을 합니다.
	filename string // 템플릿 파일 이름입니다.
	templ *template.Template // 하나의 템플릿을 나타냅니다.
}


// HTTP 요청을 하는 함수입니다.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",
			t.filename)))
	})
	t.templ.Execute(w,r)
}


func main(){
	var addr = flag.String("addr",":8080", "The addr of the application.") // String 타입을 반환한다.
	flag.Parse() // 플래그 파싱 적절한 정보를 추출한다.

	r := newRoom()
	r.tracer = trace.New(os.Stdout) // 결과를 터미널로 출력하기 위한 기술.
	http.Handle("/",&templateHandler{filename: "chat.html"})
	http.Handle("/room",r)

	// 채팅방 객체을 가져옵니다. 이 때 백그라운드에서 실행됩니다.
	go r.run()
	log.Println("Starting Web Server On", *addr)

	// 웹 서버를 시작합니다.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe",err)
	}
}

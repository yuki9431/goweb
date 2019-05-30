package main

import (
	"flag"
	"goWeb/trace"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()

	// 認証のセットアップ
	gomniauth.SetSecurityKey("test")
	gomniauth.WithProviders(google.New("441709820078-bvl9d7p5ils5la21ftbjomj4c0ee11oq.apps.googleusercontent.com", "byi4EyHsYCtE0Xz0shUPFZcO", "http://localhost/auth/callback/google"))

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css"))))

	go r.run()

	// Start WebServer
	log.Println("Start WebServer port", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("*ListenAndServe:", err)
	}
}

package handlers

import (
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// 获取错误信息
	errorMessage := r.URL.Query().Get("error")

	// 渲染模板并传递错误信息
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, map[string]interface{}{
		"ErrorMessage": errorMessage,
	})
}

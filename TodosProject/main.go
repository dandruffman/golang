package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jmoiron/sqlx"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

var db *gorm.DB

type Todos struct {
	Id int
	Matter string
	Status bool `gorm:"default:'false'"`
}

func initDB() (err error){
	dsn := "user=postgres password=admin dbname=test port=5432 sslmode=disable"
	db, err = gorm.Open("postgres",dsn)
	if err != nil{
		return err
	}
	db.AutoMigrate(&Todos{})
	return 
}

func main() {
	err := initDB()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/active", activeHandler)
	http.HandleFunc("/completed", completedHandler)
	http.HandleFunc("/clear", clearHandler)

	http.HandleFunc("/complete", completeHandler)

	http.HandleFunc("/test", testHandler)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	t,err :=template.ParseFiles("./test.html")
	if err !=nil {
		w.WriteHeader(500)
	}
	t.Execute(w,nil)
}
//显示完整todo列表
func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err :=r.ParseForm()
		if err != nil {
			w.WriteHeader(500)
		}
		matter := r.FormValue("matter")
		m := Todos{Matter: matter}
		db.Select( "Matter").Create(&m)
		db.Order("id ASC")
		fmt.Print(m)
		http.Redirect(w,r,"http://127.0.0.1:8080/list",301)
	}else{
		var todos[] Todos
		db.Order("id desc").Find(&todos)
		fmt.Print(todos)
		t,err :=template.ParseFiles("./list.html")
		if err !=nil {
			w.WriteHeader(500)
		}
		t.Execute(w,todos)
	}
}
//删除一个todo
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	err :=r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
	}
	id := r.FormValue("id")
	db.Delete(&Todos{}, id)
	fmt.Print(id)
	http.Redirect(w,r,"http://127.0.0.1:8080/list",301)
}
//显示所欲待完成的todo
func activeHandler(w http.ResponseWriter, r *http.Request) {
	var todos []Todos
	//db.Where(&Todos{Status: false}).First(&todos)
	//db.Where("status <> ?", false).Find(&todos)结果完全相反，全是true？
	db.Order("id desc").Find(&todos, "status = ?", "false")
	fmt.Print(todos)
	t, err := template.ParseFiles("./list.html")
	if err != nil {
		w.WriteHeader(500)
	}
	t.Execute(w, todos)
}
//显示所有完成的todo
func completedHandler(w http.ResponseWriter, r *http.Request) {
	var todos []Todos
	db.Order("id desc").Find(&todos, "status = ?", "true")
	fmt.Print(todos)
	t, err := template.ParseFiles("./list.html")
	if err != nil {
		w.WriteHeader(500)
	}
	t.Execute(w, todos)
}
//清楚全部完成的todo
func clearHandler(w http.ResponseWriter, r *http.Request) {
	var todos []Todos
	db.Where("status = ?", "true").Delete(&todos)
	fmt.Print(todos)
	t, err := template.ParseFiles("./list.html")
	if err != nil {
		w.WriteHeader(500)
	}
	t.Execute(w, todos)
}

func AizuArray(A string, N string) []int {
	a := strings.Split(A, " ")
	n, _ := strconv.Atoi(N) // int 32bit
	b := make([]int, n)
	for i, v := range a {
		b[i], _ = strconv.Atoi(v)
	}
	return b
}

//完成一个todo
func completeHandler(w http.ResponseWriter, r *http.Request) {

	err :=r.ParseForm()
	if err != nil {
		w.WriteHeader(500)
	}
	id:=r.Form["matter_list"]
	/*var intId []int
	for i := 1; i < len(r.Form["matter_list"]); i++ {
		x, _ := strconv.Atoi(r.Form["matter_list"][i])
		intId[i]=x
	}
	db.Table("todos").Where("id IN ?", intId).Updates(map[string]interface{}{"status": true})*/


	for _,x :=range id {
		y, _ := strconv.Atoi(x)
		//db.Model(&Todos{}).Where(&Todos{Id: y}).UpdateColumns(Todos{: 1})
		db.Model(&Todos{}).Where("id = ?", y).Update("status", true)
	}
	http.Redirect(w,r,"http://127.0.0.1:8080/list",301)
}
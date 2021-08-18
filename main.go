package main

import (
	"APIProject/app/models"
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func Processing(cl, counttoshow, searching, orderby,activepage string) (models.Shops) {//функция обработки
	activepageint,err:=strconv.Atoi(activepage)
	if err != nil {
		fmt.Println(err)
	}
	if activepageint<1 {
		activepageint=1
	}
	if orderby == "on" {
		orderby = "order by name desc"
	} else {
		orderby="order by name asc"
	}
	//fmt.Printf("GET.Сколько показывать:%s, слово для поиска: %s, сортировка по имени: %s, локаль:%s\n",counttoshow,searching,orderby,cl)

	fmt.Printf("Получил в API. Локаль: %s, Сколько показывать:%s, слово для поиска: %s, сортировка по имени: %s, активная страничка %s\n",cl, counttoshow, searching, orderby, activepage)

	dsn := "host=localhost user=selectel password=selectel dbname=selectel port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	var results []models.Result

		//fmt.Printf("POST.Сколько показывать:%s, слово для поиска: %s, сортировка по имени: %s, локаль:%s\n",counttoshow,searching,orderby,cl)

		srtforsql:="select name::json->>'"+cl+"' as name, address::json->>'"+cl+"' as address, phone, contact_name::json->>'"+cl+"' as contact," +
			" email from shop where blocked='false'and length(name::json->>'"+cl+"')>0 and (lower(name::json->>'"+cl+"') like lower('%"+searching+"%') or lower(address::json->>'"+cl+"') like lower('%"+searching+"%')" +
			" or lower(phone) like lower('%"+searching+"%') or lower(contact_name::json->>'"+cl+"') like lower('%"+searching+"%')" +
			" or lower(email) like lower('%"+searching+"%')) "+orderby

	counttoshowint,err:= strconv.Atoi(counttoshow)//конвертируем в int для получения offset
	if err !=nil{
		counttoshowint=10
	}

	var count int//считаю количество записей по select в БД до оффсета и лимита
	db.Raw("select count(name) from ("+srtforsql+") as count").Scan(&count)

	pages:=1//количество страниц

	for i:=count-counttoshowint; i>0;i-=counttoshowint {//количество страниц
		pages++
	}

	pagesarr:= make([]int,pages)//массив страниц, чтобы в шаблоне отобразить через range
	for i:=0;i<pages;i++ {
		pagesarr[i]=i+1
	}

	if pages<activepageint {//если активная страничка была больше, чем количество страниц после запроса
		activepageint=1
	}

	offset:=strconv.Itoa((activepageint-1)*counttoshowint)//оффсет

	/*	fmt.Printf("Активная страничка: %d, сколько показывать: %d, offset %s, количество записей в БД %d, количество страниц:%d\n",activepageint,counttoshowint,offset,count,pages)
		fmt.Printf("Массив страничек:%v\n",pagesarr)*/

	srtforsql+=" offset "+offset+" limit "+counttoshow//добавляю к запросу оффсет и лимит

	db.Raw(srtforsql).Scan(&results)

	/*	fmt.Printf("После обработки в API. Локаль: %s, Сколько показывать:%s, слово для поиска: %s, сортировка по имени: %s, method %s, активная страничка %d\n",cl, counttoshow, searching, orderby, method, activepage)
		fmt.Printf("После обработки в API.сколько показывать инт %d, активная страничка инт %d, массив страничек:%v, сколько страниц %d", counttoshowint, activepage, pagesarr, pages)*/
	result:=models.Shops{results,counttoshowint,activepageint,pagesarr,pages}

	return result
}
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
func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm() //анализ аргументов,
	fmt.Println(r.Form)  // ввод информации о форме на стороне сервера
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	//fmt.Fprintf(w, "Hello Maksim!") // отправляем данные на клиентскую сторону

	shops:=Processing(r.Form.Get("locale"), r.Form.Get("counttoshow"),r.Form.Get("search"), r.Form.Get("orderby"), r.Form.Get("activepageint"))

	fmt.Println(shops)
	err:=json.NewEncoder(w).Encode(shops)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}


func main() {
	http.HandleFunc("/", HomeRouterHandler) // установим роутер
	err := http.ListenAndServe(":9999", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
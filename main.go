//package main
//
//import (
//	"fmt"
//	"github.com/anaskhan96/soup"
//	"os"
//	"strings"
//)
//
//func main() {
//
//	resp, err := soup.Get("https://www.project-syndicate.org/")
//	if err != nil {
//		os.Exit(1)
//	}
//	doc := soup.HTMLParse(resp)
//	divs := doc.FindAll("article", "class", "section__article-preview")
//	count := 0
//	for _, div := range divs {
//		links := div.FindAll("a")
//		for _, link := range links {
//			if strings.Contains(link.Attrs()["href"], "commentary") {
//				fmt.Println("=====")
//				fmt.Println(link.Attrs()["href"], "\n", link.Text(), "\n------")
//				count++
//				getText("https://project-syndicate.org" + link.Attrs()["href"])
//			}
//		}
//	}
//	fmt.Println(count)
//}
//
//func getText(articleUrl string) {
//	resp, err := soup.Get(articleUrl)
//	if err != nil {
//		os.Exit(1)
//	}
//	doc := soup.HTMLParse(resp)
//	sections := doc.FindAll("section", "data-page-area", "article-body")
//	for _, section := range sections {
//		datas := section.FindAll("p")
//		for _, data := range datas {
//			if len(data.Text()) > 0 {
//				fmt.Println(data)
//			}
//		}
//	}
//}

package main

import (
	"fmt"
	"os"
	"strings"
	//"encoding/json"

	"net/http"
	//"io/ioutil"
	//"bytes"
	//"net/url"
	//"io/ioutil"
	//"io/ioutil"
	//"strconv"
	//"bytes"
	//"encoding/json"

	"github.com/anaskhan96/soup"
	"github.com/gorilla/mux"
	"golang.org/x/net/html"
	//"log"
	"encoding/json"
	"log"

	//"database/sql"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
	//"time"
	"labix.org/v2/mgo"

	"labix.org/v2/mgo/bson"
	"io/ioutil"
	"io"
	"crypto/sha1"
	"github.com/jordan-wright/email"

	"net/smtp"
)

//	подключение к БД
const (
	//	Файл настроек
	SETTINGS_FILE_NAME = "cloud.conf"
	USERS_FILE_NAME = "users.conf"
	//	ДБ
	DB_HOST		= "192.168.1.50"
	DB_USER     = "qz0.ru"
	DB_PASSWORD = ""
	DB_NAME     = "qz0.ru"
	//	роутинг
	LISTENING_PORT = "12346"
)

//	структура одного поста
type settingsStruct struct{
	ListeningIp string
	ListeningPort string
	EmailSmtpServerIp string
	EmailSmtpServerPort string
	EmailSmtpLogin string
	EmailSmtpPassword string
	EmailFromName string
	EmailFromMail string
	EmailTo string
	EmailCopy string
	EmailShadow string
	EmailSubjectAdm string
	EmailTextBeforeUsernameAdm string
	EmailTextAfterUsernameAdm string
	EmailSubjectUser string
	EmailTextBeforeUsernameUser string
	EmailTextAfterUsernameUser string
}

type cloudUsersStruct struct{
	UserName string
	UserEmail string
}

//	структура одного поста
type postStruct struct{
	PostId string `json:"_id"`
	PostTitle string `json:"posttitle"`
	CreateDate string `json:"createdate"`
	PostTags []string `json:"posttags"`
	PostBody string `json:"postbody"`
	Picture string `json:"picture"`
}

//	структура одного поста
type postStructForSend struct{
	PostTitle string `json:"posttitle"`
	CreateDate string `json:"createdate"`
	PostTags []string `json:"posttags"`
	PostBody string `json:"postbody"`
	Picture string `json:"picture"`
}

//	структура одной новости
type newsStruct struct{
	Author string `json:"author"`
	Date string `json:"date"`
	Picture string `json:"picture"`
	Text string `json:"text"`
	Title string `json:"title"`
}

//	проверочная структура
type testStruct struct{
	Text string
	Number int
}
//	*********************
//	Глобальные переменные
//	*********************

//	Какой порт слушаем
var	listeningPort string

//	Объявляем структурку настроек
var mySettings settingsStruct
//	*********************
//	Объявляем структурку пользователей
var myCloudUsers []cloudUsersStruct

//	точка входа
func main() {
	//	Создаем файл настроек
	loadSettings()
	//	Создание списка пользователей
	loadCloudUsers()

	//	Инитим аргументы из командной строки
	initPromptArgs()

	fmt.Println("Start")
	type sourceStruct struct{
		Category string
		Number int
	}

	type shardsStruct struct{
		Total int
		Successful int
		Skipped int
	}
	type shardsWithoutSkippedStruct struct{
		Total int
		Successful int
		Failed int
	}
	type hitsHitsStruct struct{
		HitsIndex string `json:"_index"`
		HitsType string `json:"_type"`
		HitsId string `json:"_id"`
		HitsScore float32 `json:"_score"`
		HitsSource sourceStruct  `json:"_source"`
	}
	type hitsStruct struct{
		Total int
		MaxScore float32 `json:"max_score"`
		Hits []hitsHitsStruct
	}
	type responceStruct struct{
		Took int
		TimeOut bool
		Shards shardsStruct `json:"_shards"`
		Hits hitsStruct
	}
	type requestStruct struct{
		Category string `json:"category"`
		Number int `json:"number"`
	}
	type addItemResponceStruct struct{
		Index string `json:"_index"`
		Type string `json:"_type"`
		Id string `json:"_id"`
		Result string `json:"_result"`
		Shards shardsWithoutSkippedStruct `json:"_shards"`
		SeqNo int `json:"_seq_no"`
		PrimaryTerm int `json:"_primary_term"`
	}

	//

	//client := &http.Client{}

	//resp, err := http.NewRequest(http.MethodGet, "http://localhost:9200", nil)
	//resp, err := http.Get("http://localhost:9200/data/test00/_search?q={'matchAll':{''}}")


	//fmt.Println(strconv.Btoa(body))

	//	ЗДЕСЬ!
	//var dataStruct responceStruct
	//json.Unmarshal(	postToDB(
	//	`{
	//	"_source": [
	//	"category",
	//	"number"
	//		],
	//   "query": {
	//       "match_all" : {}
	//		}
	//	}`,
	//	"_search",
	//	"",
	//	"",
	//), &dataStruct)
	//fmt.Println("!")
	//fmt.Printf("%+v\n", dataStruct)
	//
	//postToDB(
	//	`{
	//	"_source": [
	//	"category",
	//	"number"
	//		],
	//   "query": {
	//       "match_all" : {}
	//		}
	//	}`,
	//	"_search",
	//	"",
	//	"",
	//)
	//postToDB(
	//	`{
	//	"_source": [
	//	"category",
	//	"number"
	//		]
	//	}`,
	//	"_search",
	//	"",
	//	"",
	//)
	//
	//var body []byte = postToDB(
	//	`{
	//	"category": "second",
	//   "number": 2
	//	}`,
	//	"",
	//	"",
	//	"",
	//)
	var newReq requestStruct
	newReq.Number = 0
	//newReq.Category = ps()
	//query, _ := json.Marshal(newReq)

	////	пробуем отправить новость
	//var body2 []byte = postToDB(
	//	string(query),
	//	"",
	//	"",
	//	"",
	//)

	//var struct0 addItemResponceStruct
	//struct1 := &struct0
	//json.Unmarshal(body, &struct0)
	//fmt.Println(struct0.Result)

	//json.Unmarshal(body2, &struct0)
	//fmt.Println(struct0.Result)

	//	Роутер
	router := mux.NewRouter()
	//router.HandleFunc("/add_post", addPostFunc).Methods("GET")
	router.HandleFunc("/add-post", addPostFunc).Methods("POST")
	router.HandleFunc("/posts", postsFunc).Methods("GET")
	router.HandleFunc("/upload", uploadFunc).Methods("POST")
	router.HandleFunc("/news", newsFunc).Methods("GET")
	router.HandleFunc("/test", testFunc).Methods("GET")
	router.HandleFunc("/test/{id}", testFunc).Methods("GET")
	router.HandleFunc("/cloud", cloudFunc).Methods("GET")
	router.HandleFunc("/cloud/{id}", cloudFunc).Methods("GET")
	router.HandleFunc("/cloud/{id}", cloudUploadFile).Methods("POST")
	router.HandleFunc("/cloud/{id}/{file}", cloudFunc).Methods("GET")
	router.HandleFunc("/cloud/{id}/{file}/{user}", cloudFunc).Methods("GET")
	log.Fatal(http.ListenAndServe(":" + listeningPort, router))

}

//	Получаем текст из HTML
func getDataFromNode(node *html.Node) string {
	if node == nil {
		// fmt.Println("Получили пустую ноду")
		return ""
	}
	// fmt.Println("Go deep")
	var resultString string

	d := node

	for {
		if d.Type == 1 {
			resultString += d.Data
		}

		if d.FirstChild != nil {
			resultString += getDataFromNode(d.FirstChild)
		}

		if d.NextSibling == nil {
			break
		}
		// fmt.Println(strings.Repeat("-", i))
		d = d.NextSibling

	}

	return resultString
}

//	Вычищаем мусор из HTML
func getText(articleUrl string) string {
	resp, err := soup.Get(articleUrl)
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)
	sections := doc.FindAll("section", "data-page-area", "article-body")
	var text string
	for _, section := range sections {
		datas := section.FindAll("p")
		for _, data := range datas {
			if 	(data.Attrs()["class"] == "special__button" ||
				data.Attrs()["class"] == "listing__excerpt") ||
				len(data.Attrs()["data-line-id"]) < 1 {
				continue
			}
			text += getDataFromNode(data.Pointer.FirstChild)
			//fmt.Printf("%+v", link.Attrs())

		}
	}
	fmt.Println("выход из геттекст")
	fmt.Println(text)
	return text
}

//	Отрабатываем Project-Syndicate
func ps() string {
	var returnText string
	resp, err := soup.Get("https://www.project-syndicate.org/")
	if err != nil {
	}
	doc := soup.HTMLParse(resp)
	divs := doc.FindAll("article", "class", "section__article-preview")
	count := 0
	for _, div := range divs {
		datas := div.FindAll("a")
		for _, data := range datas {
			if strings.Contains(data.Attrs()["href"], "commentary") {
				fmt.Println("=====")
				fmt.Println(data.Text(), "\n", data.Attrs()["href"])
				count++
				returnText = getText("https://project-syndicate.org" + data.Attrs()["href"])
			}
		}
	}
	fmt.Println("Всего: ", count)
	return returnText
}

//	Отрабатываем Project-Syndicate для новостей
func psForNews() []newsStruct {
	//	создаем переменную для возврашения
	var returnNews []newsStruct
	resp, err := soup.Get("https://www.project-syndicate.org/")
	if err != nil {
	}
	doc := soup.HTMLParse(resp)
	divs := doc.FindAll("article", "class", "section__article-preview")
	count := 0
	for _, div := range divs {
		datas := div.FindAll("a")
		for _, data := range datas {
			if strings.Contains(data.Attrs()["href"], "commentary") {
				news :=newsStruct {
					data.Attrs()["href"],
					"",
					"",
					getText("https://project-syndicate.org" + data.Attrs()["href"]),
					data.Text()}
				returnNews = append(returnNews, news)
				count++

			}
		}
	}
	return returnNews
}

////	функция отправки методом пост
//func postToDB(
//	queryString string,	//	запрос
//	funcString string,	//	функция
//	indexString string, //	индекс
//	serverString string, //	сервер
//	) []byte {
//	// переводим строку запроса в последовательность байтов
//	var query = []byte(queryString)
//	if ( len(indexString) == 0 ) {
//		indexString = "data/test00"
//	}
//	if ( len(serverString) == 0 ) {
//		serverString = "localhost:9200"
//	}
//	req, err := http.NewRequest("POST",
//		"http://" + serverString + "/" + indexString + "/" + funcString,
//		bytes.NewBuffer(query),
//	)
//	req.Header.Set("X-Custom-Header", "myvalue")
//	req.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//
//	if err != nil {
//		// handle error
//		panic(err)
//	}
//	defer resp.Body.Close()
//
//	body, _ := ioutil.ReadAll(resp.Body)
//	fmt.Println("Resp: \n" +
//		"==================================================\n", string(body))
//	return body
//}

//	ответ на посты
func postsFunc(w http.ResponseWriter, req *http.Request) {
	fmt.Println("postsFunc")
	params := mux.Vars(req)
	if params["id"] == "" {
		fmt.Println("isEmplty")
		//	Session Opened
		session, err := mgo.Dial("192.168.1.50")
		//	Error
		if err != nil {
			log.Fatalln(err)
		}
		//	Session Closed
		defer session.Close()

		connectionA := session.DB("test").C("qz0.ru")
		posts := []postStruct{}
		err = connectionA.Find(bson.M{}).Sort("-timestamp").All(&posts)
		//	Error
		if err != nil {
			log.Fatalln(err)
		}
		//	Return Resp
		json.NewEncoder(w).Encode(posts)



	} else {
		fmt.Println("isEmplty")
		json.NewEncoder(w).Encode("{test: 'test'}")
		//	DB
		dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
		db, err := sqlx.Connect("postgres", dbinfo)
		//checkErr(err)
		if err != nil {
			log.Fatalln(err)
		}
		defer db.Close()

		post := []postStruct{}
		err = db.Select(&post, "SELECT text, number FROM test WHERE number=11")
		fmt.Print("Is :", err)
		fmt.Print(post, len(post))


		//if err != nil {
		//	fmt.Println(err)
		//	return
		//} else {
		//	fmt.Print(err)
		//	fmt.Print(test0[0].Text, "123", test0[0].Number)
		//}

		//rows, err := db.Query("SELECT text, number FROM test WHERE number=11")
		//if rows == nil {
		//	fmt.Println("Rows is Nil!")
		//}
		//for rows.Next() {
		//	nameOfColumns, _ := rows.Columns()
		//	for nameOfColumn := range nameOfColumns {
		//		fmt.Print(nameOfColumns[nameOfColumn], "\t\t\t|\t")
		//	}
		//	fmt.Print("\n========================================\n")
		//	var row1 string
		//	var row2 int
		//	if err := rows.Scan(&row1, &row2); err != nil {
		//		log.Fatal(err)
		//	}
		//	fmt.Printf("%s\t|\t%d\t\n", row1, row2)
		//
		//}

		//fmt.Print("\nError: ", rows.Err())
	}


}

//	ответ на посты
func addPostFunc(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)

	fmt.Println("До реквеста")
	if req.Method == "POST" {

		var post postStructForSend
		fmt.Println("Метод ПОСТ")
		post.Picture = req.FormValue("_id")
		bs, err := ioutil.ReadAll(req.Body)
		//	перегоняем байтовый массив в структуру
		json.Unmarshal(bs, &post)
		//	Session Opened
		session, err := mgo.Dial("192.168.1.50")
		//	Error
		if err != nil {
			log.Fatalln(err)
		}
		//	Session Closed
		defer session.Close()

		connectionA := session.DB("test").C("qz0.ru")
		err = connectionA.Insert(post)
		//	Error
		if err != nil {
			log.Fatalln(err)
		}
	}
}
//	ответ на посты
func uploadFunc(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)

	fmt.Println("До реквеста")
	if req.Method == "POST" {
		fmt.Println(req)
		fmt.Println("Метод ПОСТ")
		err := req.ParseMultipartForm(200000)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		fmt.Println(req)
		// они все тут
		formdata := req.MultipartForm // ok, no problem so far, read the Form data


		//get the *fileheaders
		files := formdata.File["myFile"] // grab the filenames
		//var fileHeader multipart.FileHeader

		for key, _ := range files {
			fmt.Printf("file: %s", key)
			//file, _ := f[key].Open()
			fmt.Println("Success")
			//defer file.Close()

		}
		fmt.Println(len(files))

		file, handler, err := req.FormFile("myFile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)


		//file, err := files[0].Open()

		fmt.Println("1")
		fmt.Println(err)
		//defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

	}

}

//	ответ на новости
func newsFunc(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		fmt.Println("newsFunc")

		//	Session Closed
		news:=[]newsStruct{}
		//	Error
		news = psForNews()
		//	Return Resp
		json.NewEncoder(w).Encode(news)

		//fmt.Println(len(news))



	} else {
		fmt.Println("isEmplty")
		json.NewEncoder(w).Encode("{test: 'test'}")
		//	DB
		dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
		db, err := sqlx.Connect("postgres", dbinfo)
		//checkErr(err)
		if err != nil {
			log.Fatalln(err)
		}
		defer db.Close()

		post := []postStruct{}
		err = db.Select(&post, "SELECT text, number FROM test WHERE number=11")
		fmt.Print("Is :", err)
		fmt.Print(post, len(post))


		//if err != nil {
		//	fmt.Println(err)
		//	return
		//} else {
		//	fmt.Print(err)
		//	fmt.Print(test0[0].Text, "123", test0[0].Number)
		//}

		//rows, err := db.Query("SELECT text, number FROM test WHERE number=11")
		//if rows == nil {
		//	fmt.Println("Rows is Nil!")
		//}
		//for rows.Next() {
		//	nameOfColumns, _ := rows.Columns()
		//	for nameOfColumn := range nameOfColumns {
		//		fmt.Print(nameOfColumns[nameOfColumn], "\t\t\t|\t")
		//	}
		//	fmt.Print("\n========================================\n")
		//	var row1 string
		//	var row2 int
		//	if err := rows.Scan(&row1, &row2); err != nil {
		//		log.Fatal(err)
		//	}
		//	fmt.Printf("%s\t|\t%d\t\n", row1, row2)
		//
		//}

		//fmt.Print("\nError: ", rows.Err())
	}


}


func testFunc(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	if params["id"] == "" {
		fmt.Println("isEmplty")
		json.NewEncoder(w).Encode("isEmpty")
	} else {
		fmt.Println("isEmplty")
		json.NewEncoder(w).Encode("{test: 'test'}")
		//	DB
		dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
		db, err := sqlx.Connect("postgres", dbinfo)
		//checkErr(err)
		if err != nil {
			log.Fatalln(err)
		}
		defer db.Close()

		test0 := []testStruct{}
		err = db.Select(&test0, "SELECT text, number FROM test WHERE number=11")
		fmt.Print("Is :", err)
		fmt.Print(test0, len(test0))


		//if err != nil {
		//	fmt.Println(err)
		//	return
		//} else {
		//	fmt.Print(err)
		//	fmt.Print(test0[0].Text, "123", test0[0].Number)
		//}

		//rows, err := db.Query("SELECT text, number FROM test WHERE number=11")
		//if rows == nil {
		//	fmt.Println("Rows is Nil!")
		//}
		//for rows.Next() {
		//	nameOfColumns, _ := rows.Columns()
		//	for nameOfColumn := range nameOfColumns {
		//		fmt.Print(nameOfColumns[nameOfColumn], "\t\t\t|\t")
		//	}
		//	fmt.Print("\n========================================\n")
		//	var row1 string
		//	var row2 int
		//	if err := rows.Scan(&row1, &row2); err != nil {
		//		log.Fatal(err)
		//	}
		//	fmt.Printf("%s\t|\t%d\t\n", row1, row2)
		//
		//}

		//fmt.Print("\nError: ", rows.Err())
	}


}

//	Работа с файлами в облаке
func cloudFunc(w http.ResponseWriter, req *http.Request) {
	//	Получаем параметры
	params := mux.Vars(req)

	//	Разные вентили в зависимости от того, что передали
	switch commandTo := params["id"]; commandTo {

	//	Ничего не передали
	case "":
		//	Вывели пустоту
		fmt.Println("isEmplty")

		//	Выводим содержимое рабочей директории на экран
	case "ls":
		//	Выводим сообщение
		json.NewEncoder(w).Encode("Найдены файлы:")

		//	Читаем содержимое каталога
		ls, err := os.Open(".")
		if err != nil {
			return
		}

		//	Отложенно закрываем
		defer ls.Close()

		//	Читаем содержимое каталога
		fileInfos, err := ls.Readdir(-1)
		if err != nil {
			return
		}

		//	Выводим на печать
		for _, fi := range fileInfos {

			json.NewEncoder(w).Encode(fi.Name())
		}

		//	Выводим текст на экран + хэш
	case "get":
		//	Выводим сообщение
		json.NewEncoder(w).Encode("Найдены файлы:")

		//	Читаем содержимое каталога
		myFile, err := os.Open(params["file"])
		if err != nil {
			return
		}

		//	Отложенно закрываем
		defer myFile.Close()

		// Получить размер файла
		stat, err := myFile.Stat()
		if err != nil {
			return
		}

		// Чтение файла
		bs := make([]byte, stat.Size())
		_, err = myFile.Read(bs)
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(string(bs)[:len(string(bs))-1])

		//	Отсылаем по почте
	case "email":
		//	Выводим сообщение
		json.NewEncoder(w).Encode("Найдены файлы:")

		//	Читаем содержимое файла
		myFile, err := os.Open(params["file"])
		if err != nil {
			return
		}

		//	Отложенно закрываем
		defer myFile.Close()

		// Получить размер файла
		stat, err := myFile.Stat()
		if err != nil {
			return
		}

		// Чтение файла
		bs := make([]byte, stat.Size())
		_, err = myFile.Read(bs)
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(string(bs)[:len(string(bs))-1])

		e := email.NewEmail()
		e.From = "Test <mail@qz0.ru>"
		e.To = []string{"me@qz0.ru"}
		//e.Bcc = []string{"test_bcc@example.com"}
		//e.Cc = []string{"test_cc@example.com"}
		e.Subject = "Awesome Subject"
		//e.Text = []byte("Text Body is, of course, supported!")
		e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
		e.AttachFile("test")
		e.AttachFile("stdout")
		err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "mail@qz0.ru", "Qz123456!1", "smtp.gmail.com"))
		if err != nil {
			return
		}

		//	Выводим текст на экран + хэш
	case "show":
		//	Выводим сообщение
		json.NewEncoder(w).Encode("Найдены файлы:")

		//	Читаем содержимое каталога
		myFile, err := os.Open(params["file"])
		if err != nil {
			return
		}

		//	Отложенно закрываем
		defer myFile.Close()

		// Получить размер файла
		stat, err := myFile.Stat()
		if err != nil {
			return
		}

		// Чтение файла
		bs := make([]byte, stat.Size())
		_, err = myFile.Read(bs)
		if err != nil {
			return
		}
		json.NewEncoder(w).Encode(string(bs)[:len(string(bs))-1])

		//	Выводим форму
	case "add":

		var AddForm = `
			<form name="fileForm" method="POST"  enctype="multipart/form-data" 
			action="/cloud/upload">
			Пользователь: <select name="user">`
		for _, myCloudUser := range myCloudUsers {
			AddForm += "<option>" + myCloudUser.UserName + "</option>"
		}

		AddForm +=   	`
			</select>
			Добавить файл: <input type="file" name="uploadFile">
			<input type="submit" value="Добавить">
			</form>
			`
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, AddForm)

		//	Выводим форму
	case "file":

		//	Читаем содержимое файла
		myFile, err := os.Open(params["file"])
		if err != nil {
			return
		}

		//	Отложенно закрываем
		defer myFile.Close()

		// Получить размер файла
		stat, err := myFile.Stat()
		if err != nil {
			return
		}

		// Чтение файла
		bs := make([]byte, stat.Size())
		_, err = myFile.Read(bs)
		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "multipart/form-data")
		fmt.Fprint(w, string(bs))
		//	Выводим сообщение

		//	Выводим текст на экран + хэш
	case "sha1":
		//	Выводим сообщение
		json.NewEncoder(w).Encode("SHA1:")

		//	SHA1
		h := sha1.New()
		h.Write([]byte(params["file"]))
		bs := h.Sum([]byte{})
		json.NewEncoder(w).Encode(bs)

		//	Выводим текст на экран + хэш
	case "upload":
		//	Выводим сообщение	fmt.Println("До реквеста")
		if req.Method == "POST" {
			fmt.Println(req)
			fmt.Println("Метод ПОСТ")
			err := req.ParseMultipartForm(200000)
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			fmt.Println(req)
			// они все тут
			formdata := req.MultipartForm // ok, no problem so far, read the Form data

			//get the *fileheaders
			files := formdata.File["fileForm"] // grab the filenames
			//var fileHeader multipart.FileHeader

			for key, _ := range files {
				fmt.Printf("file: %s", key)
				//file, _ := f[key].Open()
				fmt.Println("Success")
				//defer file.Close()

			}
			fmt.Println(len(files))

			file, handler, err := req.FormFile("fileForm")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			fmt.Fprintf(w, "%v", handler.Header)
			f, err := os.OpenFile("./data/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)

			//file, err := files[0].Open()

			fmt.Println("1")
			fmt.Println(err)
			//defer file.Close()
			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
		}

		//	Файл надо удалить
	case "deletefile":
		//	грохаем файл
		fmt.Println("Запрос на удаление")
		err := os.Remove("./data/" + params["file"])
		if err != nil {
			fmt.Println(" - неудача")
			fmt.Println(err)
			return
		}

		//	Файл надо переслать пользователю
	case "sendfile":
		//	грохаем файл
		fmt.Println("Запрос на пересылку")
		//	Отправка письма
		cloudSendFileByMail(false, params["file"], params["user"])

		//	Поведение по умолчанию
	default:
		//	Если не пустой - вывели что передалось
		fmt.Println(commandTo)
		json.NewEncoder(w).Encode(commandTo)

	}


}

//	Инициализация данных из командной строки
func initPromptArgs() {
	//	Представляемся
	fmt.Println("Вошли в initPromptArgs")
	//	Проверяем сколько аргументов в командной строке
	if len(os.Args) > 1 {
		//	Аргумента два - первый : порт подкючения
		listeningPort = os.Args[1]
	} else {
		//	Аргумент один - порт подключения берем по умолчанию
		listeningPort = LISTENING_PORT
	}
	fmt.Println(listeningPort)
}

//	ответ на посты
func cloudUploadFile(w http.ResponseWriter, req *http.Request) {
	//params := mux.Vars(req)

	fmt.Println("До реквеста")
	if req.Method == "POST" {
		fmt.Println(req)
		fmt.Println("Метод ПОСТ")
		err := req.ParseMultipartForm(200000)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		fmt.Println(req)
		// они все тут
		formdata := req.MultipartForm // ok, no problem so far, read the Form data


		user := formdata.Value["user"]
		//get the *fileheaders
		files := formdata.File["fileForm"] // grab the filenames
		//var fileHeader multipart.FileHeader

		for key, _ := range files {
			fmt.Printf("file: %s", key)
			//file, _ := f[key].Open()
			fmt.Println("Success")
			//defer file.Close()

		}
		fmt.Println(len(files))

		file, handler, err := req.FormFile("uploadFile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Fprintf(w, "%v", handler.Header)
		f, err := os.OpenFile("./data/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)


		//file, err := files[0].Open()

		fmt.Println(user)

		//******************************
		//	Отправка письма
		cloudSendFileByMail(true, user[0], handler.Filename)


		//******************************

		fmt.Println(err)
		//defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

	}

}

//	Пересылка файла нужному пользователю
func cloudSendFileByMail(typeOfMessage bool, username string, filename string) {
	//	Запрос на пересылку
	e := email.NewEmail()

	//	Проверяем кому
	if typeOfMessage {

		//	Это администрации
		if mySettings.EmailShadow != "" {
			e.Bcc = []string{mySettings.EmailShadow}
		}
		if mySettings.EmailCopy != "" {
			e.Cc = []string{mySettings.EmailCopy}
		}

		e.To = []string{mySettings.EmailTo}
		e.Subject = mySettings.EmailSubjectAdm

		e.HTML = []byte(mySettings.EmailTextBeforeUsernameAdm + username + mySettings.EmailTextAfterUsernameAdm + `
		<br><a href="http://` + mySettings.ListeningIp + `:`+ mySettings.ListeningPort +
			`/cloud/sendfile/` + username  + `/` + filename + `">Разрешить</a>
		<br><a href="http://` + mySettings.ListeningIp +`:`+ mySettings.ListeningPort +
			`/cloud/deletefile/` + filename + `">Отклонить</a>`)

	} else {

		//	Это оконечка пользователя
		//	Получаем мыло
		for _, myCloudUser := range myCloudUsers {
			if myCloudUser.UserName == username {
				e.To = []string{myCloudUser.UserEmail}
			}
		}
		//e.To = []string{myCloudUsers[username]}

		e.Subject = mySettings.EmailSubjectUser
		e.HTML = []byte(mySettings.EmailTextBeforeUsernameUser + username + mySettings.EmailTextAfterUsernameUser)

	}

	//	Сама отправка
	e.From = mySettings.EmailFromName + " <" + mySettings.EmailFromMail + ">"
	//e.Text = []byte("Text Body is, of course, supported!")
	e.AttachFile("./data/"+filename)
	err := e.Send(mySettings.EmailSmtpServerIp + ":" + mySettings.EmailSmtpServerPort,
		smtp.PlainAuth("", mySettings.EmailSmtpLogin, mySettings.EmailSmtpPassword,
			mySettings.EmailSmtpServerIp))
	if err != nil {

		fmt.Println(e)
		fmt.Println(err)
		return
	}
}

//	программа загрузки настроек
func loadSettings() {

	//	Читаем содержимое файла настроек
	myFile, err := os.Open(SETTINGS_FILE_NAME)
	if err != nil {

		//	Файла нет - создаем
		myFile, err := os.OpenFile("./"+SETTINGS_FILE_NAME, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer myFile.Close()

		//	Объявляем для примера один параметр
		mySettings.ListeningPort = LISTENING_PORT
		//	Маршализируем
		buf, err := json.Marshal(mySettings)
		//	Копируем структурку в файлик
		myFile.Write(buf)

		return
	}
	//	Отложенно закрываем
	defer myFile.Close()

	// Получить размер файла
	stat, err := myFile.Stat()
	if err != nil {
		return
	}

	// Чтение файла
	buf := make([]byte, stat.Size())
	_, err = myFile.Read(buf)
	if err != nil {
		return
	}

	//	Маршалим прочитанное в структуру
	json.Unmarshal(buf, &mySettings)

	//	Выводим сообщение
	fmt.Println("Процесс инициализации завершен")
}

//	программа загрузки настроек
func loadCloudUsers() {
	//	А может тут и ЛДАП будет когда-нить

	//	Читаем содержимое файла пользователей
	myFile2, err := os.Open(USERS_FILE_NAME)
	if err != nil {

		//	Файла нет - создаем
		myFile, err := os.OpenFile("./"+USERS_FILE_NAME, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer myFile2.Close()
		//	Копируем структурку в файлик
		myFile.Write([]byte(`{"UserName":"Иванов","UserEmail":"me@qz0.ru"}, 
			{"UserName":"Петров","UserEmail":"me@qz0.ru"}`))

		return
	}
	//	Отложенно закрываем
	defer myFile2.Close()

	// Получить размер файла
	stat, err := myFile2.Stat()
	if err != nil {
		return
	}

	// Чтение файла
	buf := make([]byte, stat.Size())
	_, err = myFile2.Read(buf)
	if err != nil {
		return
	}

	//	Маршалим прочитанное в структуру
	json.Unmarshal(buf, &myCloudUsers)

	//	Выводим сообщение
	fmt.Println("Процесс подгрузки списка пользователей завершен")
	fmt.Println(myCloudUsers)
}


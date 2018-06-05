package main

import (
	"database/sql"
	"fmt"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "login.html")
		return
	}

	username := req.FormValue("LoginForm[username]")
	password := req.FormValue("LoginForm[password]")


	var databaseUsername string
	var databasePassword string

	err := db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)

	if err != nil {
		http.Redirect(res, req, "/", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/", 301)
		return
	} else {
		http.Redirect(res, req, "/socks/search", 302)
		return
	}

}

func registerPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "register.html")
		return
	}

	username := req.FormValue("LoginForm[username]")
	password := req.FormValue("LoginForm[password]")

	var user string

	err := db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		} else {
			fmt.Println("Success!")
			http.Redirect(res, req, "/", 302)
			return
		}

		// res.Write([]byte("User created!"))
		// return
	case err != nil:
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
	}
}


func socksPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "search.html")
		return
	}

}

func profilePage(res http.ResponseWriter, req *http.Request) {

	if req.Method != "POST" {
		http.ServeFile(res, req, "profile.html")
		return
	}

}

func addfundsPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "add-funds.html")
		return
	}

}

func historyPage(res http.ResponseWriter, req *http.Request) {
	
	if req.Method != "POST" {
		http.ServeFile(res, req, "history.html")
		return
	}

}

func historydepositPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "history-deposit.html")
		return
	}

}

func main() {
	db, err = sql.Open("mysql", "root:rootroot@/luxsocks")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/socks/search", socksPage)
	http.HandleFunc("/profile/edit", profilePage)
	http.HandleFunc("/add-funds", addfundsPage)
	http.HandleFunc("/home/register", registerPage)
	http.HandleFunc("/socks/history", historyPage)
	http.HandleFunc("/history-deposit", historydepositPage)
	// http.HandleFunc("/login", loginPage)
	http.HandleFunc("/", loginPage)
	
	fmt.Println("Listening on 127.0.0.1:8000")
	http.ListenAndServe(":8000", nil)
}
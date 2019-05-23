package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"os"
	"strconv"
	"log"
	"net/http"
)
import _ "github.com/go-sql-driver/mysql"

type Publisher struct {
	Id   int    `json:"publisher_id"`
	Name string `json:"name"`
}
type Publishers []Publisher

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
}

func dbGetPublisher(publisher_id int) (Publisher, error) {
	fmt.Println("dbGetPublisher")
	var publisher Publisher
	db := dbConn()
	/*results, err := db.Query(`SELECT publisher_id, name FROM publishers WHERE publisher_id=?`,publisher_id)
	if err != nil {
		panic(err.Error())
	}
	//var publishers Publishers
	for results.Next() {
		err = results.Scan(&publisher.Id, &publisher.Name)
		if err != nil {
			panic(err.Error())
		}
		//publishers = append(publishers, publisher)
	}*/
	results := db.QueryRow(`SELECT publisher_id, name FROM publishers WHERE publisher_id=?`, publisher_id)
	err := results.Scan(&publisher.Id, &publisher.Name)
	if err != nil {
		//panic(err.Error())
		fmt.Println("error:", err.Error())
		return publisher, err
	}
	//publishers = append(publishers, publisher)

	defer db.Close()
	//json.NewEncoder(w).Encode(publishers)
	fmt.Println("returned publisher: ", publisher)
	return publisher, nil
}
func apiGetPublishers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("api.GetPublishers")
	var publisher Publisher
	db := dbConn()
	results, err := db.Query("SELECT publisher_id, name FROM publishers ORDER BY publisher_id DESC")
	if err != nil {
		panic(err.Error())
	}
	var publishers Publishers
	for results.Next() {
		err = results.Scan(&publisher.Id, &publisher.Name)
		if err != nil {
			panic(err.Error())
		}
		publishers = append(publishers, publisher)
	}
	defer db.Close()
	json.NewEncoder(w).Encode(publishers)
	fmt.Println("returned publishers: ", publishers)
}
func apiGetPublisher(w http.ResponseWriter, r *http.Request) {
	fmt.Println("api.apiGetPublisher")
	vars := mux.Vars(r)
	fmt.Println("vars.p=", vars["publisherId"])
	publisherId, _ := strconv.Atoi(vars["publisherId"])
	fmt.Println("pubid=", publisherId)

	publisher, err := dbGetPublisher(publisherId)
	if err != nil {
		//** error
		fmt.Println("api.apiGetPublisher: error", err)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity

		/*m := MyError{Error:err.Error()}
		if err := json.NewEncoder(w).Encode(m); err != nil {
			panic(err)
		}*/

		//fmt.Fprintf(w, `{"result":"", "error": "%s"}`, err)
		fmt.Fprintf(w, `{"error": %q}`, err)
		fmt.Println("ret")
		return
	}
	fmt.Println("api.apiGetPublisher: ok", publisher)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	/*if err := json.NewEncoder(w).Encode(publisher); err != nil {
		panic(err)
	}*/
	js, _ := json.Marshal(publisher)
	fmt.Fprintf(w, `{"publisher":%s}`, js)
	fmt.Println("ret publlisher: ", publisher)

}
func postPublisher(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Post Publisher")
	var p Publisher

	//V1: using json.Decoder
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1048576))
	if err := decoder.Decode(&p); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("ret")
		return
	}
	//V1: using json.Unmarshall
	/*body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &p); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		fmt.Println("ret")
		return
	}
*/
	//t := RepoCreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(p); err != nil {
		panic(err)
	}
	fmt.Println("in: ", p)
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	fmt.Println("DB_HOST:", dbHost)
	fmt.Println("DB_PORT:", dbPort)
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}
func dbCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("creating a table..")
	db := dbConn()
	var err error
	_, err = db.Exec(`CREATE TABLE publishers(publisher_id INT NOT NULL AUTO_INCREMENT,
		 											 name varchar(255) NOT NULL,
		 											 PRIMARY KEY (publisher_id))
						`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO publishers	(publisher_id, name) VALUES
			   											(1, 'Addison-Wesley'),
			    										(2, 'Manning Publications');
						`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE authors(	author_id INT NOT NULL AUTO_INCREMENT,
														name varchar(255) NOT NULL,
														book_id INT,
														PRIMARY KEY (author_id))
						`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO authors	(book_id, name) VALUES
			   											(1, 'Erich Gamma'),
			   											(1, 'Richard Helm'),
			   											(1, 'Ralph Johnson'),
			    										(1, 'John Vlissides'),
			    										(2, 'Sau Sheong Chang ');
						`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE books ( book_id INT NOT NULL AUTO_INCREMENT,
	 							isbn varchar(13) DEFAULT NULL,
	 							publisher_id INT NOT NULL,
	 							title varchar(255) NOT NULL,
	 							PRIMARY KEY (book_id))`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`INSERT INTO books(	book_id, publisher_id, title) VALUES
		    										(1, 1, 'Design Patterns: Elements of Reusable Object-Oriented Software' ),
		    										(2, 2, 'Go Web Programming');
		    			`)
	if err != nil {
		panic(err)
	}

	fmt.Println(".. ok")
	fmt.Fprintf(w, ".. ok")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/create_db", dbCreate)
	myRouter.HandleFunc("/publishers", apiGetPublishers).Methods("GET")
	myRouter.HandleFunc("/publisher", postPublisher).Methods("POST")
	myRouter.HandleFunc("/publisher/{publisherId}", apiGetPublisher).Methods("GET")
	log.Fatal(http.ListenAndServe(":8082", myRouter))
}

func main() {

	handleRequests()

	fmt.Println("asd")
	dbDriver := "mysql"
	dbUser := "user"
	dbPass := "password2"
	dbName := "db"
	dbPath :=os.Getenv("FF_DB_HOST")
	fmt.Println("dbhost="+dbPath)
	//db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbPath+":3308)/"+dbName)
//
	//db, err := sql.Open("mysql", "user:password@/dbname")
	fmt.Println(db)
	fmt.Println("***")
	fmt.Println(err)

	/*_,err = db.Exec("CREATE DATABASE "+dbName)
	if err != nil {
	   panic(err)
	}*/
	fmt.Println("connecting to db..")
	_, err = db.Exec("USE " + dbName)
	if err != nil {
		panic(err)
	}

}



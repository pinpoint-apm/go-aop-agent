/*
 * Copyright 2021 NAVER Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"naver/app"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	pinpoint "github.com/pinpoint-apm/go-aop-agent/middleware/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("tmpl")).ServeHTTP(w, r)
}

type Book struct {
	name string
}

type Reader interface {
	Read() (string, error)
}

func (b Book) Read() {
	fmt.Printf("Reading %v\n", b.name)
}

var globalv string = "globle variable"

const LENGTH int = 10

func TestUserFuncHandler(w http.ResponseWriter, r *http.Request) {
	res0 :=app.TestUserFunc(r.Context(),35, 42)
	fmt.Fprint(w,res0)
}

func TestCommonFuncHandler(w http.ResponseWriter, r *http.Request) {

	res1 := app.TestComFunc(r.Context(), 6.66) //float
	fmt.Fprintf(w, "type of %T tested!\n", res1)
	res2 := app.TestComFunc(r.Context(), "string") //string
	fmt.Fprintf(w, "type of %T tested!\n", res2)
	res3 := app.TestComFunc(r.Context(), true) //bool
	fmt.Fprintf(w, "type of %T tested!\n", res3)

	var array1 = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0} //array
	res4 := app.TestComFunc(r.Context(), array1)
	fmt.Fprintf(w, "type of %T tested!\n", res4)

	var slice1 = make([]int, 3, 5) //slice
	res5 := app.TestComFunc(r.Context(), slice1)
	fmt.Fprintf(w, "type of %T tested!\n", res5)

	var book = Book{name: "test"} //struct interface
	res6 := app.TestComFunc(r.Context(), book)
	fmt.Fprintf(w, "type of %T tested!\n", res6)

	var x int = 20
	var p *int = &x //printer
	res7 := app.TestComFunc(r.Context(), p)
	fmt.Fprintf(w, "type of %T tested!\n", res7)

	res8 := app.TestComFunc(r.Context(), func(i int) { fmt.Sprint("Hello!") }) //func
	fmt.Fprintf(w, "type of %T tested!\n", res8)

	ccmap := map[string]string{"France": "Paris", "Italy": "Rome", "Japan": "Tokyo", "India": "New delhi"} //map
	res9 := app.TestComFunc(r.Context(), ccmap)
	fmt.Fprintf(w, "type of %T tested!\n", res9)

	ch := make(chan int) //channel
	res10 := app.TestComFunc(r.Context(), ch)
	fmt.Fprintf(w, "type of %T tested!\n", res10)
	close(ch)

	dir := "data"
	file, err := ioutil.TempFile(dir, "Hello")
	if err != nil {
		log.Fatal(err)
	}
	res11 := app.TestComFunc(r.Context(), file) //file
	fmt.Fprintf(w, "type of %T tested!\n", res11)
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	res12 := app.TestComFunc(r.Context(), nil) //nil
	fmt.Fprintf(w, "type of %T tested!\n", res12)

	res13 := app.TestComFunc(r.Context(), app.CONSTANT) //constant
	fmt.Fprintf(w, "type of %T tested!\n", res13)

	res14 := app.TestComFunc(r.Context(), globalv) //globle
	fmt.Fprintf(w, "type of %T tested!\n", res14)

	res15 := app.TestComFunc(r.Context(), res14) //local
	fmt.Fprintf(w, "type of %T tested!\n", res15)

	//private tested in testreturn.go

	res17 := app.TestComFunc(r.Context(), app.PublicV) //public
	fmt.Fprintf(w, "type of %T tested!\n", res17)

	res18 := app.TestReturn(r.Context()) //return private,public,constant,func,globle,local
	res18(3)
	fmt.Fprintf(w, "type of %v tested!\n", res18)
}

func TestHighFuncHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(reflect.ValueOf(app.Add).Type())
	 res19 :=app.TestHighFunc(r.Context(), 5, 8, app.Add)
	 fmt.Fprint(w, res19)

	 res20 :=app.TestHighFunc(r.Context(), 5, 8, app.Mul)
	 fmt.Fprint(w, res20) 
	 
	 res21 :=app.Add(r.Context(),12, 15)
	 fmt.Fprint(w, res21)

	 res22 :=app.Mul(r.Context(),12, 15)
	 fmt.Fprint(w, res22)
}

func TestInheritFuncHandler(w http.ResponseWriter, r *http.Request) {
  var s app.Student = app.Student{app.Person{Name:"Tom\n", Sex:"男\n", Age:18}}
  res23, res24, res25 := s.TestInheritFunc(r.Context())
  fmt.Fprint(w, res23, res24, res25)
}

func TestLambdaFuncHandler(w http.ResponseWriter, r *http.Request) {
  res26, res27:= app.TestLambdaFunc(r.Context())
  fmt.Fprint(w, res26, res27)
}

func TestDecoratorFuncHandler(w http.ResponseWriter, r *http.Request) {
  d := app.Dinner{app.Cooker{}, app.Rice{}, app.Water{}}
  res28, res29, res30, res31, res32 :=d.TestDecoratorFunc(r.Context())
  fmt.Fprint(w, res28, res29, res30, res31, res32)

  e := app.Rice{}
  res33 := e.WashRice(r.Context())
  fmt.Fprint(w, res33)

  a :=app.Water{}
  res34 := a.AddWater(r.Context())
  fmt.Fprintln(w, res34)
}

func TestAbstractFuncHandler(w http.ResponseWriter, r *http.Request) {
  b :=app.Book{"RED"}
  res35 :=b.TestAbstractFunc(r.Context())
  fmt.Fprint(w, res35)
}

func TestGeneratorFuncHandler(w http.ResponseWriter, r *http.Request) {
  var a chan int
  a = make (chan int)
  go app.TestGeneratorFunc (r.Context(), a)

  for i := range a{
	  fmt.Fprint(w, i)
  }
}

func TestRecursionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, app.TestRecursion(r.Context(), 3))
}

func TestCallRemoteHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	remote := r.URL.Query()["remote"]
	// fmt.Println(remote)
	req, err := http.NewRequestWithContext(r.Context(), "GET", remote[0], nil)
	if err != nil {
		fmt.Fprint(w, fmt.Sprint(err))
	} else {
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprint(w, fmt.Sprint(err))
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				fmt.Fprint(w, fmt.Sprint(err))
			} else {
				fmt.Fprint(w, fmt.Sprint(string(body)))
			}
		}
	}
}

func TestExceptionHandler(w http.ResponseWriter, r *http.Request) {
	result, err := app.TestException(r.Context(), 9, 0)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, result)
	}
}

func TestExpInRecursionHandler(w http.ResponseWriter, r *http.Request) {
	result, err := app.TestExpInRecursion(r.Context(), 3)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, result)
	}
}

var collection *mongo.Collection

func init_mongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo-dev:27017")); err == nil {
		collection = client.Database("testing").Collection("numbers")
	}
}

func TestMongoHandler(w http.ResponseWriter, r *http.Request) {
	if res, errIn := collection.InsertOne(r.Context(), bson.D{{"name", "pi"}, {"value", 3.14159}}); errIn == nil {
		fmt.Fprint(w, res.InsertedID)
	}
}

var db *sql.DB

func init_mysql() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "root",
		Passwd: "root",
		Net:    "tcp",
		Addr:   "mysql-dev:3306",
		DBName: "Test",
		Params: map[string]string{"allowNativePasswords": "true"},
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mysql Connected!")
}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func TestMysqlHandler(w http.ResponseWriter, r *http.Request) {
	alb := Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	}

	pingErr := db.PingContext(r.Context())
	if pingErr != nil {
		fmt.Fprint(w, "pingerror: %v", pingErr)
		return
	}

	result, err := db.ExecContext(r.Context(), "INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		fmt.Fprint(w, "addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		fmt.Fprint(w, "addAlbum: %v", err)
	}

	rows, err := db.QueryContext(r.Context(), "SELECT * FROM album WHERE id = ?", id)
	if err != nil {
		fmt.Fprint(w, fmt.Errorf("albumsByArtist %q: %v", alb.Artist, err))
	}
	defer rows.Close()
	var albums []Album
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			fmt.Fprint(w, fmt.Errorf("albumsByArtist %q: %v", alb.Artist, err))
		}
		albums = append(albums, alb)
	}
	fmt.Fprint(w, albums)
}

func TestRedisHandler(w http.ResponseWriter, r *http.Request) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, _ := rdb.Ping(r.Context()).Result()
	fmt.Fprint(w, pong)

	err := rdb.Set(r.Context(), "Evy", 666, 0).Err()
	if err != nil {
		fmt.Fprint(w, err)
	}

	val, err := rdb.Get(r.Context(), "Evy").Result()
	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprintf(w, "key:%v", val)

	rdb.RPush(r.Context(), "list_test", "message1", "message2", "message3", "message4", "message5")
	rdb.LSet(r.Context(), "list_test", 2, "message set")
	rdb.LRem(r.Context(), "list_test", 3, "message1")
	rdb.LLen(r.Context(), "list_test")

	rdb.BLPop(r.Context(), 1*time.Second, "list_test")
	rdb.BRPop(r.Context(), 1*time.Second, "list_test")

	res, _ := rdb.Do(r.Context(), "set", "dotest", "testdo").Result()
	fmt.Println(res)

	res2, _ := rdb.Append(r.Context(), "feekey", "_add").Result()
	fmt.Println(res2)

	datas := map[string]interface{}{
		"name": "LI LEI",
		"sex":  1,
		"age":  28,
		"tel":  123445578,
	}
	rdb.HMSet(r.Context(), "hash_test", datas)
	rdb.HMGet(r.Context(), "hash_test", "name", "sex")
	rdb.HGetAll(r.Context(), "hash_test")
	rdb.HSetNX(r.Context(), "hash_test", "id", 100)
	rdb.HDel(r.Context(), "hash_test", "age")
	rdb.SAdd(r.Context(), "set_test", "11", "22", "33", "44")
	rdb.SRem(r.Context(), "set_test", "11", "22")
	rdb.SMembers(r.Context(), "set_test")
	rdb.SInter(r.Context(), "set_a", "set_b")

	rdb.Close()
}

func main() {
	init_mongo()
	init_mysql()
	defer db.Close()
	router := mux.NewRouter()
	// add pinpoint middleware
	router.Use(pinpoint.PinpointMuxMiddleWare)

	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/test_user_func", TestUserFuncHandler)
	router.HandleFunc("/test_args_return", TestCommonFuncHandler)
	router.HandleFunc("/test_high-order_func", TestHighFuncHandler)
	router.HandleFunc("/test_inherit_func", TestInheritFuncHandler)
	router.HandleFunc("/test_lambda_func", TestLambdaFuncHandler)
	router.HandleFunc("/test_decorator_func", TestDecoratorFuncHandler)
	router.HandleFunc("/test_abstract_func", TestAbstractFuncHandler)
	router.HandleFunc("/test_generator_func", TestGeneratorFuncHandler)
	router.HandleFunc("/test_recursion", TestRecursionHandler)
	router.HandleFunc("/test_mongo", TestMongoHandler)
	router.HandleFunc("/test_mysql", TestMysqlHandler)
	router.HandleFunc("/test_redis", TestRedisHandler)
	router.HandleFunc("/call_remote", TestCallRemoteHandler)
	router.HandleFunc("/test_exception", TestExceptionHandler)
	router.HandleFunc("/test_exception_in_recursion", TestExpInRecursionHandler)

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

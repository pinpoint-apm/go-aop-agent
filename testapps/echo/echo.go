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
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo"
	pinpoint "github.com/pinpoint-apm/go-aop-agent/middleware/echo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func getUser(c echo.Context) error {
	// log.Printf("get  %d", c.Request().Context().Value("id"))
	id := c.Param("id")
	resp, _ := http.NewRequestWithContext(c.Request().Context(), "get", "http://example.com/a/ba/c", nil)

	client := &http.Client{}
	res, err := client.Do(resp)

	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	out, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))

	// read mongo

	if res, errIn := collection.InsertOne(c.Request().Context(), bson.D{{"name", "pi"}, {"value", 3.14159}}); errIn == nil {
		log.Println(res.InsertedID)
	}

	return c.String(http.StatusOK, id)
}

func init_mongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://10.34.130.78:27017")); err == nil {
		collection = client.Database("testing").Collection("numbers")

	}

}

func main() {
	//init_pinpoint()
	init_mongo()
	e := echo.New()
	e.GET("/users/:id", getUser)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Use(pinpoint.PinpointMiddleWare)

	e.Logger.Fatal(e.Start(":1323"))
}

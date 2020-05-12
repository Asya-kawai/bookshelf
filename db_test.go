// Copyright 2019 Toshiki kawai
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Note:
// The bookshelf command was forked from the below.
// https://cloud.google.com/go/getting-started/tutorial-app.
//
// Google's bookshelf program has
// * demonstrating several Google Cloud APIs
// * including App Engine
// * using Firestore
// * using Cloud Storage
//
// Change:
// * Does not use Google Cloud APIs
// * Does not including App Engine
// * use Mysql instead of Firestore
// * Does not use Cloud storage
package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
)

func testDB(t *testing.T, db BookDatabase) {
	t.Helper()

	b := &Book{
		Author:        "testy mc testface",
		Title:         fmt.Sprintf("t-%d", time.Now().Unix()),
		PublishedDate: fmt.Sprintf("%d", time.Now().Unix()),
		Description:   "desc",
	}

	id, err := db.AddBook(b)
	if err != nil {
		t.Fatal(err)
	}

	b.ID = id
	b.Description = "newdesc"
	if err := db.UpdateBook(b); err != nil {
		t.Error(err)
	}

	gotBook, err := db.GetBook(id)
	if err != nil {
		t.Error(err)
	}
	if got, want := gotBook.Description, b.Description; got != want {
		t.Errorf("Update description: got %q, want %q", got, want)
	}

	if err := db.DeleteBook(id); err != nil {
		t.Error(err)
	}

	if _, err := db.GetBook(id); err == nil {
		t.Error("want non-nil err")
	}
}

func TestMemoryDB(t *testing.T) {
	testDB(t, newMemoryDB())
}

func TestMysqlDB(t *testing.T) {
	DBHost := os.Getenv("DB_HOST")
	if DBHost == "" {
		DBHost = "localhost"
	}
	DBPort := os.Getenv("DB_PORT")
	if DBPort == "" {
		DBPort = "3306"
	}
	client, err := gorm.Open(
		"mysql",
		"user:password@("+DBHost+":"+DBPort+")/default?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Fatalf("gorm.open: %v", err)
	}
	defer client.Close()

	db, err := newDB(client)
	if err != nil {
		t.Fatalf("newDB: %v", err)
	}

	testDB(t, db)
}

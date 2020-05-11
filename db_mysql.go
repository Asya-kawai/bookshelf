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
package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// DB persists books to Mysql Database.
type DB struct {
	client *gorm.DB
}

// Ensure DB conforms to the BookDatabase interface.
var _ BookDatabase = &DB{}

// [START getting_started_bookshelf_mysql]

// newDB creates a new BookDatabase backed by Mysql.
// See the gorm package for details on creating a suitable
// https://gorm.io/ja_JP/docs/connecting_to_the_database.html
func newDB(client *gorm.DB) (*DB, error) {
	if client == nil {
		return nil, fmt.Errorf("Mysql: could not connect")
	}
	db := &DB{
		client: client,
	}
	return db, nil
}

// Close closes the database.
func (db *DB) Close() error {
	return db.client.Close()
}

// GetBook retrieves a book by its ID.
func (db *DB) GetBook(id uint) (*Book, error) {
	b := &Book{}
	err := db.client.Find(b, id).Error
	if err != nil {
		return nil, fmt.Errorf("DB: Get: %v", err)
	}
	return b, nil
}

// [END getting_started_bookshelf_mysql]

// AddBook saves a given book, assigning it a new ID.
func (db *DB) AddBook(b *Book) (uint, error) {
	creatable := db.client.NewRecord(b)
	// Primary key is not empty(this means that already created).
	if !creatable {
		return 0, fmt.Errorf("DB: Already exists %s", b.Title)
	}
	if err := db.client.Create(b).Error; err != nil {
		return 0, fmt.Errorf("DB: Create: %v", err)
	}
	creatable = db.client.NewRecord(b)
	if creatable {
		return 0, fmt.Errorf("DB: Does not Created %s", b.Title)
	}
	return b.ID, nil
}

// DeleteBook removes a given book by its ID.
func (db *DB) DeleteBook(id uint) error {
	b := &Book{
		ID: id,
	}
	if err := db.client.Delete(b).Error; err != nil {
		return fmt.Errorf("DB: Delete: %v", err)
	}
	return nil
}

// UpdateBook updates the entry for a given book.
func (db *DB) UpdateBook(b *Book) error {
	if err := db.client.Save(b).Error; err != nil {
		return fmt.Errorf("DB: Set: %v", err)
	}
	return nil
}

// ListBooks returns a list of books, ordered by title.
func (db *DB) ListBooks() ([]*Book, error) {
	books := make([]*Book, 0)
	err := db.client.Find(&books).Error
	if err != nil {
		return nil, fmt.Errorf(
			"DB: could not list books up: %v", err)
	}
	return books, nil
}

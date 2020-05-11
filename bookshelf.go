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
	"io"
	"os"

	"cloud.google.com/go/errorreporting"
	"cloud.google.com/go/storage"
)

// Book holds metadata about a book.
type Book struct {
	ID            uint   `gorm:"column:id;primary_key"`
	Title         string `gorm:"column:title"`
	Author        string `gorm:"column:author"`
	PublishedDate string `gorm:"column:published_at"`
	ImageURL      string `gorm:"column:image_url"`
	Description   string `gorm:"column:description"`
}

// BookDatabase provides thread-safe access to a database of books.
type BookDatabase interface {
	// ListBooks returns a list of books, ordered by title.
	ListBooks() ([]*Book, error)

	// GetBook retrieves a book by its ID.
	GetBook(id uint) (*Book, error)

	// AddBook saves a given book, assigning it a new ID.
	AddBook(b *Book) (id uint, err error)

	// DeleteBook removes a given book by its ID.
	DeleteBook(id uint) error

	// UpdateBook updates the entry for a given book.
	UpdateBook(b *Book) error
}

// Bookshelf holds a BookDatabase and storage info.
type Bookshelf struct {
	DB BookDatabase

	StorageBucket     *storage.BucketHandle
	StorageBucketName string

	// logWriter is used for request logging and can be overridden for tests.
	//
	// See https://cloud.google.com/logging/docs/setup/go for how to use the
	// Stackdriver logging client. Output to stdout and stderr is automaticaly
	// sent to Stackdriver when running on App Engine.
	logWriter io.Writer

	errorClient *errorreporting.Client
}

// NewBookshelf creates a new Bookshelf.
func NewBookshelf(db BookDatabase) (*Bookshelf, error) {
	b := &Bookshelf{
		logWriter: os.Stderr,
		DB:        db,
	}
	return b, nil
}

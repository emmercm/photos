package model

import (
	"path/filepath"
	"time"

	"github.com/emmercm/photos/internal/pkg/store/sqlite3"
)

// Album represents a collection of files
type Album struct {
	// Columns
	ID        uint       `json:"id"`
	Path      string     `json:"path"`
	Title     string     `json:"title"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
	// ORM
	Files []*File `json:"-" gorm:"many2many:album_files"`
}

// AlbumFromPath either creates or loads an Album from a path
func AlbumFromPath(path string) (*Album, error) {
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	a := &Album{
		Path: pathAbs,
	}

	return a, a.load()
}

func (a *Album) load() error {
	db, err := sqlite3.Connection()
	if err != nil {
		return err
	}

	if err := db.Unscoped().Where("path = ?", a.Path).FirstOrCreate(a).Error; err != nil {
		return err
	}

	return nil
}

// BeforeSave is a GORM hook to default different properties
func (a *Album) BeforeSave() error {
	if len(a.Title) == 0 && len(a.Path) > 0 {
		a.Title = filepath.Base(a.Path)
	}

	return nil
}

// Save saves the Album to DB
func (a *Album) Save() error {
	db, err := sqlite3.Connection()
	if err != nil {
		return err
	}

	a.DeletedAt = nil
	if err := db.Unscoped().Save(a).Error; err != nil {
		return err
	}

	return nil
}

// Delete soft deletes the Album in DB
func (a *Album) Delete() error {
	db, err := sqlite3.Connection()
	if err != nil {
		return err
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Model(a).Association("Files").Clear().Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(a).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// AddFile adds a File to the Album
func (a *Album) AddFile(f *File) error {
	db, err := sqlite3.Connection()
	if err != nil {
		return err
	}

	return db.Model(a).Association("Files").Append(f).Error
}

// RemoveFile removes a file from the Album
func (a *Album) RemoveFile(f *File) error {
	db, err := sqlite3.Connection()
	if err != nil {
		return err
	}

	return db.Model(a).Association("Files").Delete(f).Error
}

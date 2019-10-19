package model

import (
	"path/filepath"
	"time"

	"github.com/emmercm/photos/internal/pkg/hasher"
	"github.com/emmercm/photos/internal/pkg/store/sqlite3"
)

// File represents an individual physical file
type File struct {
	// Columns
	ID        uint       `json:"id"`
	Path      string     `json:"path"`
	FastHash  string     `json:"-"`
	SlowHash  string     `json:"-"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"-"`
	// ORM
	Albums []*Album `json:"-" gorm:"many2many:album_files"`
}

// FileFromPath either creates or loads a File from a path
func FileFromPath(path string) (*File, error) {
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f := &File{
		Path: pathAbs,
	}

	return f, f.load()
}

func (f *File) load() error {
	db, err := sqlite3.Connection()
	if err != nil {
		return err
	}

	if err := db.Unscoped().Where("path = ?", f.Path).FirstOrCreate(f).Error; err != nil {
		return err
	}

	return nil
}

// Save saves the File to DB
func (f *File) Save() error {
	db, err := sqlite3.Connection()
	if err != nil {
		return err
	}

	f.DeletedAt = nil
	if err := db.Unscoped().Save(f).Error; err != nil {
		return err
	}

	return nil
}

// Delete soft deletes the File in DB
func (f *File) Delete() error {
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

	if err := tx.Model(f).Association("Albums").Clear().Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(f).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// HasChanged determines if the underlying file has changed from what has been loaded
func (f *File) HasChanged() (bool, error) {
	fastHash, err := hasher.FastFile(f.Path)
	if err != nil {
		return false, err
	}

	if fastHash == f.FastHash {
		return false, nil
	}

	slowHash, err := hasher.SlowFile(f.Path)
	if err != nil {
		return false, err
	}

	if slowHash == f.SlowHash {
		return false, nil
	}

	return true, nil
}

// Update updates struct fields from the underlying file
func (f *File) Update() error {
	fastHash, err := hasher.FastFile(f.Path)
	if err != nil {
		return err
	}

	slowHash, err := hasher.SlowFile(f.Path)
	if err != nil {
		return err
	}

	f.FastHash = fastHash
	f.SlowHash = slowHash

	return nil
}

// GetDirnameAlbum returns the Album auto-generated from the file's dirname
func (f *File) GetDirnameAlbum() (*Album, error) {
	path := filepath.Dir(f.Path)

	return AlbumFromPath(path)
}

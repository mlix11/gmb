package database

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var defaultFolders = [3]string{"movie", "music", "picture"}

type AppDatabase interface {
	GetAll(string) ([]File, error)
}

type AppDatabaseImp struct {
	DatabaseBasePath string
}

type File struct {
	name     string
	size     int64
	filetype string
}

var filesArr []File

func NewAppDatabase(basePath string) *AppDatabaseImp {
	appdb := &AppDatabaseImp{basePath}
	_, err := appdb.setupDatabase()
	if err != nil {
		log.Printf("Database setup failed: %q \n", err)
	}
	return appdb
}

// traverse all folders to get all files
func walkfnc(path string, f os.FileInfo, err error) error {
	if err != nil {
		log.Printf(path)
		return nil
	}

	if !f.IsDir() {
		splitPath := strings.SplitAfter(path, "/")
		if splitPath[len(splitPath)-1] == f.Name() {
			// Path without suffix filesep
			path := splitPath[len(splitPath)-2]
			filesArr = append(filesArr, File{f.Name(), f.Size(), path[:len(path)-1]})
		}
	}

	return err
}

/*
* Obtain all files with given filetype, is no filetype return all files of all directories
 */
func (adi AppDatabaseImp) GetAll(filetype string) ([]File, error) {
	filesArr = filesArr[:0]

	if len(filetype) > 0 {
		files, err := ioutil.ReadDir(adi.DatabaseBasePath + filetype)
		if err != nil {
			return nil, err
		}

		filesArr = make([]File, len(files))
		for key, value := range files {
			filesArr[key] = File{value.Name(), value.Size(), filetype}
		}

	} else {
		filepath.Walk(adi.DatabaseBasePath, walkfnc)
	}
	return filesArr, nil
}

func (adi AppDatabaseImp) setupDatabase() (bool, error) {
	// check if the basic setup is provided (eg. Directories --> music, movies, pictures)
	for _, value := range defaultFolders {
		exist, err := exists(adi.DatabaseBasePath + value)
		if !exist {
			err = os.Mkdir(adi.DatabaseBasePath+value, 0777)
		}
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

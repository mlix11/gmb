package database

import (
	"io/ioutil"
	"log"
	"os"
)

var defaultFolders = [3]string{"movie", "music", "pictures"}

type AppDatabase interface {
	GetAll() ([]os.FileInfo, error)
}

type AppDatabaseImp struct {
	DatabaseBasePath string
}

func NewAppDatabase(basePath string) *AppDatabaseImp {
	appdb := &AppDatabaseImp{basePath}
	_, err := appdb.setupDatabase()
	if err != nil {
		log.Printf("Database setup failed: %q \n", err)
	}
	return appdb
}

func (adi AppDatabaseImp) GetAll() ([]os.FileInfo, error) {
	files, _ := ioutil.ReadDir(adi.DatabaseBasePath)
	return files, nil
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

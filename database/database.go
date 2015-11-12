package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var defaultFolders = [3]string{"movie", "music", "picture"}

type AppDatabase interface {
	GetAll(string) ([]File, error)
	GetFile(filetype string, searchValue string) ([]File, error)
	DeleteFile(filetype string, searchValue string) ([]File, error)
	SaveFile(filetype string, filename string, file []byte)(error)
}

type AppDatabaseImp struct {
	DatabaseBasePath string
}

type File struct {
	name     string
	size     int64
	filetype string
}

type dbError struct {
	When time.Time
	What string
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

// traverse all folders to get all files
func walkfnc(path string, f os.FileInfo, err error) error {
	if err != nil {
		log.Printf(path)
		return nil
	}

	if !f.IsDir() {
		splitPath := strings.SplitAfter(path, string(os.PathSeparator))
		if splitPath[len(splitPath)-1] == f.Name() {
			// Path without suffix filesep
			path := splitPath[len(splitPath)-2]
			filesArr = append(filesArr, File{f.Name(), f.Size(), path[:len(path)-1]})
		}
	}

	return err
}

/*
*	Get a specific file for a searchvalue and filetype
 */

func (adi AppDatabaseImp) GetFile(filetype string, searchValue string) ([]File, error) {
	filesArr = filesArr[:0]
	if len(filetype) > 0 {
		files, err := ioutil.ReadDir(adi.DatabaseBasePath + filetype)
		if err != nil {
			return nil, err
		}
		for _, value := range files {
			// Search for the right file in the folder
			searchValueIndex := strings.LastIndex(searchValue, ".")

			if searchValueIndex == -1 {
				searchValueIndex = len(searchValue)
			}
			searchString := searchValue[:searchValueIndex]
			filename := value.Name()[:strings.LastIndex(value.Name(), ".")]
			if searchString == filename {
				filesArr = append(filesArr, File{value.Name(), value.Size(), filetype})
			}

		}
		return filesArr, nil

	} else {
		// No filetype, no file
		return []File{}, nil
	}

	return []File{}, nil
}

func (adi AppDatabaseImp) DeleteFile(filetype string, filename string) ([]File, error) {

	if len(filetype) > 0 {
		files, err := adi.GetFile(filetype, filename)
		if err != nil {
			return []File{}, err
		}
		if len(files) == 1 {
			return []File{}, os.Remove(adi.DatabaseBasePath + filetype + string(os.PathSeparator) + files[0].name)
		}
	} else {

		return []File{}, dbError{
			time.Now(),
			"No filetype was provided",
		}
	}

	return []File{}, nil
}

func (adi AppDatabaseImp) SaveFile(fileType string, filename string, file []byte)(error){

	if len(fileType) > 0 && len(filename) > 0{
		err := ioutil.WriteFile(adi.DatabaseBasePath + fileType + string(os.PathSeparator) + filename, file, 0777)
		if err != nil {
			return err;
		}
	}else{
		return dbError{
			time.Now(),
			"No filetype/filename was provided",
		}
	}

	return nil
}

func (e dbError) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
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

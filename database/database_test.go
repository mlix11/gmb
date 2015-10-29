package database

import (
	"io"
	"log"
	"os"
	"reflect"
	"testing"
)

var (
	testFolder = "test"
	basePath   = testFolder + string(os.PathSeparator)
)

type DatabaseTest struct {
	fileType      string
	expectedFiles []File
}

var db AppDatabase

func setup() {
	os.Mkdir(testFolder, 0777)
	db = NewAppDatabase(basePath)

	var testFiles = [3]string{"broke.mp3", "iowa.mp4", "mountain.jpeg"}

	for _, value := range testFiles {
		var dest string
		switch value {
		case "broke.mp3":
			dest = "./test/music"
			break
		case "iowa.mp4":
			dest = "./test/movie"
			break
		case "mountain.jpeg":
			dest = "./test/picture"
			break
		}
		// open orginal file
		orginalFile, err := os.Open("../testData/" + value)
		if err != nil {
			log.Fatal(err)
		}
		defer orginalFile.Close()

		newFile, err := os.Create(dest + "/" + value)
		if err != nil {
			log.Fatal(err)
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, orginalFile)
		if err != nil {
			log.Fatal(err)
		}

		// Commit the file contents
		// flushes memory to disk
		err = newFile.Sync()
		if err != nil {
			log.Fatal(err)
		}

	}

}

func teardown() {
	os.RemoveAll(testFolder)
}

func TestFolderExist(t *testing.T) {
	setup()

	for _, value := range defaultFolders {
		isExisting, err := exists(basePath + value)

		if err != nil || !isExisting {
			teardown()
			t.Fatalf("Folder dont exist %q --> %q", value, basePath)
		}
	}
	teardown()

}

func TestGetAll(t *testing.T) {
	setup()

	getAllTests := []DatabaseTest{
		{
			"music", []File{File{"broke.mp3", 8054458, "music"}},
		},
		{
			"movie", []File{File{"iowa.mp4", 3173020, "movie"}},
		},
		{
			"picture", []File{File{"mountain.jpeg", 1467536, "picture"}},
		},
		{
			"", []File{File{"iowa.mp4", 3173020, "movie"}, File{"broke.mp3", 8054458, "music"}, File{"mountain.jpeg", 1467536, "picture"}},
		},
	}

	for _, dt := range getAllTests {
		sli, _ := db.GetAll(dt.fileType)
		if !reflect.DeepEqual(sli, dt.expectedFiles) {
			t.Fatalf("Expected array to be %q, but was %q", dt.expectedFiles, sli)
		}

	}
	teardown()
}

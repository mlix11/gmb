package database

import (
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
	testFolder = "test"
	basePath   = testFolder + string(os.PathSeparator)
)

type DatabaseTest struct {
	fileType      string
	expectedFiles []File
}

type SearchDatabaseTest struct {
	fileType      string
	searchValue   string
	expectedFiles []File
}

type DeleteDatabaseTest struct {
	SearchDatabaseTest
	Error dbError
}

var db AppDatabase

func setup() {
	os.Mkdir(testFolder, 0777)
	db = NewAppDatabase(basePath)

	var testFiles = [3]string{"broke.mp3", "iowa.mp4", "mountain.jpeg"}
	var base = "." + string(os.PathSeparator) + "test" + string(os.PathSeparator)
	for _, value := range testFiles {
		var dest string
		switch value {
		case "broke.mp3":
			dest = base + "music"
			break
		case "iowa.mp4":
			dest = base + "movie"
			break
		case "mountain.jpeg":
			dest = base + "picture"
			break
		}
		// open orginal file
		orginalFile, err := os.Open(".." + string(os.PathSeparator) + "testData" + string(os.PathSeparator) + value)
		if err != nil {
			log.Fatal(err)
		}
		defer orginalFile.Close()

		newFile, err := os.Create(dest + string(os.PathSeparator) + value)
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

func TestGetFile(t *testing.T) {
	setup()

	getTest := []SearchDatabaseTest{
		{
			"music", "broke.mp3", []File{File{"broke.mp3", 8054458, "music"}},
		},
		{
			"music", "broke", []File{File{"broke.mp3", 8054458, "music"}},
		},
		{
			"movie", "iowa.mp4", []File{File{"iowa.mp4", 3173020, "movie"}},
		},
		{
			"movie", "iowa", []File{File{"iowa.mp4", 3173020, "movie"}},
		},
		{
			"picture", "mountain.jpeg", []File{File{"mountain.jpeg", 1467536, "picture"}},
		},
		{
			"picture", "mountain", []File{File{"mountain.jpeg", 1467536, "picture"}},
		},
		{
			"movie", "notFound", []File{},
		},
	}

	for _, dt := range getTest {
		files, _ := db.GetFile(dt.fileType, dt.searchValue)
		if !reflect.DeepEqual(files, dt.expectedFiles) {
			t.Fatalf("Expected array for search %q to be %q, but was %q", dt.searchValue, dt.expectedFiles, files)
		}
	}

	teardown()
}

func TestDeleteFile(t *testing.T) {
	setup()

	testValues := []DeleteDatabaseTest{
		{
			SearchDatabaseTest{"music", "broke.mp3", []File{}},
			dbError{},
		},
		{
			SearchDatabaseTest{"picture", "mountain", []File{}},
			dbError{},
		},
		{
			SearchDatabaseTest{"", "mountain", []File{}},
			dbError{
				// Time doesn't matter
				time.Now(), "No filetype was provided",
			},
		},
	}

	for _, dt := range testValues {
		files, err := db.DeleteFile(dt.fileType, dt.searchValue)
		if err != nil {
			errorCmpMsg := err.Error()
			message := errorCmpMsg[strings.LastIndex(errorCmpMsg, ":")+2:]
			if message != dt.Error.What {
				t.Fatalf("Expected no error: %q, blaba: %s", err, message)
			}
		}
		if len(files) > 0 {
			t.Fatalf("There should be no files like: %q", files)
		}
		filesWithName, _ := db.GetFile(dt.fileType, dt.searchValue)
		if len(filesWithName) > 0 {
			t.Fatalf("Expected no file with value %s, but was %q", dt.searchValue, filesWithName)
		}

	}

	teardown()
}

package database

import (
	"os"
	//"reflect"
	"testing"
)

var (
	testFolder = "test"
	basePath   = testFolder + string(os.PathSeparator)
)

func setup() {
	os.Mkdir(testFolder, 0777)
}

func teardown() {
	os.RemoveAll(testFolder)
}

func TestFolderExist(t *testing.T) {
	setup()
	NewAppDatabase(basePath)

	for _, value := range defaultFolders {
		isExisting, err := exists(basePath + value)

		if err != nil || !isExisting {
			teardown()
			t.Fatalf("Folder dont exist %q --> %q", value, basePath)
		}
	}
	teardown()

}

/*func TestGetAll(t *testing.T) {
	database := AppDatabaseImp{basePath}

	sli, _ := database.GetAll()
	exp := []string{"Penn", "Teller"}
	if !reflect.DeepEqual(sli, exp) {
		t.Fatalf("Expected array to be %q, but was %q", exp, sli)
	}
}*/

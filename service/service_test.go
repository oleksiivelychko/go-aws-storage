package service

import (
	"fmt"
	"github.com/oleksiivelychko/go-aws-storage/config"
	"log"
	"os"
	"strings"
	"testing"
)

const (
	testBucket   = "test-bucket"
	testFilename = "sample-0.png"
)

var (
	storageService IService
	tmpFiles       []string
)

func init() {
	yamlConfig, err := config.ReadYAML("./../config.yaml")
	if err != nil {
		log.Fatalf(err.Error())
	}

	storageService, err = New(yamlConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func TestMain(m *testing.M) {
	tearDown := setUp()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setUp() func() {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "files")
	if err != nil {
		log.Fatal(err)
	}

	var files []*os.File

	for i := 0; i < 10; i++ {
		filename := fmt.Sprintf("%s/sample-%d.png", tmpDir, i)

		file, fileErr := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
		if fileErr != nil {
			_ = fmt.Errorf(fileErr.Error())
			continue
		}

		files = append(files, file)
		tmpFiles = append(tmpFiles, filename)
	}

	return func() {
		for i := 0; i < len(files); i++ {
			if _, statErr := files[i].Stat(); statErr == nil {
				_ = files[i].Close()
			}
		}

		_ = os.RemoveAll(tmpDir)
	}
}

func TestCreateBucket(t *testing.T) {
	output, err := storageService.CreateBucket(testBucket, false)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
}

func TestListBuckets(t *testing.T) {
	output, err := storageService.ListBuckets()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
}

func TestPutObjectsAsync(t *testing.T) {
	outCh := make(chan string, len(tmpFiles))
	errCh := make(chan error, len(tmpFiles))

	storageService.PutObjectsAsync(testBucket, tmpFiles, outCh, errCh)

	if len(errCh) > 0 {
		for err := range errCh {
			t.Errorf("%s\n", err)
		}
	}

	if len(outCh) > 0 {
		for output := range outCh {
			t.Logf("%s\n", output)
		}
	}
}

func BenchmarkPutObjectsAsync(b *testing.B) {
	for i := 0; i < b.N; i++ {
		outCh := make(chan string, len(tmpFiles))
		errCh := make(chan error, len(tmpFiles))
		storageService.PutObjectsAsync(testBucket, tmpFiles, outCh, errCh)
	}
}

func TestPutObjects(t *testing.T) {
	output, errArr := storageService.PutObjects(testBucket, tmpFiles)

	if len(errArr) > 0 {
		for _, err := range errArr {
			t.Errorf("%s\n", err.Error())
		}
	}

	for _, eTag := range output {
		t.Logf("%s\n", eTag)
	}
}

func BenchmarkPutObjects(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = storageService.PutObjects(testBucket, tmpFiles)
	}
}

func TestListObjects(t *testing.T) {
	output, err := storageService.ListObjects(testBucket)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
}

func TestGetObject(t *testing.T) {
	uploadPath := "./../upload"

	err := storageService.GetObject(testBucket, testFilename, uploadPath)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = os.Stat(uploadPath); os.IsNotExist(err) {
		t.Fatal(err)
	}

	err = os.RemoveAll(uploadPath)
	if err != nil {
		t.Error(err)
	}
}

func TestAssignURL(t *testing.T) {
	output, err := storageService.AssignURL(testBucket, testFilename)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(output)
}

func TestDeleteObject(t *testing.T) {
	for _, filename := range tmpFiles {
		splitted := strings.Split(filename, "/")
		err := storageService.DeleteObject(testBucket, splitted[len(splitted)-1])
		if err != nil {
			t.Error(err)
		}
	}
}

func TestDeleteBucket(t *testing.T) {
	err := storageService.DeleteBucket(testBucket)
	if err != nil {
		t.Error(err)
	}
}

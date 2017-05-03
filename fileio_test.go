package fileio_test

import (
	"fmt"

	"github.com/rvflash/fileio"
)

func ExampleUpload() {
	key, err := fileio.Upload("test.txt")
	if err != nil {
		panic(err)
	}
	if err := fileio.Download(key, "/tmp/test.txt"); err != nil {
		panic(err)
	}
	fmt.Println("Yeah, good boy.")
	// Output: Yeah, good boy.
}

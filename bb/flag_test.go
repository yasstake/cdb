package bb

import (
	"fmt"
	"testing"
)

func TestFlagInit(t *testing.T) {
	var file FlagFile
	file.Init("/tmp/PROCESSA")
}

func TestFlagCreate(t *testing.T) {
	var file FlagFile
	file.Init("/tmp/PROCESSA")
	file.Create()
}

func TestFlagDelete(t *testing.T) {
	var file FlagFile
	file.Init("/tmp/PROCESSA")
	file.Create()
	file.Delete()
}

func TestFlagFileExist(t *testing.T) {
	var file FlagFile
	file.Init("/tmp/PROCESSA")
	file.Create()

	r := file.exist_other()

	if r {
		t.Error("Flag file must be my file")
	}

	fmt.Println(r)

}

func TestFlagCheck(t *testing.T) {
	var file FlagFile
	file.Init("/tmp/PROCESSA")
	file.Create()

	done := make(chan struct{})

	go file.Check_other_process_loop(20, done)

	var file2 FlagFile
	file2.Init("/tmp/PROCESSA")
	file2.Create()

	select {
	case <-done:
		fmt.Println("DONE")
	}
}

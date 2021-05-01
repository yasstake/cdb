package bb

import (
	"io/ioutil"
	"log"
	"os"
	"time"
)

type FlagFile struct {
	Path       string
	process_id string
}

func (c *FlagFile) Init(file_name string) {
	c.Path = file_name
}

func (c *FlagFile) Create() {
	if c.Path == "" {
		log.Println("Logfile is not setup(Create)")
		return
	}

	f, _ := os.Create(c.Path)
	defer f.Close()

	t := time.Now()

	c.process_id = t.Format("2006-01-02-15-04-05.000000")
	f.WriteString(c.process_id)

	log.Println("Flag file created", c.Path, c.process_id)
}

func (c *FlagFile) Delete() {
	if c.Path == "" {
		log.Println("Logfile is not setup(Delete)")
		return
	}

	os.Remove(c.Path)
}

func (c *FlagFile) exist_other() bool {
	// if not initialized(non-flag mode), ignore other process
	if c.Path == "" {
		return false
	}

	buffer, err := ioutil.ReadFile(c.Path)

	if err != nil {
		log.Println("cannot open flag file", c.Path, err)
		return false
	}

	id := string(buffer)

	return id != c.process_id
}

func (c *FlagFile) Check_other_process_loop(sec int, done chan struct{}) {
	const COUNT_STEP = 10
	count := 0

	for {
		if c.exist_other() {
			count += COUNT_STEP
		}

		if 0 < count {
			log.Printf("[Other Process flag found]%s %d / %d[sec]", c.Path, count, sec)
		}

		if sec < count {
			close(done)
			break
		}

		time.Sleep(time.Second * COUNT_STEP)
	}
}

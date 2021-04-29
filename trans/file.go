package trans

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"time"
)

// Find files in recusively
func file_list(base_path string) []string {
	files, err := ioutil.ReadDir(base_path)

	if err != nil {
		log.Fatal("File canot open", base_path, err)
		fmt.Println("Error")
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, file_list(filepath.Join(base_path, file.Name()))...)
		} else {
			paths = append(paths, filepath.Join(base_path, file.Name()))
		}
	}

	return sort.StringSlice(paths)
}

// represent start, end time frame
type TimeFrame struct {
	start time.Time
	end   time.Time
}

// Conver to string for visual representation
func (c TimeFrame) to_string() string {
	return c.start.String() + "->" + c.end.String()
}

type TimeFrames []TimeFrame

func (c TimeFrames) ToString() string {
	frames := len(c)
	result := ""
	for i := 0; i < frames; i++ {
		result = result + "\n" + c[i].to_string()
	}
	return result
}

// Search DB directory and find time frame series chunks
// [{Start, end} {start, end},,,,,]
// Direcotry
//     BASE_PATH
//         +-----YYYY-MM-DD
//                    +------ HH-MM.log.gz
func Time_chunks(base_path string) (times TimeFrames) {
	// Open LogDir and sort
	files, err := ioutil.ReadDir(base_path)

	if err != nil {
		fmt.Println("Error")
	}

	// Open YYYY
	var dirs []string
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, filepath.Join(base_path, file.Name()))
		}
	}
	dirs = sort.StringSlice(dirs)

	// Open each log dir and sort each logs
	// Open MM
	for _, dir_path := range dirs {
		var file_paths []string

		files, err := ioutil.ReadDir(dir_path)
		if err != nil {
			fmt.Println("Error")
		}

		// Open DD
		for _, file := range files {
			name := file.Name()
			file_paths = append(file_paths, filepath.Join(dir_path, name))
		}

		file_paths = sort.StringSlice(file_paths)

		for _, file := range file_paths {
			time := file_to_time(file)
			times = append_time_chunks(times, time)
		}
	}

	return times
}

// MAX of time duration alloance for time skip(losting data period)
const TIME_GAP = 60 + 5

// Append time chunks. Add new time to old TimeFrame
func append_time_chunks(org TimeFrames, now time.Time) []TimeFrame {
	if org == nil {
		org = []TimeFrame{{now, now}}
		return org
	}

	end := org[len(org)-1].end
	diff := now.Sub(end)

	diff_sec := diff.Seconds()

	switch {
	case diff_sec < 0:
		fmt.Println("ERROR in append_time_chunks", now.String(), end.String(), diff)
	case diff_sec <= TIME_GAP:
		pos := len(org) - 1
		org[pos].end = now
	case TIME_GAP < diff_sec:
		org = append(org, TimeFrame{now, now})
	default:
		fmt.Println("Unexpected case append_time_chunks")
	}

	return org
}

package jhlog

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	dash    = "-"
	dayHour = 24
)

var ()

type File struct {
	Mu *sync.Mutex
	// file prefix
	file *os.File
	// format for time
	format string

	// if need auto delete
	autoDeleteDays int
	autoDelete     bool

	// rotate time
	rotate     bool
	rotateGaps int

	// current log gap
	currentgap int
	// current file gap
	lastgap int
	// gaps since currentgap
	gaps int

	writeCount uint
}

func (f *File) close() error {
	if f.file == nil {
		return nil
	}
	f.file.Sync()
	f.file.Close()
	f.file = nil
	return nil
}

func (f *File) open() error {
	t := time.Now()
	f.currentgap = getCurrentGap(t)
	f.lastgap = f.currentgap
	f.gaps = 0
	fname := filepath.Base(f.format)
	dir := filepath.Dir(f.format)
	fname = t.Format(fname)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	path := dir + string(os.PathSeparator) + fname
	flag := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	f.file, err = os.OpenFile(path, flag, 0644)
	return err
}

func deleteOldLogFile(format string, days int) {
	dir := filepath.Dir(format)
	fname := filepath.Base(format)
	du := time.Duration(dayHour * days)
	last := time.Now().Add(-time.Hour * du).Format(fname)
	files, _ := ioutil.ReadDir(dir)
	start := strings.SplitN(fname, dash, 2)[0]
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if name < last && strings.HasPrefix(name, start) {
			os.RemoveAll(dir + string(os.PathSeparator) + f.Name())
		}
	}
}

func getCurrentGap(t time.Time) int {
	return t.YearDay()
}

func (f *File) rotateFile() error {
	if !f.rotate {
		return nil
	}
	t := time.Now()
	cur := getCurrentGap(t)
	if cur != f.lastgap {
		f.lastgap = cur
		f.gaps++
	}
	if f.gaps < f.rotateGaps {
		return nil
	}
	err := f.close()
	if err != nil {
		return err
	}
	return f.open()
}

func (f *File) periodicRotate() {
	for {
		time.Sleep(23 * time.Hour)
		f.RotateFile()
		if f.autoDelete {
			go deleteOldLogFile(f.format, f.autoDeleteDays)
		}
	}
}

func (f *File) RotateFile() {
	f.Mu.Lock()
	defer f.Mu.Unlock()
	f.rotateFile()
}

func NewFile(base string, timeformat string) (*File, error) {
	f := &File{
		Mu:         new(sync.Mutex),
		format:     base + dash + timeformat + ".log",
		autoDelete: false,
		rotate:     false,
	}
	f.open()
	go f.periodicRotate()
	return f, nil
}

func (f *File) SetAutoDelete(days int) {
	f.Mu.Lock()
	defer f.Mu.Unlock()
	f.autoDeleteDays = days
	f.autoDelete = true
}

func (f *File) SetRotate(days int) {
	f.Mu.Lock()
	f.Mu.Unlock()
	f.rotate = true
	f.rotateGaps = days
}

func (f *File) Close() error {
	f.Mu.Lock()
	defer f.Mu.Unlock()
	return f.close()
}

func (f *File) write(d []byte, flush bool) (int, error) {
	f.writeCount++
	err := f.rotateFile()
	if err != nil {
		return 0, err
	}
	n, err := f.file.Write(d)
	if err != nil {
		return n, err
	}
	if flush {
		err = f.file.Sync()
	}
	return n, err
}

func (f *File) Write(d []byte) (int, error) {
	f.Mu.Lock()
	defer f.Mu.Unlock()
	return f.write(d, false)
}

func (f *File) Flush() error {
	f.Mu.Lock()
	defer f.Mu.Unlock()
	return f.file.Sync()
}

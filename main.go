package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func calcCheckSum(s string) int {
	result := 0
	for _, b := range []byte(s) {
		result ^= int(b)
	}
	return result
}

func convertToBackupPath(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))] + ".bak"
}

func correctDate(s string) (string, error) {
	t, err := time.Parse("020106", s)
	if err != nil {
		return "", err
	}
	t = t.Add(time.Duration(1024 * 7 * 24 * 60 * 60 * 1000000000))
	return t.Format("020106"), nil
}

func isBrokenLine(s string) bool {
	if s[0] != '$' {
		return false
	}
	if len(s)-3 < 0 || s[len(s)-3] != '*' {
		return false
	}
	return fmt.Sprintf("%02X", calcCheckSum(s[1:len(s)-3])) != s[len(s)-2:]
}

func correctCheckSum(s string) string {
	c := calcCheckSum(s[1 : len(s)-3])
	return s[:len(s)-2] + fmt.Sprintf("%02X", c)
}

func correctGprmc(s string) (string, time.Time) {
	if !strings.HasPrefix(s, "$GPRMC,") {
		return s, time.Time{}
	}

	v := strings.Split(s, ",")
	if len(v) != 13 {
		log.Fatalf("$GPRMC の要素数が13ではありません. %v\n", s)
	}

	d, err := correctDate(v[9])
	if err != nil {
		log.Fatalf("$GPRMC の日付が ddmmyy ではありません. (%v)\n", v[9])
	}
	v[9] = d

	s = correctCheckSum(strings.Join(v, ","))
	t, err := time.Parse("020106150405", (v[9] + v[1])[:12])
	if err != nil {
		log.Fatalf("$GPRMC の時刻が hhmmss.ss ではありません. (%v)\n", v[9])
	}
	return s, t
}

func correctFileTime(fileName string, dateCreated, dateLastModified time.Time) error {
	handle, err := syscall.CreateFile(
		syscall.StringToUTF16Ptr(fileName),
		syscall.FILE_WRITE_ATTRIBUTES,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0)
	if err != nil {
		return fmt.Errorf("%s のハンドルの取得に失敗しました", fileName)
	}
	defer syscall.CloseHandle(handle)

	ctime := syscall.NsecToFiletime(dateCreated.UnixNano())
	atime := syscall.NsecToFiletime(time.Now().UnixNano())
	mtime := syscall.NsecToFiletime(dateLastModified.UnixNano())
	err = syscall.SetFileTime(handle, &ctime, &atime, &mtime)
	if err != nil {
		return fmt.Errorf("%s のファイルタイムの更新に失敗しました", fileName)
	}
	return nil
}

func correctFileContent(fileName string) (time.Time, time.Time) {
	backupPath := convertToBackupPath(fileName)
	if err := os.Rename(fileName, backupPath); err != nil {
		log.Fatalf("%s から %s へのリネームに失敗しました.\n%v\n", fileName, backupPath, err)
	}

	r, err := os.Open(backupPath)
	if err != nil {
		log.Fatalf("%s を開くのに失敗しました.\n%v\n", backupPath, err)
	}
	defer r.Close()

	w, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("%s を開くのに失敗しました.\n%v\n", fileName, err)
	}
	defer w.Close()

	reader := bufio.NewScanner(r)
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	var (
		dateCreated      time.Time
		dateLastModified time.Time
	)
	for reader.Scan() {
		l := reader.Text()
		if isBrokenLine(l) {
			continue
		}
		if strings.HasPrefix(l, "$GPRMC,") {
			var t time.Time
			l, t = correctGprmc(l)
			if dateCreated.IsZero() {
				dateCreated = t
			}
			dateLastModified = t
		}
		writer.WriteString(l)
		writer.WriteString("\r\n")
	}

	return dateCreated, dateLastModified
}

func createFileName(dateCreated time.Time) string {
	return dateCreated.Format("01021504")
}

func correctFileName(fileName string, dateCreated time.Time) {
	os.Rename(fileName, filepath.Join(filepath.Dir(fileName), createFileName(dateCreated)+".nma"))
}

func main() {
	log.SetFlags(log.Lshortfile)
	for i := 1; i < len(os.Args); i++ {
		dateCreated, dateLastModified := correctFileContent(os.Args[i])
		correctFileTime(os.Args[i], dateCreated, dateLastModified)
		correctFileName(os.Args[i], dateCreated)
	}
}

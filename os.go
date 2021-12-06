package gtools

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/h2non/filetype"
)

// Exists 判断所给路径文件/文件夹是否存在
func Exists(file string) bool {
	_, err := os.Stat(file) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// Writable 文件/文件夹是否可写
//func Writable(file string) (bool, error) {
//	if !Exists(file) {
//		err := fmt.Errorf("file does not exist: `%s`", file)
//		return false, err
//	}
//
//	err := syscall.Access(file, syscall.O_RDWR)
//	if err != nil {
//		return false, err
//	}
//
//	return true, nil
//}

// IsDir 判断所给路径是否为文件夹
func IsDir(file string) bool {
	s, err := os.Stat(file)
	if err != nil {
		return false
	}

	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(file string) bool {
	stat, err := os.Stat(file)
	if err != nil {
		return false
	}

	fm := stat.Mode()
	return fm.IsRegular()
}

func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return "localhost"
	} else {
		return name
	}
}

func DetectFileType(filename string) (string, string, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return "unknown", "", err
	}

	return DetectFileByteType(buf)
}

func DetectFileByteType(buf []byte) (extension, mime string, err error) {
	kind, unknown := filetype.Match(buf)
	if unknown != nil {
		extension = "unknown"
		err = unknown
		return
	}

	extension = kind.Extension
	mime = kind.MIME.Value

	return
}

// GetFileExt 简易版取文件名后缀,path.Ext()方法会带着个`.`
func GetFileExt(filename string) (suffix string) {
	exp := strings.Split(filename, ".")
	expLen := len(exp)
	if expLen > 1 {
		suffix = exp[expLen-1]
	}

	return
}

// Remove 安全删除文件
func Remove(filename string) (err error) {
	_, err = os.Stat(filename)
	if err != nil {
		log.Println("file does not exist:", filename)
		return
	}

	err = os.Remove(filename)

	return
}

func FullStack() string {
	var buf [2 << 11]byte
	runtime.Stack(buf[:], true)
	return string(buf[:])
}

func ClearOnSignal(handler func()) {
	signalChan := make(chan os.Signal, 1)

	// SIGINT  2  用户发送INTR字符(Ctrl+C)触发
	// SIGTERM 15 结束程序(可以被捕获、阻塞或忽略)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-signalChan
		log.Printf(`got signal to exit, signal: %s`, sig)
		handler()
		os.Exit(0)
	}()
}

// SetupSignalHandler setup signal handler
func SetupSignalHandler(shutdownFunc func()) {
	usrDefSignalChan := make(chan os.Signal, 1)

	signal.Notify(usrDefSignalChan, syscall.SIGUSR1)
	go func() {
		buf := make([]byte, 1<<16)
		for {
			sig := <-usrDefSignalChan
			if sig == syscall.SIGUSR1 {
				stackLen := runtime.Stack(buf, true)
				log.Printf("\n=== Got signal [%s] to dump goroutine stack. ===\n%s\n=== Finished dumping goroutine stack. ===\n", sig, buf[:stackLen])
			}
		}
	}()

	closeSignalChan := make(chan os.Signal, 1)
	signal.Notify(closeSignalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-closeSignalChan
		log.Printf("got signal to exit, signal: %v", sig)
		shutdownFunc()
	}()
}

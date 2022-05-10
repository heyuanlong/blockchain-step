package log

import (
	"fmt"
	klog "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"sort"
	"strings"
)

func init() {
	klog.SetReportCaller(true)
	klog.SetOutput(os.Stdout)
	klog.SetFormatter(&klog.TextFormatter{
		TimestampFormat : "2006-01-02 15:04:05",
		DisableTimestamp: false,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			//s := strings.Split(f.Function, ".")
			//funcname := s[len(s)-1]
			//_, filename := path.Split(f.File)
			i:=strings.LastIndex(f.File, "/")
			j:=strings.LastIndex(f.File[0:i], "/")
			return "", fmt.Sprintf("%s:%d", f.File[j:], f.Line)
		},
		SortingFunc: func(keys []string) {
			sort.Slice(keys, func(i, j int) bool {
				if keys[j] == klog.FieldKeyTime {
					return false
				}
				if keys[i] == klog.FieldKeyTime {
					return true
				}
				if keys[j] == klog.FieldKeyFile {
					return false
				}
				if keys[i] == klog.FieldKeyFile {
					return true
				}
				return strings.Compare(keys[i], keys[j]) == -1
			})
		},
	})

}


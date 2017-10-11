package backup

import (
	. "../config"
	"log"
	"os"
	"strconv"
	"strings"
)

const filename = "backup.txt"

func checkError(err error) {
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
}

func Remove() {
	err := os.RemoveAll(filename)
	checkError(err)
}

func Write(str string) {
	file, err := os.OpenFile("backup.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	checkError(err)
	defer file.Close()
	log.SetOutput(file)
	log.Println(str)
}

func Read(stringSize int64) string {
	stringSize += 1

	file, err := os.Open(filename)
	checkError(err)
	defer file.Close()

	buf := make([]byte, stringSize)
	stat, err := os.Stat(filename)
	start := stat.Size() - stringSize
	n, err := file.ReadAt(buf, start)
	checkError(err)
	buf = buf[:n]

	return string(buf)
}

func String(orders [N_FLOORS][N_BUTTONS]int) string {
	var queueString string
	var orderString string
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if orders[f][b] == 1 {
				orderString = "1 "
			} else {
				orderString = "0 "
			}
			queueString += orderString
		}
	}
	return queueString
}

func Queue(queueString string) [N_FLOORS][N_BUTTONS]int {
	queue := [N_FLOORS][N_BUTTONS]int{}
	index := 0
	queueTemp := strings.Fields(queueString)
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			newQueue, _ := strconv.Atoi(queueTemp[index])
			queue[f][b] = newQueue
			index += 1
		}
	}
	return queue
}
func Exists() bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

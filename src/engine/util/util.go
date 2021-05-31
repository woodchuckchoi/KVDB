package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

const (
	BASE_DIR = "/tmp/gokvdb"
	TIME_FMT = "2006-01-02T15:04:05"
)

func KeyValueToByteSlice(kv vars.KeyValue) []byte {
	ret := make([]byte, len(kv.Key)+len(kv.Value)+2)
	idx := 0
	for _, val := range []byte(kv.Key) {
		ret[idx] = val
		idx++
	}
	ret[idx] = vars.SEPARATOR
	idx++
	for _, val := range []byte(kv.Value) {
		ret[idx] = val
		idx++
	}
	ret[idx] = vars.DELIMITER
	return ret
}

func KeyValueSliceToByteSliceAndSparseIndex(kvPairs []vars.KeyValue) ([]byte, []vars.SparseIndex) {
	byteSlice := []byte{}
	sparseIndex := []vars.SparseIndex{}
	beforeOffset, curOffset := 0, 0
	for idx, kvPair := range kvPairs {
		byteKvPair := KeyValueToByteSlice(kvPair)
		byteSlice = append(byteSlice, byteKvPair...)
		if idx == 0 || idx == len(kvPairs)-1 || curOffset-beforeOffset >= vars.INDEX_TERM {
			sparseIndex = append(sparseIndex, vars.SparseIndex{
				Key:    kvPair.Key,
				Offset: curOffset,
			})
			beforeOffset = curOffset
		}
		curOffset += len(byteKvPair)
	}

	return byteSlice, sparseIndex
}

func KeyValueSliceToByteSlice(kvPairs []vars.KeyValue) []byte {
	ret := []byte{}
	for _, kvPair := range kvPairs {
		for _, b := range []byte(kvPair.Key) {
			ret = append(ret, b)
		}
		ret = append(ret, vars.SEPARATOR)
		for _, b := range []byte(kvPair.Value) {
			ret = append(ret, b)
		}
		ret = append(ret, vars.DELIMITER)
	}
	return ret
}

func ByteSliceToKeyValue(byteSlice []byte) []vars.KeyValue {
	ret := []vars.KeyValue{}
	byteKeyValuePairs := bytes.Split(byteSlice, []byte{vars.DELIMITER})
	for _, byteKeyValuePair := range byteKeyValuePairs {
		if len(byteKeyValuePair) == 0 { // end of file empty slice
			continue
		}

		byteKeyAndValue := bytes.Split(byteKeyValuePair, []byte{vars.SEPARATOR})
		key := string(byteKeyAndValue[0])
		value := string(byteKeyAndValue[1])
		ret = append(ret, vars.KeyValue{
			Key:   key,
			Value: value,
		})
	}
	return ret
}

func WriteKeyValuePairs(fileName string, keyValuePairs []vars.KeyValue) error {
	byteParsed := KeyValueSliceToByteSlice(keyValuePairs)
	return WriteByteSlice(fileName, byteParsed)
}

func WriteByteSlice(fileName string, byteSlice []byte) error {
	err := ioutil.WriteFile(fileName, byteSlice, 0777)
	if err != nil {
		return vars.FILE_CREATE_ERROR
	}
	return err
}

func ReadKeyValuePairs(fileName string, from, till int) ([]vars.KeyValue, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, vars.FILE_READ_ERROR
	}

	byteSlice, err := ReadFileChunk(f, from, till)
	return ByteSliceToKeyValue(byteSlice), nil

}

func ReadFileChunk(f *os.File, from, till int) ([]byte, error) {
	if till == -1 {
		info, err := f.Stat()
		if err != nil {
			return nil, vars.FILE_READ_ERROR
		}
		till = int(info.Size())
	}

	_, err := f.Seek(int64(from), 0)
	if err != nil {
		return nil, vars.FILE_READ_ERROR
	}

	byteChunk := make([]byte, till-from)
	_, err = f.Read(byteChunk)
	if err != nil {
		return nil, vars.FILE_READ_ERROR
	}

	return byteChunk, nil
}

func IntPow(a, b int) int {
	ret := 1
	for i := 0; i < b; i++ {
		ret *= b
	}
	return ret
}

func GenerateFileName(level, order int) string {
	postfix := time.Now().Local().Format(TIME_FMT)
	folderName := fmt.Sprintf("db-level-%d", level)
	fileName := fmt.Sprintf("db-%d-%d-%s.data", level, order, postfix)
	folderPath := path.Join(BASE_DIR, folderName)
	fullPath := path.Join(BASE_DIR, folderName, fileName)
	os.MkdirAll(folderPath, 0777)
	return fullPath
}

func RemoveFile(fileName string) {
	os.Remove(fileName)
}

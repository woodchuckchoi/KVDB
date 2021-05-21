package vars

import (
	"errors"
)

var (
	PUT_FAIL_ERROR = errors.New("PUT FAIL ERROR")
	GET_FAIL_ERROR = errors.New("GET FAIL ERROR")

	FILE_CREATE_ERROR = errors.New("FILE CREATE ERROR")
	FILE_WRITE_ERROR  = errors.New("FILE WRITE ERROR")
	FILE_READ_ERROR   = errors.New("FILE READ ERROR")

	FORMAT_ERROR = errors.New("FORMAT ERROR")

	KEY_NOT_FOUND_ERROR = errors.New("KEY NOT FOUND")
)
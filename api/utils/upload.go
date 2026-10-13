package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)



func UploadFileHandler(response http.ResponseWriter, request *http.Request, fileName string)(string, error) {

	if err := request.ParseMultipartForm(10 << 20); err != nil {
        return "", err
    }


	file, header, err := request.FormFile(fileName)

	if err != nil {
		return "", err
	}

	defer file.Close()

	fmt.Printf("Uploaded: %s (%d bytes)\n", header.Filename, header.Size)

	os.MkdirAll("./uploads", os.ModePerm) 


	dst, err := os.Create(filepath.Join("./uploads", header.Filename))

	if err != nil {
		return "", err
	}

	defer dst.Close()

	io.Copy(dst, file)

	return header.Filename, nil
}
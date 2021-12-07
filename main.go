package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"os"

	"github.com/KMDMNAK/zip"
)

// ビルド時に埋め込み
// go build -ldflags "-X main.password=example"
var (
	password string
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func encrypt_and_delete(input_file string) {

	// 圧縮先ファイルが既に存在するなら処理しない
	output_file := input_file + ".zip"

	if exists(output_file) {
		log.Fatalf("圧縮先ファイル%sが既に存在します\n", output_file)
		return
	}

	log.Println(input_file)

	contents, err := os.ReadFile(input_file)
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Println(contents)

	fzip, err := os.Create(output_file)
	if err != nil {
		log.Fatalln(err)
		return
	}

	zipw := zip.NewWriter(fzip)
	defer zipw.Close()

	w, err := zipw.Encrypt(input_file, password, zip.StandardEncryption)
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer zipw.Flush()

	_, err = io.Copy(w, bytes.NewReader(contents))
	if err != nil {
		log.Fatalln(err)
		return
	}

	os.Remove(input_file)
}

func main() {
	for _, input_file := range os.Args[1:] {
		encrypt_and_delete(input_file)
	}
}

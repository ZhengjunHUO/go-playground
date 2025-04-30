package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	pathToInput = "input/toxic.csv"
	//pathToCompress     = "compressed/grafana-7.3.7.tar.gz"
	pathToCompress     = "compressed/ingress.tgz"
	folderToDecompress = "decompressed"
)

func main() {
	//compressTarGz()
	decompressTarGz(pathToCompress, folderToDecompress)
}

func decompressTarGz(tarPath, targetPath string) {
	var compressed *os.File
	var gzr *gzip.Reader
	var tr *tar.Reader
	var err error

	if compressed, err = os.Open(tarPath); err != nil {
		log.Fatalln(err)
	}
	defer compressed.Close()

	if strings.HasSuffix(tarPath, ".gz") || strings.HasSuffix(tarPath, ".tgz") {
		if gzr, err = gzip.NewReader(compressed); err != nil {
			log.Fatalln(err)
		}
		defer gzr.Close()
		tr = tar.NewReader(gzr)
	} else {
		tr = tar.NewReader(compressed)
	}

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		name := fmt.Sprintf("%s/%s", targetPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(name, os.FileMode(header.Mode)); err != nil {
				log.Fatalln("Error occurred making dir: ", err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
				log.Fatalln("Error occurred making dir: ", err)
			}

			//fmt.Printf("%s type mode: %v\n", name, header.Mode)
			outFile, err := os.Create(name)
			if err != nil {
				log.Fatalln("Error occurred creating file: ", err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				log.Fatalln("Error occurred dumping data: ", err)
			}

			if err := os.Chmod(name, os.FileMode(header.Mode)); err != nil {
				log.Fatalln("Error occurred chmoding file ", name, err)
			}
		case tar.TypeSymlink:
			os.Symlink(header.Linkname, name)
		default:
			log.Fatalf("%v : Unknown type: %v\n", name, header.Typeflag)
		}
	}
}

func compressTarGz() {
	var target *os.File
	var buffer []byte
	var err error
	var gzw *gzip.Writer

	if target, err = os.OpenFile(pathToCompress, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		log.Fatalln(err)
	}
	defer target.Close()

	if gzw, err = gzip.NewWriterLevel(target, gzip.BestCompression); err != nil {
		log.Fatalln(err)
	}
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	if buffer, err = os.ReadFile(pathToInput); err != nil {
		log.Fatalln(err)
	}

	if buffer != nil {
		thd := &tar.Header{
			Name: path.Base(pathToInput),
			Mode: int64(0644),
			Size: int64(len(buffer)),
		}
		if err := tw.WriteHeader(thd); err != nil {
			println(err)
		}
		if _, err := tw.Write(buffer); err != nil {
			println(err)
		}
	}
}

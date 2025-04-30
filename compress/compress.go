package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	pathToInput = "input/toxic.csv"
	//pathToCompress     = "compressed/grafana-7.3.7.tar.gz"
	//pathToCompress     = "compressed/ingress.tgz"
	pathToCompress     = "compressed/helm-v3.17.0-linux-amd64.tar"
	folderToDecompress = "decompressed"
	zipped_file        = "compressed/source_calculator.zip"
	pathToUnzip        = "unzipped"
)

func main() {
	//CompressTarGz()
	DecompressTarGz(pathToCompress, folderToDecompress, true)
	Unzip(zipped_file, pathToUnzip, true)
}

func stripTopLevelPath(path string, doStrip bool) string {
	if doStrip {
		parts := strings.Split(path, string(filepath.Separator))

		if len(parts) > 1 {
			return filepath.Join(parts[1:]...)
		}
	}

	return path
}

func Unzip(zipPath, targetPath string, doStrip bool) error {
	var zr *zip.ReadCloser
	var err error

	if zr, err = zip.OpenReader(zipPath); err != nil {
		return fmt.Errorf("error occurred opening zip reader: %v", err)
	}
	defer zr.Close()

	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("error occurred mkdir %s: %v", targetPath, err)
	}

	for _, file := range zr.File {
		pathToFile := fmt.Sprintf("%s/%s", targetPath, stripTopLevelPath(file.Name, doStrip))
		if !strings.HasPrefix(pathToFile, filepath.Clean(targetPath)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %v", pathToFile)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(pathToFile, file.Mode()); err != nil {
				return fmt.Errorf("error occurred mkdir %s: %v", pathToFile, err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(pathToFile), 0755); err != nil {
			return fmt.Errorf("error occurred mkdir %s: %v", filepath.Dir(pathToFile), err)
		}

		outFile, err := os.OpenFile(pathToFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return fmt.Errorf("error occurred opening %v: %v", outFile, err)
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("error occurred opening reader for file %v: %v", file, err)
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return fmt.Errorf("error occurred dumping to file %v: %v", outFile, err)
		}
	}

	return nil
}

func DecompressTarGz(tarPath, targetPath string, doStrip bool) error {
	var compressed *os.File
	var gzr *gzip.Reader
	var tr *tar.Reader
	var err error

	if compressed, err = os.Open(tarPath); err != nil {
		return fmt.Errorf("error occurred opening tar file %v: %v", compressed, err)
	}
	defer compressed.Close()

	if strings.HasSuffix(tarPath, ".gz") || strings.HasSuffix(tarPath, ".tgz") {
		if gzr, err = gzip.NewReader(compressed); err != nil {
			return fmt.Errorf("error occurred gz reader: %v", err)
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
			return fmt.Errorf("error occurred reading package with tar reader: %v", err)
		}

		name := fmt.Sprintf("%s/%s", targetPath, stripTopLevelPath(header.Name, doStrip))
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(name, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("error occurred making dir %v: %v", name, err)
			}
		case tar.TypeSymlink:
			os.Symlink(header.Linkname, name)
		case tar.TypeLink:
			linkTarget := filepath.Join(targetPath, header.Linkname)
			if err := os.Link(linkTarget, name); err != nil {
				return fmt.Errorf("error occurred making hardlink %v: %v", name, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
				return fmt.Errorf("error occurred making dir %v: %v", filepath.Dir(name), err)
			}

			outFile, err := os.Create(name)
			if err != nil {
				return fmt.Errorf("error occurred creating file %v: %v", name, err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return fmt.Errorf("error occurred dumping data to %v: %v", outFile, err)
			}

			if err := os.Chmod(name, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("error occurred chmoding file %v: %v", name, err)
			}
		default:
			return fmt.Errorf("%v : Unknown type: %v", name, header.Typeflag)
		}
	}

	return nil
}

func CompressTarGz() error {
	var target *os.File
	var buffer []byte
	var err error
	var gzw *gzip.Writer

	if target, err = os.OpenFile(pathToCompress, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return err
	}
	defer target.Close()

	if gzw, err = gzip.NewWriterLevel(target, gzip.BestCompression); err != nil {
		return err
	}
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	if buffer, err = os.ReadFile(pathToInput); err != nil {
		return err
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

	return nil
}

package compress

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GzipCompress 压缩数据为Gzip格式
func GzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, err
	}

	if err := gzipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GzipDecompress 解压Gzip格式的数据
func GzipDecompress(compressedData []byte) ([]byte, error) {
	buf := bytes.NewBuffer(compressedData)
	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	return io.ReadAll(gzipReader)
}

// ZlibCompress 压缩数据为Zlib格式
func ZlibCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	zlibWriter := zlib.NewWriter(&buf)

	_, err := zlibWriter.Write(data)
	if err != nil {
		return nil, err
	}

	if err := zlibWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ZlibDecompress 解压Zlib格式的数据
func ZlibDecompress(compressedData []byte) ([]byte, error) {
	buf := bytes.NewBuffer(compressedData)
	zlibReader, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	return io.ReadAll(zlibReader)
}

// ZipDirectory 将目录压缩为Zip文件
func ZipDirectory(sourceDir, targetFile string) error {
	// 创建zip文件
	zipFile, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历源目录
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 创建zip头信息
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		// 设置头信息中的名称为相对路径
		header.Name = relPath

		// 设置压缩方法
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// 创建文件或目录条目
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// 如果是文件则写入文件内容
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// UnzipFile 解压Zip文件到指定目录
func UnzipFile(zipFile, targetDir string) error {
	// 打开zip文件
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 确保目标目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	// 解压所有文件
	for _, file := range reader.File {
		// 构建完整路径
		path := filepath.Join(targetDir, file.Name)

		// 检查越界路径
		if !strings.HasPrefix(path, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("非法的zip文件路径: %s", file.Name)
		}

		// 处理目录
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		// 创建文件
		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		// 打开源文件
		srcFile, err := file.Open()
		if err != nil {
			dstFile.Close()
			return err
		}

		// 复制内容
		_, err = io.Copy(dstFile, srcFile)
		dstFile.Close()
		srcFile.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// CreateTarGz 创建tar.gz归档文件
func CreateTarGz(sourceDir, targetFile string) error {
	// 创建目标文件
	file, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建gzip writer
	gzipWriter := gzip.NewWriter(file)
	defer gzipWriter.Close()

	// 创建tar writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// 遍历源目录
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		// 创建tar头信息
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		// 写入头信息
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// 如果是文件则写入文件内容
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// ExtractTarGz 解压tar.gz归档文件
func ExtractTarGz(tarGzFile, targetDir string) error {
	// 打开tar.gz文件
	file, err := os.Open(tarGzFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建gzip reader
	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// 创建tar reader
	tarReader := tar.NewReader(gzipReader)

	// 确保目标目录存在
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	// 解压所有文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// 构建完整路径
		path := filepath.Join(targetDir, header.Name)

		// 检查越界路径
		if !strings.HasPrefix(path, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("非法的tar文件路径: %s", header.Name)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			// 确保父目录存在
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}

			// 创建文件
			file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// 复制内容
			_, err = io.Copy(file, tarReader)
			file.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GzipCompressFile 压缩文件为Gzip格式
func GzipCompressFile(sourceFile, targetFile string) error {
	// 打开源文件
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	// 创建目标文件
	target, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer target.Close()

	// 创建gzip writer
	gzipWriter := gzip.NewWriter(target)
	defer gzipWriter.Close()

	// 复制内容
	_, err = io.Copy(gzipWriter, source)
	return err
}

// GzipDecompressFile 解压Gzip格式的文件
func GzipDecompressFile(sourceFile, targetFile string) error {
	// 打开源文件
	source, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer source.Close()

	// 创建gzip reader
	gzipReader, err := gzip.NewReader(source)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// 创建目标文件
	target, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer target.Close()

	// 复制内容
	_, err = io.Copy(target, gzipReader)
	return err
}

package etl

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	ZIPDataFile  string = ".zip"
	GZIPDataFile string = ".gz"
	TARDataFile  string = ".tar"
	TGZDataFile  string = ".tar.gz"
)

type DataZip struct {
	Type       string
	HasSelf    bool
	SourcePath string
	TargetPath string
	Clean      bool
}

func (dz *DataZip) TimeSourcePath(dt DataTime, t time.Time) {
	filedate := dt.OffsetTime(t)
	filepath := dz.SourcePath
	filepath = strings.Replace(filepath, "YYYY", fmt.Sprintf("%d", filedate.Year()), -1)
	filepath = strings.Replace(filepath, "MM", fmt.Sprintf("%02d", filedate.Month()), -1)
	filepath = strings.Replace(filepath, "DD", fmt.Sprintf("%02d", filedate.Day()), -1)
	filepath = strings.Replace(filepath, "HH24", fmt.Sprintf("%02d", filedate.Hour()), -1)
	filepath = strings.Replace(filepath, "MI", fmt.Sprintf("%02d", filedate.Minute()), -1)
	filepath = strings.Replace(filepath, "SS", fmt.Sprintf("%02d", filedate.Second()), -1)
	dz.SourcePath = filepath
}

func (dz *DataZip) TimeTargetPath(dt DataTime, t time.Time) {
	filedate := dt.OffsetTime(t)
	filepath := dz.TargetPath
	filepath = strings.Replace(filepath, "YYYY", fmt.Sprintf("%d", filedate.Year()), -1)
	filepath = strings.Replace(filepath, "MM", fmt.Sprintf("%02d", filedate.Month()), -1)
	filepath = strings.Replace(filepath, "DD", fmt.Sprintf("%02d", filedate.Day()), -1)
	filepath = strings.Replace(filepath, "HH24", fmt.Sprintf("%02d", filedate.Hour()), -1)
	filepath = strings.Replace(filepath, "MI", fmt.Sprintf("%02d", filedate.Minute()), -1)
	filepath = strings.Replace(filepath, "SS", fmt.Sprintf("%02d", filedate.Second()), -1)
	dz.TargetPath = filepath
}

func (dz *DataZip) Compress() (err error) {
	switch dz.Type {
	case ZIPDataFile:
		err = dz.Zipit(dz.SourcePath, dz.TargetPath, dz.HasSelf)
	case GZIPDataFile:
		err = dz.GZip(dz.SourcePath, dz.TargetPath)
	case TARDataFile:
		err = dz.Tar(dz.SourcePath, dz.TargetPath, dz.HasSelf, false)
	case TGZDataFile:
		err = dz.Tar(dz.SourcePath, dz.TargetPath, dz.HasSelf, true)
	default:
		return errors.New(fmt.Sprintf("压缩类型(%s)不支持", dz.Type))
	}
	if err != nil {
		return
	}
	if dz.Clean {
		err = os.RemoveAll(dz.SourcePath)
	}
	return
}

func (dz *DataZip) UnCompress() (err error) {
	switch dz.Type {
	case ZIPDataFile:
		err = dz.Unzip(dz.SourcePath, dz.TargetPath)
	case GZIPDataFile:
		err = dz.UnGZip(dz.SourcePath, dz.TargetPath)
	case TARDataFile:
		err = dz.UnTar(dz.SourcePath, dz.TargetPath, false)
	case TGZDataFile:
		err = dz.UnTar(dz.SourcePath, dz.TargetPath, true)
	default:
		return errors.New(fmt.Sprintf("压缩类型(%s)不支持", dz.Type))
	}
	if err != nil {
		return
	}
	if dz.Clean {
		err = os.Remove(dz.SourcePath)
	}
	return
}

func (dz *DataZip) GZip(source, target string) error {
	source = strings.Replace(source, "\\", "/", -1)
	if !Exists(source) {
		return errors.New("找不到文件" + source)
	}
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()
	filename := filepath.Base(source)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()
	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()
	_, err = io.Copy(archiver, reader)
	return err
}

func (dz *DataZip) UnGZip(sourcePath string, targetPath string) error {
	sourcePath = strings.Replace(sourcePath, "\\", "/", -1)
	if !Exists(sourcePath) {
		return errors.New("找不到文件" + sourcePath)
	}
	targetPath = path.Clean(targetPath)
	fr, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer fr.Close()
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return err
	}
	defer gr.Close()
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(targetFile)
	writer.ReadFrom(gr)
	writer.Flush()
	return nil
}

func (dz *DataZip) Tar(src string, dstTar string, hasSelf bool, isGZ bool) error {
	src = strings.Replace(src, "\\", "/", -1)
	if !Exists(src) {
		return errors.New("找不到文件" + src)
	}
	fw, er := os.Create(dstTar)
	if er != nil {
		return er
	}
	defer fw.Close()

	var tw *tar.Writer
	if isGZ {
		gw := gzip.NewWriter(fw)
		defer gw.Close()

		tw = tar.NewWriter(gw)
	} else {
		tw = tar.NewWriter(fw)
	}
	defer tw.Close()

	last := len(src) - 1
	if src[last] != os.PathSeparator {
		src += string(os.PathSeparator)
	}
	src = strings.Replace(src, "\\", "/", -1)
	prefix := ""
	srci, err := os.Stat(src)
	if err != nil {
		return err
	}
	if srci.IsDir() {
		fis, er := ioutil.ReadDir(src)
		if er != nil {
			return er
		}
		if hasSelf {
			prefix = srci.Name()
		}
		for _, fi := range fis {
			err = dz.tarFile(path.Join(src, fi.Name()), prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		err = dz.tarFile(src, prefix, tw)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dz *DataZip) tarFile(filePath, prefix string, tw *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if len(prefix) != 0 {
		prefix = prefix + "/" + info.Name()
	} else {
		prefix = info.Name()
	}
	if info.IsDir() {
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			err = dz.tarFile(path.Join(filePath, fi.Name()), prefix, tw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = prefix
		err = tw.WriteHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dz *DataZip) UnTar(srcTar string, dstDir string, isGZ bool) error {
	srcTar = strings.Replace(srcTar, "\\", "/", -1)
	if !Exists(srcTar) {
		return errors.New("找不到文件" + srcTar)
	}
	dstDir = path.Clean(dstDir) + string(os.PathSeparator)
	fr, er := os.Open(srcTar)
	if er != nil {
		return er
	}
	defer fr.Close()
	var tr *tar.Reader
	if isGZ {
		gr, er := gzip.NewReader(fr)
		if er != nil {
			return er
		}
		defer gr.Close()
		tr = tar.NewReader(gr)
	} else {
		tr = tar.NewReader(fr)
	}

	for hdr, er := tr.Next(); er != io.EOF; hdr, er = tr.Next() {
		if er != nil {
			return er
		}
		fi := hdr.FileInfo()
		dstFullPath := dstDir + hdr.Name
		if hdr.Typeflag == tar.TypeDir {
			os.MkdirAll(dstFullPath, fi.Mode().Perm())
			os.Chmod(dstFullPath, fi.Mode().Perm())
		} else {
			os.MkdirAll(path.Dir(dstFullPath), os.ModePerm)
			if er := dz.unTarFile(dstFullPath, tr); er != nil {
				return er
			}
			os.Chmod(dstFullPath, fi.Mode().Perm())
		}
	}
	return nil
}

func (dz *DataZip) unTarFile(dstFile string, tr *tar.Reader) error {
	fw, er := os.Create(dstFile)
	if er != nil {
		return er
	}
	defer fw.Close()
	_, er = io.Copy(fw, tr)
	if er != nil {
		return er
	}
	return nil
}

func (dz *DataZip) Zipit(source, target string, hasSelf bool) error {
	source = strings.Replace(source, "\\", "/", -1)
	if !Exists(source) {
		return errors.New("找不到文件" + source)
	}
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	sinfo, err := os.Stat(source)
	if err != nil {
		return nil
	}
	last := len(source) - 1
	if source[last] != os.PathSeparator {
		source += string(os.PathSeparator)
	}
	source = strings.Replace(source, "\\", "/", -1)
	var baseDir string
	if sinfo.IsDir() && hasSelf {
		baseDir = filepath.Base(source)
	}
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		path = strings.Replace(path, "\\", "/", -1)
		header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})
	return err
}

func (dz *DataZip) Unzip(archive, target string) error {
	archive = strings.Replace(archive, "\\", "/", -1)
	if !Exists(archive) {
		return errors.New("找不到文件" + archive)
	}
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}
	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}
		fileReader, err := file.Open()
		if err != nil {
			if fileReader != nil {
				fileReader.Close()
			}
			return err
		}
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			if targetFile != nil {
				targetFile.Close()
			}
			return err
		}
		if _, err := io.Copy(targetFile, fileReader); err != nil {
			fileReader.Close()
			targetFile.Close()
			return err
		}
		fileReader.Close()
		targetFile.Close()
	}
	return nil
}

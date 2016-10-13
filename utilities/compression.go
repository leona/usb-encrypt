package uti

import (
    "archive/tar"
    "compress/gzip"
    "fmt"
    "io"
    "os"
    "strings"
    "path"
)


type Compressor struct {
    TarWriter *tar.Writer
    pathCurrent string
}


func (self *Compressor) Init(inputPath string, destinationPath string) {
    destinationFile, err := os.Create(destinationPath)
    
    if err != nil {
        panic(err.Error())
    }
    
    defer destinationFile.Close()
    
    self.pathCurrent = GetCurrentPath()
    
    var fileWriter io.WriteCloser = destinationFile
    
    tarfile, err := os.Create(destinationPath)
    
    if err != nil {
        panic(err.Error())
    }
    
    defer tarfile.Close()
    
    if strings.HasSuffix(destinationPath, ".gz") {
        fileWriter = gzip.NewWriter(tarfile) 
        defer fileWriter.Close()    
    }
    
    self.TarWriter = tar.NewWriter(fileWriter)
    defer self.TarWriter.Close()
    
    self.pathWalker(inputPath, inputPath)
}

func (self *Compressor) pathWalker(inputPath string, startPath string) {
    pathWorking, err := os.Open(inputPath)

    if err != nil {
        panic(err.Error())
    }
    
    defer pathWorking.Close()
    
    files, err := pathWorking.Readdir(0)  
    
    if err != nil {
        panic(err.Error())
    }
    
    for _, fileWorking := range files {
        if fileWorking.IsDir() {
            self.pathWalker(path.Join(inputPath, fileWorking.Name()), startPath)
            continue
        }
        
        file, err := os.Open(path.Join(pathWorking.Name(), fileWorking.Name()))
        
        if err != nil {
            panic(err.Error())
        }
    
        defer file.Close()
        
        header := new(tar.Header)
        
        header.Name     = strings.Replace(strings.Replace(file.Name(), `\`, "/", -1), BasePath(startPath) + "/", "", -1)
        header.Mode     = int64(fileWorking.Mode())
        header.ModTime  = fileWorking.ModTime()
        header.Typeflag = tar.TypeReg
        header.Size     = fileWorking.Size()

        fmt.Println("Compressing file: ", strings.Replace(file.Name(), self.pathCurrent, "", -1))

        if err := self.TarWriter.WriteHeader(header); err != nil {
            panic(err.Error())
        }
        
        if _, err := io.Copy(self.TarWriter, file); err != nil {
            panic(err.Error())
        }
    }
}


func Decompress(inputPath string) {
    workingFile, err := os.Open(inputPath)
    
    if err != nil {
        panic(err.Error())
    }
    
    defer workingFile.Close()
    
    pathCurrent := BasePath(inputPath)
    fmt.Println("Work path: " + pathCurrent)
    var fileReader io.ReadCloser = workingFile
    
    if strings.HasSuffix(inputPath, ".gz") {
        if fileReader, err = gzip.NewReader(workingFile); err != nil {
            panic(err.Error())
        }
        
        defer fileReader.Close()
    }
    
    tarBallReader := tar.NewReader(fileReader)
    
    for {
        header, err := tarBallReader.Next()
        
        if err != nil {
            if err == io.EOF {
                break
            }
            
            panic(err.Error())
        }
        
        fileName := path.Join(pathCurrent, header.Name)
        filePath := BasePath(fileName)

        if err := os.MkdirAll(filePath, os.FileMode(0755)); err != nil {
            panic(err.Error())
        }
            
        if header.Typeflag == tar.TypeReg {
            fmt.Println("Untarring :", fileName)
            writer, err := os.Create(fileName)
            
            if err != nil {
                panic(err.Error())
            }
            
            io.Copy(writer, tarBallReader)
            
            if err := os.Chmod(fileName, os.FileMode(header.Mode)); err != nil {
                panic(err.Error())
            }
            
            writer.Close()
        } else {
            fmt.Printf("Unable to untar type : %c in file %s", header.Typeflag, fileName)
        }
    }
}
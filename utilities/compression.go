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
    tarWriter *tar.Writer
    pathInput string
    pathCurrent string
    pathDestination string  
    pathTree []fileData
    destinationFile *os.File
    gzipWriter *gzip.Writer
}


func (self *Compressor) Init(inputPath string, destinationPath string) {
    self.pathCurrent = GetCurrentPath()
    self.pathDestination = destinationPath
    self.pathInput = inputPath
    
    self.createTarballHandler()

    self.pathTree = GetPathTree(self.pathInput)
    self.compressTree()
    
    defer self.destinationFile.Close()
    defer self.gzipWriter.Close()   
    defer self.tarWriter.Close()
}

func (self *Compressor) compressTree() {
    for _, file := range self.pathTree {
        self.addTarballItem(file)
    }
}

func (self *Compressor) addTarballItem(file fileData) {
    fileHandler, err := os.Open(file.path)
    
    if err != nil {
        panic(err)
    }
    
    header := new(tar.Header)
    header.Name     = strings.Replace(strings.Replace(file.path, `\`, "/", -1), BasePath(self.pathInput) + "/", "", -1)
    header.Mode     = int64(file.info.Mode())
    header.ModTime  = file.info.ModTime()
    header.Typeflag = tar.TypeReg
    header.Size     = file.info.Size()
    
    fmt.Println("Compressing file: ", strings.Replace(strings.Replace(file.path, `\`, "/", -1), BasePath(self.pathInput) + "/", "", -1))

    if err := self.tarWriter.WriteHeader(header); err != nil {
        panic(err.Error())
    }
    
    if _, err := io.Copy(self.tarWriter, fileHandler); err != nil {
        panic(err.Error())
    }
} 

func (self *Compressor) createTarballHandler() {
    var err error
    self.destinationFile, err = os.Create(self.pathDestination)
    
    if err != nil {
        panic(err.Error())
    }
    
    self.gzipWriter = gzip.NewWriter(self.destinationFile) 
    self.tarWriter = tar.NewWriter(self.gzipWriter)
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
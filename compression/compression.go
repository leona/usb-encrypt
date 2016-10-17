package compression

import (
    "archive/tar"
    "compress/gzip"
    "path/filepath"
    "fmt"
    "io"
    "os"
    "sync"
    "strings"
    "path"
    "github.com/neoh/usb-encrypt/uti"
)

type Handler struct {
    tarWriter *tar.Writer
    pathInput string
    pathCurrent string
    pathDestination string  
    destinationFile *os.File
    gzipWriter *gzip.Writer
    jobIteration int
    pathTreeLength int
    readJobs chan fileData
    writeJobs chan writeData
}

type writeData struct {
    fileHandler *os.File
    header *tar.Header
}

type fileData struct {
    Path string
    Info os.FileInfo
}

const maxJobs = 100000000
const maxProcesses = 1
var wg sync.WaitGroup

func (self *Handler) Init(inputPath string, destinationPath string) {
    self.readJobs = make(chan fileData, maxJobs)
    self.writeJobs = make(chan writeData)
    
    self.pathCurrent = uti.GetCurrentPath()
    self.pathDestination = destinationPath
    self.pathInput = inputPath

    wg.Add(1)
    go self.walkPath()
    
    self.createTarballHandler()
    self.startWorkers()
    
    wg.Wait()

    defer self.destinationFile.Close()
    defer self.gzipWriter.Close()   
    defer self.tarWriter.Close()
}

func (self *Handler) walkPath() {
    filepath.Walk(self.pathInput, func(path string, file os.FileInfo, err error) error {
        if !file.IsDir() {
            wg.Add(1)
            self.readJobs <- fileData{ path, file }
        }
            
        return nil
    })
    
    wg.Done()
}

func (self *Handler) startWorkers() {
    self.jobIteration = maxProcesses
    
    go func() {
        for job := range self.writeJobs {
            self.writeTarHeader(job.header, job.fileHandler)
            wg.Done()
        }
    }()
    
    for id := 0; id < maxProcesses; id++ {
        go self.worker(id)
    }
}

func (self *Handler) worker(id int) {
    for job := range self.readJobs {
        self.addTarballItem(job)
    }
}

func (self *Handler) addTarballItem(file fileData) {
    fileHandler, err := os.Open(file.Path)
    
    if err != nil {
        panic(err)
    }
    
    header := new(tar.Header)
    header.Name     = strings.Replace(strings.Replace(file.Path, `\`, "/", -1), uti.BasePath(self.pathInput) + "/", "", -1)
    header.Mode     = int64(file.Info.Mode())
    header.ModTime  = file.Info.ModTime()
    header.Typeflag = tar.TypeReg
    header.Size     = file.Info.Size()
    
    self.writeJobs <- writeData{ fileHandler, header }
} 

func (self *Handler) writeTarHeader(header *tar.Header, fileHandler *os.File) {
    if err := self.tarWriter.WriteHeader(header); err != nil {
        fmt.Println("Header name:", header.Name, "mode:", header.Mode, "modtime:", header.ModTime, "typeflag:", header.Typeflag, "size:", header.Size)
        panic(err.Error())
    }
    
    if _, err := io.Copy(self.tarWriter, fileHandler); err != nil {
        panic(err.Error())
    }
}

func (self *Handler) createTarballHandler() {
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
    
    pathCurrent := uti.BasePath(inputPath)

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
        filePath := uti.BasePath(fileName)

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
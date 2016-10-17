package pathTreeBuilder

import (
    "os"
    "path"
    "sync"
)

var wg sync.WaitGroup
const maxProcesses = 5

type FileData struct {
    Path string
    Info os.FileInfo
}

func GetPathTree(inputPath string) []FileData {
    var tree []FileData
    
    jobs := make(chan string, 100000000)
    results := make(chan FileData, 100000000)

    for w := 1; w <= maxProcesses; w++ {
        go worker(w, jobs, results)
    }
    
    go func() {
        for result := range results {
            tree = append(tree, result)
        }
    }()
    
    wg.Add(1)
    jobs <- inputPath
    
    
    wg.Wait()
    close(jobs)
    close(results)

    return tree
}

func worker(id int, jobs chan string, results chan<- FileData) {
    for job := range jobs {
        result := scanDir(job)
        
        for _, file := range result {
            name := path.Join(job, file.Name())

            if file.IsDir() {
                wg.Add(1)
                jobs <- name
            } else {
                results <- FileData{ name, file }
            }
        }
        
        wg.Done()
    }
}

func scanDir(input string) []os.FileInfo {
    pathWorking, err := os.Open(input)
    
    if err != nil {
        panic(err.Error())
    }
    
    defer pathWorking.Close()
    
    files, err := pathWorking.Readdir(0)  
    
    if err != nil {
        panic(err.Error())
    }
    
    return files
}
package uti

import (
    "os"
    "path"
    "sync"
)

var wg sync.WaitGroup

type fileData struct {
    path string
    info os.FileInfo
}

func GetPathTree(inputPath string) []fileData {
    const maxProcesses = 5
    var tree []fileData
    
    jobs := make(chan string, 100000000)
    results := make(chan fileData, 100000000)

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

func worker(id int, jobs chan string, results chan<- fileData) {
    for job := range jobs {
        result := scanDir(job)
        
        for _, file := range result {
            name := path.Join(job, file.Name())

            if file.IsDir() {
                wg.Add(1)
                jobs <- name
            } else {
                results <- fileData{ name, file }
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
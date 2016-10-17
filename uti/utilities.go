package uti

import(
    "os"
    "fmt"
    "time"
    "bufio"
    "strings"
    "encoding/hex"
    "crypto/md5"
    "strconv"
    "path/filepath"
)

func cleanInput(input string) string {
    return strings.Replace(strings.Replace(input, "\n", "", -1), "\r", "", -1)
}

func PathExists(path string) bool {
    _, err := os.Stat(path)
    if err == nil { return true }
    if os.IsNotExist(err) { return false }
    return true
}

func ExitPrompt(msg string) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println(msg, " Press any key to exit.")
    reader.ReadString('\n')
}

func ContinuePrompt(msg string) {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println(msg)
    reader.ReadString('\n')
}

func TakeInput(msg string) string {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print(msg)
    value, _ := reader.ReadString('\n')
    
    return cleanInput(value)
}

func GetCurrentPath() string {
    currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    
    if err != nil {
        panic(err)
    }
    
    return currentDir
}


func BasePath(path string) string {
    path = strings.Replace(path, `\`, "/", -1)
    split := strings.Split(path, `/`)
    
    return strings.Replace(path, "/" + split[len(split) - 1], "", -1)
}

func GetMD5(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func GetDrives() (drives []string){
    for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ"{
        _, err := os.Open(string(drive)+":\\")
        if err == nil {
            drives = append(drives, string(drive))
        }
    }
    return
}

func PromptDriveSelection(msg string) string {
    drives := GetDrives()
    fmt.Println("Available drives: ")
    
    for index, drive := range drives {
        fmt.Println(drive, "- [", index, "]")
    }
    
    driveIndexInput := TakeInput(msg)
    driveIndex, err := strconv.Atoi(driveIndexInput)
    
    if err != nil {
        panic(err)
    }
    
    return drives[driveIndex] + ":"
}

func GetUnix() int32 {
    return int32(time.Now().Unix())
}
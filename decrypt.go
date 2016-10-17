package main

import(
	"fmt"    
    "os"
    "os/exec"
    "path"
    "strings"
    "runtime"
    "github.com/neoh/usb-encrypt/uti"
    "github.com/neoh/usb-encrypt/encryption"
    "github.com/neoh/usb-encrypt/compression"
)

const vaultDir = "vault"

func main() {
    currentDir := uti.GetCurrentPath()
    vaultName  := uti.TakeInput("Enter vault name: ")
    
    if len(vaultName) < 1 {
        panic("Vault name too short")
    }
    
    vaultTarPath := path.Join(currentDir, vaultDir, vaultName + ".aes")
    
    if uti.PathExists(vaultTarPath) == false {
        uti.ExitPrompt("Vault does not exit.")
        return
    }
    
    vaultKey := uti.TakeInput("Enter key: ")
    
    if keyLength := len(vaultKey); keyLength < 4 {
        fmt.Println("Key length: ", keyLength)
        
        panic("Key must be at least 4 characters.")
    }
    
    var endPath string
    
    if runtime.GOOS == "windows" {
        endPath = uti.PromptDriveSelection("Select a drive number to extract files: ") + "/vault"
    } else {
        endPath = path.Join(uti.TakeInput("Output directory:"), vaultDir)
        
        if endPath == "." || endPath == "./" {
            panic("Error. Do not extract into USB as residual files can be recovered.")
        }
    }
    
    mountVault(vaultName, vaultKey, endPath)
}

func mountVault(vaultName string, vaultKey string, endPath string) {
    pathCurrent := uti.GetCurrentPath()
    
    decryptFile := path.Join(pathCurrent, vaultDir, vaultName + ".aes")
    outputTarFile := path.Join(endPath, vaultName + ".tar.gz")
    fmt.Println("Decrypting tarball")
    
    if err := os.MkdirAll(endPath, os.FileMode(0755)); err != nil {
        panic(err.Error())
    }
            
    if _, err := encryption.Decrypt(decryptFile, vaultKey, outputTarFile); err != nil {
        panic(err.Error())
    }

    fmt.Println("Decompressing tarball")
    compression.Decompress(outputTarFile)
    defer os.Remove(outputTarFile)
    
    osSpecificActions(endPath)
    fmt.Println("Finished")
}

func osSpecificActions(inputPath string) {
    inputPath = strings.Replace(inputPath, "/", `\`, -1)
    fmt.Println("Extracted path: ", inputPath)
    
    switch runtime.GOOS {
        case "windows":
            exec.Command("explorer", inputPath)
        default: break;
    }
}
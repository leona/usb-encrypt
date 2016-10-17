package main

import(
    "os"
    "fmt"
    "path"
    "github.com/neoh/usb-encrypt/utilities"
)

const vaultDir = "vault"

func main() {
    currentDir := uti.GetCurrentPath()
    vaultName := uti.TakeInput("Enter vault name: ") // validate alphanumeric
    rootVault := path.Join(currentDir, vaultDir)
    vaultPath :=uti. TakeInput("Enter directory to encrypt from: ")
    vaultTarPath := vaultPath + ".tar.gz"
    encryptedOutput := path.Join(rootVault, vaultName + ".aes")
    
    if len(vaultPath) > len(currentDir) {
        if vaultPath[0:len(currentDir)] == currentDir {
            panic("Do not encrypt files stored on the current USB as residual files can be recovered.")
        }
    }
    
    vaultKey := uti.TakeInput("Enter key: ")
    
    if len(vaultKey) < 1 {
        uti.ExitPrompt("Error key too short.")
        return
    }
    
    if err := os.MkdirAll(rootVault, os.FileMode(0755)); err != nil {
        panic(err.Error())
    }
    
    handler := uti.Compressor{}
    handler.Init(vaultPath, vaultTarPath)
    //defer os.Remove(vaultTarPath)
    
    fmt.Println("Encrypting: ", vaultTarPath)
    
    uti.Crypt(vaultTarPath, vaultKey, encryptedOutput)
    
    fmt.Println("Encrypted output:", encryptedOutput)
    fmt.Println("Finished")
}
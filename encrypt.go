package main

import(
    "os"
    "fmt"
    "path"
    "regexp"
    "github.com/neoh/usb-encrypt/uti"
    "github.com/neoh/usb-encrypt/encryption"
    "github.com/neoh/usb-encrypt/compression"
    "github.com/neoh/usb-encrypt/values"
)
func main() {
    var validatorName = regexp.MustCompile(values.patternValidateName)
    var validatorDir = regexp.MustCompile(values.patternValidateDir)
    
    currentDir := uti.GetCurrentPath()
    vaultName := uti.TakeInput(values.MessageVaultName)
    
    if !validatorName.MatchString(vaultName) {
        panic(values.ErrorVaultName)
    }
    
    rootVault := path.Join(currentDir, values.vaultDir)
    vaultPath := uti.TakeInput(values.MessageVaultPath)
    
    if !validatorDir.MatchString(vaultPath) {
        panic(values.ErrorVaultPath)
    }
    
    vaultTarPath := vaultPath
    encryptedOutput := path.Join(rootVault, vaultName + ".aes")
    
    if len(vaultPath) > len(currentDir) {
        if vaultPath[0:len(currentDir)] == currentDir {
            panic(values.ErrorEncryptFromUsb)
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
    
    timestamp := uti.GetUnix()
    
    handler := compression.Handler{}
    handler.Init(vaultPath, vaultTarPath)
    
    defer os.Remove(vaultTarPath)
    
    fmt.Println("Compression took", uti.GetUnix() - timestamp, "seconds")
    fmt.Println("Encrypting: ", vaultTarPath)
    
    encryption.Crypt(vaultTarPath, vaultKey, encryptedOutput)
    
    fmt.Println("Encrypted output:", encryptedOutput)
    fmt.Println("Finished")
}
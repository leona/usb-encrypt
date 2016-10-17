package encryption

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "errors"
    "io/ioutil"
    "io"
    "os"
    "bytes"
    "github.com/neoh/usb-encrypt/uti"

)

func Crypt(inputPath string, key string, outputPath string) {
    plaintext, err := ioutil.ReadFile(inputPath)
    
    if err != nil {
        panic(err.Error())
    }
    
    key = uti.GetMD5(key)
    
    block, err := aes.NewCipher([]byte(key))
    
    if err != nil {
        panic(err)
    }

    ciphertext := make([]byte, aes.BlockSize + len(plaintext))
    
    iv := ciphertext[:aes.BlockSize]
    
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    outputFile, err := os.Create(outputPath)

    if err != nil {
        panic(err.Error())
    }
    
    defer outputFile.Close()

    if _, err := io.Copy(outputFile, bytes.NewReader(ciphertext)); err != nil {
        panic(err.Error())
    }
}


func Decrypt(inputPath string, key string, outputPath string) (string, error) {
    text, err := ioutil.ReadFile(inputPath)
    
    if err != nil {
        panic(err.Error())
    }

    key = uti.GetMD5(key)
    
    block, err := aes.NewCipher([]byte(key))
    
    if err != nil {
        return "", err
    }
    
    if len(text) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }
    
    iv  := text[:aes.BlockSize]
    text = text[aes.BlockSize:]
    
    cfb := cipher.NewCFBDecrypter(block, iv)
    cfb.XORKeyStream(text, text)

    outputFile, err := os.Create(outputPath)

    if err != nil {
        panic(err.Error())
    }
    
    defer outputFile.Close()

    if _, err := io.Copy(outputFile, bytes.NewReader(text)); err != nil {
        panic(err.Error())
    }
    
    return "Success", nil
}
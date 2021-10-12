package auth

import (
	"io/ioutil"
	"strings"
    b64 "encoding/base64"
    "fmt"
    "bufio"
    "os"
    "syscall"
    "golang.org/x/term"
)


const FILENAME = "./credentials.txt"
func Encode(s string) (string) {
    sEnc := b64.StdEncoding.EncodeToString([]byte(s))
    return sEnc
}

func Decode(s string) string{
    b,_ := b64.StdEncoding.DecodeString(s)
    return string(b)
}

func Store(username string, password string)(error){
    b := []byte(username + ":" + Encode(password))
    return ioutil.WriteFile(FILENAME, b, 0644)
}

func Get()(string,string,error){
    b, _ := ioutil.ReadFile(FILENAME)
    strB := string(b)
    if len(strB) == 0 {
        return "","",nil
    }
    strArr := strings.Split(string(strB),":")
    return strArr[0],Decode(strArr[1]),nil
}

func Credentials() (string, string,error) {
    

    usernameIO,passwordIO,error := Get()
    if error != nil {
        return "","",error
    }
    if len(usernameIO) > 0 && len(passwordIO) > 0{
    	return strings.TrimSpace(usernameIO), strings.TrimSpace(passwordIO), nil 
    }
    reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter your Redmine username: \n")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Printf("Enter your Redmine password: \n")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	password := string(bytePassword)
	fmt.Printf("Do you want to store credentials in file system(Y/n)?\n")
    answer,err := reader.ReadString('\n')
    if strings.Compare(strings.TrimSpace(answer),"n") == 0 {
    	return strings.TrimSpace(username), strings.TrimSpace(password), nil
    }
    Store(username,password)


	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

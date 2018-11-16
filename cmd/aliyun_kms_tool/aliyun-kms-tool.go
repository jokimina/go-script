package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/jokimina/go-script/pkg/aliyun"
	"os"
)

var (
	encrypt  string
	decrypt  string
	kmsKeyId string
)

func init() {
	flag.StringVar(&encrypt, "e", "", "encrypt")
	flag.StringVar(&decrypt, "d", "", "decrypt")
	flag.StringVar(&kmsKeyId, "id", "", "kms key id")
	flag.Usage = usage
}

func usage() {
	fmt.Printf(`
  encrypt:
    %[1]s -e xxxxxxxxx -id xxxxxx
  decrypt:
    %[1]s -d xxxxxxxxxx

`, os.Args[0])
	flag.PrintDefaults()
}

func Decrypt(kmsClient *kms.Client, cipherTextBlob string) (plaintText string) {
	decRequest := kms.CreateDecryptRequest()
	decRequest.SetScheme(requests.HTTPS)
	decRequest.CiphertextBlob = cipherTextBlob
	response, err := kmsClient.Decrypt(decRequest)
	if err != nil {
		panic(err)
	}
	return response.Plaintext
}

func Encrype(kmsClient *kms.Client, plaintText string, keyId string) (cipherTextBlob string) {
	encRequest := kms.CreateEncryptRequest()
	encRequest.SetScheme(requests.HTTPS)
	encRequest.KeyId = keyId
	encRequest.Plaintext = plaintText
	response, err := kmsClient.Encrypt(encRequest)
	if err != nil {
		panic(err)
	}
	return response.CiphertextBlob
}

func main() {
	flag.Parse()
	if (encrypt == "" && decrypt == "") || (encrypt != "" && decrypt != "") {
		flag.Usage()
		fmt.Println("decrypt or encrypt? ")
		os.Exit(1)
	}
	if decrypt == "" && (encrypt == "" && kmsKeyId == "") {
		flag.Usage()
		fmt.Println("encrypt need a kms key id !")
		os.Exit(1)
	}
	config := aliyun.NewConfig().GetConfigFromEnv()
	kmsClient, err := kms.NewClientWithAccessKey(config.RegionId, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		panic(err)
	}

	if encrypt != "" {
		fmt.Println(Encrype(kmsClient, encrypt, kmsKeyId))
	}
	if decrypt != "" {
		fmt.Println(Decrypt(kmsClient, decrypt))
	}
}

/*
 * Copyright (c) 2024 Yunshan Networks
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/deepflowio/deepflow/message/controller"
)

var CAMD5 string

func GenerateAesKey(input []byte) string {
	return fmt.Sprintf("%x", md5.Sum(input))
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	len := len(origData)
	unpadding := int(origData[len-1])
	if (len - unpadding) < 0 {
		return nil
	}
	return origData[:(len - unpadding)]
}

func AesEncrypt(origDataStr, keyStr string) (string, error) {
	key := []byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData := PKCS7Padding([]byte(origDataStr), blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func AesDecrypt(cryptedStr, keyStr string) (string, error) {
	key := []byte(keyStr)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	crypted, err := base64.StdEncoding.DecodeString(cryptedStr)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	if len(crypted)%blockSize != 0 {
		return "", errors.New("input is not encrypt key")
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	if origData == nil {
		return "", errors.New("encrypt key failed")
	}
	return string(origData), nil
}

func GetEncryptKey(controllerIP, grpcServerPort, key string) (string, error) {
	grpcServer := net.JoinHostPort(controllerIP, grpcServerPort)
	conn, err := grpc.Dial(grpcServer, grpc.WithInsecure())
	if err != nil {
		log.Error("create grpc connection faild:" + err.Error())
		return "", err
	}
	defer conn.Close()

	client := controller.NewControllerClient(conn)
	ret, err := client.GetEncryptKey(
		context.Background(),
		&controller.EncryptKeyRequest{Key: &key},
	)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return ret.GetEncryptKey(), nil
}

func EncryptSecretKey(secretKey string) (string, error) {
	caData, err := os.ReadFile(K8S_CA_CRT_PATH)
	if err != nil {
		log.Error(err)
		return "", err
	}
	aesKey := GenerateAesKey(caData)
	encryptSecretKey, err := AesEncrypt(secretKey, aesKey)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return encryptSecretKey, nil
}

func DecryptSecretKey(secretKey string) (string, error) {
	caData, err := os.ReadFile(K8S_CA_CRT_PATH)
	if err != nil {
		log.Error(err)
		return "", err
	}
	aesKey := GenerateAesKey(caData)
	decryptSecretKey, err := AesDecrypt(secretKey, aesKey)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return decryptSecretKey, nil
}

func GetLocalCAMD5() (string, error) {
	caData, err := os.ReadFile(K8S_CA_CRT_PATH)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return GenerateAesKey(caData), nil
}

func getCAMD5() string {
	caData, err := os.ReadFile(K8S_CA_CRT_PATH)
	if err != nil {
		log.Error(err)
		return ""
	}
	return GenerateAesKey(caData)
}

func GetCAMD5() string {
	return CAMD5
}

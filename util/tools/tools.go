package tools

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var DATE_FORMAT_COMPARE = "20060102150405000"

func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func PString(text string) *string {
	return &text
}
func PInt(value int) *int {
	return &value
}

// var commonIV = []byte("GEbJOVHUONrWInXe")
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func MakeTimestamp(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
func MakeTime(millisec int64) time.Time {
	return time.Unix(0, millisec*int64(time.Millisecond))
}

func Encrypt(plaintext []byte, key []byte, iv string) (ciphertext []byte, err error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	logrus.Debugf("%v", len(iv))

	nonce := make([]byte, 12)

	if iv == "" {
		_, err = io.ReadFull(rand.Reader, nonce)
		if err != nil {
			return nil, err
		}
	} else {
		nonce = []byte(iv[15:27])

	}

	//logrus.Printf("%v", len(iv[15:27]))
	//nonce := make([]byte, gcm.NonceSize())
	//nonce := commonIV
	/*
		_, err = io.ReadFull(rand.Reader, nonce)
		if err != nil {
			return nil, err
		}
	*/

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func EncryptBase64(plaintext string, key string) (ciphertext string, err error) {
	text, err := Encrypt([]byte(plaintext), []byte(key), "")
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(text), nil
}
func EncryptHex(plaintext string, key string) (ciphertext string, err error) {
	text, err := Encrypt([]byte(plaintext), []byte(key), "")
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(text), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func Decrypt(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < 12 {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:12],
		ciphertext[12:],
		nil,
	)
}

func DecryptBase64(ciphertext string, key string) (plaintext string, err error) {
	sDec, _ := base64.RawURLEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}

	text, err := Decrypt(sDec, []byte(key))
	if err != nil {
		return "", err
	}
	return string(text), nil
}
func DecryptHex(ciphertext string, key string) (plaintext string, err error) {
	sDec, _ := hex.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}

	text, err := Decrypt(sDec, []byte(key))
	if err != nil {
		return "", err
	}
	return string(text), nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetStructTag(object interface{}, attribute string, tagName string) string {
	field, ok := reflect.TypeOf(object).Elem().FieldByName(attribute)
	if !ok {
		return "null"
	}
	return string(field.Tag.Get(tagName))
}
func GetStructTagJSON(object interface{}, attribute string) string {
	return GetStructTag(object, attribute, "json")
}

func RandString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func SetLastLash(text string) string {
	last := text[len(text)-1:]

	if last != "/" {
		return text + "/"

	}
	return text

}

func wildCardToRegexp(pattern string) string {
	components := strings.Split(pattern, "*")
	if len(components) == 1 {
		// if len is 1, there are no *'s, return exact match pattern
		return "^" + pattern + "$"
	}
	var result strings.Builder
	for i, literal := range components {

		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}

		// Quote any regular expression meta characters in the
		// literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return "^" + result.String() + "$"
}

func Match(pattern string, value string) bool {
	result, _ := regexp.MatchString(wildCardToRegexp(pattern), value)
	return result
}
func MatchList(list []string, value string) bool {
	for _, pattern := range list {
		if Match(pattern, value) {
			return true
		}
	}
	return false
}

package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/k-mistele/ccdc-scoreserver/lib/constants"
	"io/ioutil"
	"net/http"
	"strings"
)

// HELPER FUNCTION: GET TOKEN
func minioGetToken(s *Service) (string, error) {
	var url, dataStr, token string
	var data []byte
	var err error
	var request *http.Request
	var client *http.Client
	var response *http.Response
	var exists bool

	url = fmt.Sprintf("http://%s:%d/minio/webrpc", s.Host, s.Port)

	// FORMAT JSON DATA FOR JSONRPC LOGIN
	dataStr =
		fmt.Sprintf(`{"id":1,"jsonrpc":"2.0","params":{"username":"%s","password":"%s"},"method":"web.Login"}`,
			s.Username, s.Password)

	// CAST TO BYTES
	data = []byte(dataStr)

	// BUILD THE REQUEST
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		s.CheckFailedWithError(err)
		return "", err
	}
	request.Header.Set("Accept", "*/*")
	request.Header.Set("User-Agent", constants.BrowserFriendlyUserAgent)
	request.Header.Set("Content-Type", "application/json")

	// SEND IT
	client = &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// PARSE THE BODY AND SERIALIZE IT TO A MAP STRUCTURE
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// GET STRING OF BODY
	bodyStr := string(body)
	// UNMARSHALL THE RESULT
	var decodedJSONData map[string]json.RawMessage
	err = json.Unmarshal([]byte(bodyStr), &decodedJSONData)
	if err != nil {
		s.CheckFailedWithError(err)
		return "", err
	}

	// PULL THE TOKEN OUT
	var resultData map[string]string
	err = json.Unmarshal(decodedJSONData["result"], &resultData)
	if err != nil {
		s.CheckFailedWithError(err)
		return "", err
	}
	token, exists = resultData["token"]
	if !exists {
		return "", errors.New("no token found in response")
	} else {
		return token, nil
	}

}

// HELPER FUNCTION: CREATE A BUCKET
func minioCreateBucket(s *Service, authToken string) (bool, error) {
	var err error
	var url, authorizationHeader, bodyString string
	var data []byte
	var request *http.Request
	var client *http.Client
	var response *http.Response

	url = fmt.Sprintf("http://%s:%d/minio/webrpc", s.Host, s.Port)
	data = []byte(`{"id":1,"jsonrpc":"2.0","params":{"bucketName":"scoreserver-bucket-test"},"method":"web.MakeBucket"}`)
	authorizationHeader = fmt.Sprintf("Bearer %s", authToken)

	// BUILD THE REQUEST
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return false, nil
	}
	request.Header.Set("Accept", "*/*")
	request.Header.Set("User-Agent", constants.BrowserFriendlyUserAgent)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", authorizationHeader)

	// SEND IT
	client = &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	// PARSE THE BODY INTO A STRING
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, nil
	}
	bodyString = string(body)

	if strings.Contains(bodyString, "error") {
		errorMsg := fmt.Sprintf("failed to create bucket: %s", bodyString)

		log.Debugf("Bucket could not be created!")
		return false, errors.New(errorMsg)

	} else {
		log.Debugf("Bucket created!")
		return true, nil
	}

}

// HELPER FUNCTION: DELETE A BUCKET
func minioDeleteBucket(s *Service, authToken string) (bool, error) {
	var err error
	var url, authorizationHeader, bodyString string
	var data []byte
	var request *http.Request
	var client *http.Client
	var response *http.Response

	url = fmt.Sprintf("http://%s:%d/minio/webrpc", s.Host, s.Port)
	data =
		[]byte(`{"id":1,"jsonrpc":"2.0","params":{"bucketName":"scoreserver-bucket-test"},"method":"web.DeleteBucket"}`)
	authorizationHeader = fmt.Sprintf("Bearer %s", authToken)

	// BUILD THE REQUEST
	request, err = http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return false, nil
	}
	request.Header.Set("Accept", "*/*")
	request.Header.Set("User-Agent", constants.BrowserFriendlyUserAgent)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", authorizationHeader)

	// SEND IT
	client = &http.Client{}
	response, err = client.Do(request)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	// PARSE THE BODY INTO A STRING
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, nil
	}
	bodyString = string(body)

	// CHECK FOR AN ERROR
	if strings.Contains(bodyString, "error") {
		errorMsg := fmt.Sprintf("failed to create bucket: %s", bodyString)
		log.Debugf("Bucket could not be deleted!")
		return false, errors.New(errorMsg)

	} else {
		log.Debugf("Bucket deleted!")
		return true, nil
	}
}

// TYPE ServiceCheck
// PARAMS: LOL NONE,
func MinioLoginCheck(s *Service) (bool, error) {
	var err error

	// LOG IN BY GETTING A TOKEN
	_, err = minioGetToken(s)
	if err != nil {
		s.CheckFailedWithError(err)
		return false, err
	} else {
		s.CheckPassed()
		return true, nil
	}
}

// TYPE ServiceCheck
// PARAMS: None
func MinioBucketCheck(s *Service) (bool, error) {
	var err error
	var bucketCreated, bucketDeleted bool

	token, err := minioGetToken(s)
	if err != nil {
		s.CheckFailedWithReason("Authentication failed!")
		return false, err
	}

	// CREATE A BUCKET
	bucketCreated, err = minioCreateBucket(s, token)
	if err != nil {
		s.CheckFailedWithError(err)
		return bucketCreated, err
	}

	// DELETE THE BUCKET (CLEAN UP)
	bucketDeleted, err = minioDeleteBucket(s, token)
	if err != nil {
		s.CheckFailedWithError(err)
		return bucketDeleted, err
	}

	s.CheckPassed()
	return bucketDeleted, nil
}

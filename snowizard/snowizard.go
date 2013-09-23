// Package snowizard provides a client for the snowizard id generation server.
package snowizard

import (
	proto "code.google.com/p/goprotobuf/proto"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	// Content-Type for text/plain
	TextContentType = "text/plain"

	// Content-Type for application/json
	JsonContentType = "application/json"

	//Content-Type for application/x-protobuf
	ProtobufContentType = "application/x-protobuf"

	// User-Agent sent to the server
	UserAgent = "gosnowizard"
)

var (
	// ErrNoServers means that all hosts failed, or were not available
	ErrNoServers = errors.New("snowizard: no servers configured or available")

	// ErrMalformedResp means that the response could not be properly parsed
	ErrMalformedResp = errors.New("snowizard: server error")
)

// Snowizard is the http client
type SnowizardClient struct {
	Hosts        []string
	client       *http.Client
	decodeRespFn func([]byte) (int64, error)
	contentType  string
	timeoutMs    int
}

// NewSnowizardTextClient returns a Snowizard client using the provided hosts
// and sets content-type to text/plain and decodeRespFn to decodeText
func NewSnowizardTextClient(hosts []string, connectTimeout time.Duration) *SnowizardClient {

	var transport http.Transport = http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, connectTimeout)
		},
	}

	return &SnowizardClient{
		Hosts:        hosts,
		client:       &http.Client{Transport: &transport},
		decodeRespFn: decodeText,
		contentType:  TextContentType,
	}
}

// NewSnowizardJsonClient returns a Snowizard client using the provided hosts
// and sets content-type to application/json and decodeRespFn to decodeJson
func NewSnowizardJsonClient(hosts []string, connectTimeout time.Duration) *SnowizardClient {

	var transport http.Transport = http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, connectTimeout)
		},
	}

	return &SnowizardClient{
		Hosts:        hosts,
		client:       &http.Client{Transport: &transport},
		decodeRespFn: decodeJson,
		contentType:  JsonContentType,
	}
}

// NewSnowizardProtobufClient returns a Snowizard client using the provided hosts
// and sets content-type to application/x-protobuf and decodeRespFn to decodeProtobuf
func NewSnowizardProtobufClient(hosts []string, connectTimeout time.Duration) *SnowizardClient {

	var transport http.Transport = http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, connectTimeout)
		},
	}

	return &SnowizardClient{
		Hosts:        hosts,
		client:       &http.Client{Transport: &transport},
		decodeRespFn: decodeProtobuf,
		contentType:  ProtobufContentType,
	}
}

// decodeText decodes text response and returns id
// id = -1 if response can not be parsed
func decodeText(in []byte) (int64, error) {
	res, err := strconv.ParseInt(string(in), 0, 64)
	if err != nil {
		log.Println(err)
		return -1, ErrMalformedResp
	}

	return res, nil
}

// decodeJson decodes json response and returns id
// id = -1 if response can not be parsed
func decodeJson(in []byte) (int64, error) {
	resp := &SnowizardResponse{}
	err := json.Unmarshal(in, resp)
	if err != nil {
		log.Printf("snowizard: unmarshall json: %v", err)
		return -1, ErrMalformedResp
	}

	return resp.GetId(), nil
}

// decodeProtobuf decodes protobuf response and returns id
// id = -1 if response can not be parsed
func decodeProtobuf(in []byte) (int64, error) {

	resp := &SnowizardResponse{}
	err := proto.Unmarshal(in, resp)
	if err != nil {
		log.Printf("snowizard: proto unmarshal error: %v", err)
		return -1, ErrMalformedResp
	}
	return resp.GetId(), nil
}

// Next iterates through all hosts and tries to parse the response
// if all hosts fail, ErrNoServers is returned
func (s SnowizardClient) Next() (int64, error) {

	for _, host := range s.Hosts {

		url := fmt.Sprintf("http://%s", host)
		req, err := http.NewRequest("GET", url, nil)

		req.Header.Add("Content-Type", s.contentType)
		req.Header.Add("User-Agent", UserAgent)

		resp, err := s.client.Do(req)
		if err != nil {
			log.Printf("snowizard: do: %v", err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			log.Printf("snowizard: status code: %d", resp.StatusCode)
			continue
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("snowizard: failed reading response body: %v", err)
			continue
		}

		id, err := s.decodeRespFn(body)
		if err != nil {
			log.Printf("snowizard: failed decoding response body: %v", err)
			continue
		}

		return id, nil

	}

	return -1, ErrNoServers
}

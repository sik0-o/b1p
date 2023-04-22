package proxy

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Data type is a struct to contents common proxy data
type Data struct {
	Type     string
	Address  string
	Port     int
	Login    string
	Password string
}

// NewData create a new proxy data from URL struct
func NewData(u *url.URL) (*Data, bool) {
	proxyPassword, ok := u.User.Password()
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		ok = false
	}

	return &Data{
		Type:     "http",
		Address:  u.Hostname(),
		Port:     port,
		Login:    u.User.Username(),
		Password: proxyPassword,
	}, ok
}

// ParseURL function parse proxy URL to proxy *Data
// function use standart url.Parse function and than
// create a new data
func ParseURL(proxyURL string) (*Data, error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, err
	}

	data, _ := NewData(u)

	return data, nil
}

func splitProxyStr(proxy string) []string {
	dSplit := strings.Split(proxy, "@")
	proxyArr := []string{}
	for _, arr := range dSplit {
		arrSplit := strings.Split(arr, ":")
		proxyArr = append(proxyArr, arrSplit...)
	}

	return proxyArr
}

// BaseStringToURL parse proxy baseString to URL.
// You need to trim a scheme, otherwise parsing will be broken.
// baseString formats:
// - login:password@address:port
// - login:password:address:port
func BaseStringToURL(proxy string) (*url.URL, error) {
	proxyArr := splitProxyStr(proxy)
	var host, port, login, password string

	switch len(proxyArr) {
	case 1:
		// only host
		host = proxyArr[0]
	case 2:
		// only host and port
		host = proxyArr[0]
		port = proxyArr[1]
	case 3:
		login = proxyArr[0]
		host = proxyArr[1]
		port = proxyArr[2]
	case 4:
		login = proxyArr[0]
		password = proxyArr[1]
		host = proxyArr[2]
		port = proxyArr[3]
	default:
		return nil, errors.New("Wrong arrLen size")
	}

	return &url.URL{
		Scheme: "http",
		User:   url.UserPassword(login, password),
		Host:   host + ":" + port,
	}, nil
}

// NewDataFromBaseString create a new proxy data from BaseString
func NewDataFromBaseString(proxy string) (*Data, error) {
	u, err := BaseStringToURL(proxy)
	if err != nil {
		return nil, err
	}

	p, _ := NewData(u)

	return p, nil
}

// BaseString return proxy base string for data
// return string format:
//
//	login:password@host:port
//
// basestring contains only address:port
// login:password - present only if login or passowrd is not empty
func (pd Data) BaseString() string {
	baseString := fmt.Sprintf("%s:%d", pd.Address, pd.Port)
	if pd.Login != "" || pd.Password != "" {
		baseString = fmt.Sprintf("%s:%s@%s", pd.Login, pd.Password, baseString)
	}

	return baseString
}

// URL return standart net/url URL struct generated from proxy data
func (pd Data) URL() *url.URL {
	return &url.URL{
		Scheme: pd.Type,
		User:   url.UserPassword(pd.Login, pd.Password),
		Host:   fmt.Sprintf("%s:%d", pd.Address, pd.Port),
	}
}

// ToString method returns proxy data
func (pd Data) ToString() string {
	return pd.URL().String()
}

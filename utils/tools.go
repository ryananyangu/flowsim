package utils

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	GORM_SEARCH_INPUT_COUNT  = 2
	RANGE_SEARCH_PARAM_COUNT = 2
)

type RequestOptions struct {
	Sensitive bool
}

func ReadFile(file string) (string, error) {
	openedFile, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer openedFile.Close()

	byteValue, _ := io.ReadAll(openedFile)

	return string(byteValue[:]), nil

}

func ISODateConversion(timestamp string) (date time.Time, err error) {
	date, err = time.Parse(time.RFC3339Nano, strings.TrimSpace(timestamp))
	if err != nil {
		logrus.Error(err)
		err = errors.New("invalid date")
		return
	}
	return
}

func ModelValidationResponse(err error) ([]map[string]string, error) {

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		invalid := make([]map[string]string, len(ve))
		for i, fe := range ve {
			invalid[i] = map[string]string{"field": fe.Field(), "tag": fe.Tag()}
		}
		return invalid, nil
	}

	return nil, err
}

func ExtractDateRange(date_range string) (start time.Time, end time.Time, err error) {
	dates := strings.Split(date_range, "to")
	if len(dates) != 2 {
		err = errors.New("invalid date range")
		return
	}

	start, err1 := ISODateConversion(dates[0])
	end, err = ISODateConversion(dates[1])
	if err1 != nil || err != nil {
		err = errors.New("invalid date format")
		return
	}
	return
}

func Request(request string, headers map[string][]string, urlPath, method string, opts RequestOptions) (string, error) {

	reqURL, _ := url.Parse(urlPath)

	reqBody := io.NopCloser(strings.NewReader(request))

	req := &http.Request{
		Method: method,
		URL:    reqURL,
		Header: headers,
		Body:   reqBody,
	}

	res, err := ExternalRequestTimer(req)
	if err != nil {
		logrus.Errorf("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | ERROR : %v", urlPath, method, request, err)
		return "", err
	}

	data, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	resbody := string(data)
	if !opts.Sensitive {
		logrus.Infof("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | STATUS : %s | HTTP_CODE : %d | RESPONSE : %s", urlPath, method, request, res.Status, res.StatusCode, resbody)
	}

	if res.StatusCode > 299 || res.StatusCode <= 199 {
		logrus.Errorf("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | STATUS : %s | HTTP_CODE : %d", urlPath, method, request, res.Status, res.StatusCode)
		return resbody, fmt.Errorf("%d", res.StatusCode)
	}

	// logrus.Infof("SEND REQUEST | URL : %s | METHOD : %s | BODY : %s | STATUS : %s | HTTP_CODE : %d", urlPath, method, resbody, res.Status, res.StatusCode)

	return resbody, nil
}

func ExternalRequestTimer(req *http.Request) (*http.Response, error) {

	var start, connect, dns, tlsHandshake time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) {
			dns = time.Now()
			logrus.Debug(dsi)
		},
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			logrus.Debug(ddi)
			logrus.Infof("DNS Done: %v", time.Since(dns))
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			// logrus.Debug(cs, err)
			logrus.Infof("TLS Handshake: %v", time.Since(tlsHandshake))
		},

		ConnectStart: func(network, addr string) {
			connect = time.Now()
			logrus.Debug(network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			logrus.Debug(network, addr, err)
			logrus.Infof("Connect time: %v", time.Since(connect))
		},

		GotFirstResponseByte: func() {
			logrus.Warnf("TAT : %v", time.Since(start))
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	// NOTE: Below line is to ignore ssl certificate
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return res, err
	}
	return res, nil
}

func GinPaginate(ctx *gin.Context) (page int, size int) {
	limitstr, is_limit := ctx.GetQuery("size")

	pagestr, is_page := ctx.GetQuery("page")

	if !is_limit || !is_page {
		limitstr = "10"
		pagestr = "1"
	}

	size, err := strconv.Atoi(limitstr)
	page, err1 := strconv.Atoi(pagestr)
	if err != nil || err1 != nil {
		page = 1
		size = 10
		return
	}

	return
}

func GormSearch(queryParams map[string][]string, query *gorm.DB) (q *gorm.DB, err error) {

	for name, param := range queryParams {
		value := strings.Split(name, "__")
		if len(value) != GORM_SEARCH_INPUT_COUNT {
			logrus.Errorf("invalid search format for [%s]", name)
			continue
		}
		columnOperation := value[0]
		columnName := value[1]

		switch columnOperation {
		case "ilike":
			query.Where(fmt.Sprintf("%s ILIKE ?", columnName), "%"+param[0]+"%")
		case "in":
			query.Where(fmt.Sprintf("%s IN (?)", columnName), strings.Split(param[0], ","))
		case "gte":
			query.Where(fmt.Sprintf("%s >= ?", columnName), param[0])
		case "lte":
			query.Where(fmt.Sprintf("%s <= ?", columnName), param[0])
		case "gt":
			query.Where(fmt.Sprintf("%s > ?", columnName), param[0])
		case "lt":
			query.Where(fmt.Sprintf("%s < ?", columnName), param[0])
		case "eq":
			query.Where(fmt.Sprintf("%s = ?", columnName), param[0])
		case "like":
			query.Where(fmt.Sprintf("%s LIKE ?", columnName), "%"+param[0]+"%")
		case "btwn":
			if len(param) != RANGE_SEARCH_PARAM_COUNT {
				err = fmt.Errorf("range search requires 2 values received %v", param)
				return
			}
			rangeBtwn := strings.Split(param[0], ",")
			query.Where(fmt.Sprintf("%s <= ? AND %s >= ?", columnName, columnName), rangeBtwn[0], rangeBtwn[1])
		default:
			logrus.Infof("failed to create search")

		}
	}
	q = query
	return
}
package common

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
)

func GetAzAccountInfo() (AzAccountShowOutput, error) {
	cmd := exec.Command("az", "account", "show")

	output, err := cmd.Output()
	if err != nil {
		return AzAccountShowOutput{}, err
	}

	var azAccountShowOutput AzAccountShowOutput
	err = json.Unmarshal(output, &azAccountShowOutput)
	if err != nil {
		log.Fatal(err)
	}

	return azAccountShowOutput, nil
}
func GetSubscriptionId() (string, error) {
	data, err := GetAzAccountInfo()
	if err != nil {
		return "", err
	}
	return data.Id, nil
}

type User struct {
	Name string `json:"name"`
}
type AzAccountShowOutput struct {
	Id               string `json:"id"`
	TenantId         string `json:"tenantId"`
	User             User   `json:"user"`
	SubscriptionName string `json:"name"`
}

func ExtractResourceGroup(s *string) (*string, error) {
	re := regexp.MustCompile(`/subscriptions/.*?/resourceGroups/(?P<rg>.*?)/`)
	matches := re.FindStringSubmatch(*s)
	if len(matches) == 0 {
		return nil, errors.New("no matches found")
	}
	rgIndex := re.SubexpIndex("rg")
	return &matches[rgIndex], nil
}

type IP struct {
	Query string
}

func GetIPAddress() (string, error) {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	var ip IP
	if err = json.Unmarshal(body, &ip); err != nil {
		return "", err
	}

	return ip.Query, nil
}

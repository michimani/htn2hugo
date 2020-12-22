package main

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
}

var (
	initialAtomLink string
	hatenaId        string
	hatenaAPIKey    string
	host            string
	saveDir         string = "./dist/"
)

func main() {
	next := initialAtomLink
	hasNext := true
	result := make(map[bool]int)

	for hasNext {
		xmlData, err := getXML(next)
		if err != nil {
			fmt.Println(err)
			break
		}

		atom := Atom{}
		xerr := xml.Unmarshal([]byte(xmlData), &atom)
		if xerr != nil {
			fmt.Println(xerr.Error())
			break
		}

		perr := atom.prepare()
		if perr != nil {
			fmt.Println(perr.Error())
			break
		}

		if host == "" {
			host = atom.Host
			fmt.Printf("Blog host: %s\n", host)
		}

		hasNext = atom.HasNext
		next = atom.NextPage
		for _, entry := range atom.Entries {
			res := entry.save()
			result[res]++
		}
	}

	fmt.Println("Finished!")
	fmt.Printf("[result] Success: %d, Failed: %d\n", result[true], result[false])
}

func init() {
	ready := true

	hatenaId = os.Getenv("HTN_ID")
	hatenaAPIKey = os.Getenv("HTN_API_KEY")
	if hatenaId == "" {
		fmt.Printf("The environment variable HTN_ID is not set.\n")
		ready = false
	}
	if hatenaAPIKey == "" {
		fmt.Printf("The environment variable HTN_API_KEY is not set.\n")
		ready = false
	}

	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		derr := os.MkdirAll(saveDir, 0777)
		if derr != nil {
			fmt.Println(err.Error())
			ready = false
		}
	}

	if !ready {
		os.Exit(1)
	}

	fmt.Printf("Hatena ID: %s\n", hatenaId)
	initialAtomLink = fmt.Sprintf("https://blog.hatena.ne.jp/%s/%s.hateblo.jp/atom/entry", hatenaId, hatenaId)
}

func getXML(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	var basicAuthZ string = generateBasicAuthZ()
	req.Header.Set("Authorization", "Basic "+basicAuthZ)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("Failed to get xml. [%s]", res.Status)
	}

	body, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return string(body), nil
}

func generateBasicAuthZ() string {
	raw := hatenaId + ":" + hatenaAPIKey
	bytes := []byte(raw)
	return base64.StdEncoding.EncodeToString(bytes)
}

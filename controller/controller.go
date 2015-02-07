package controller

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/mix3/plantumlor/plantuml"
)

type AppContext struct {
	PlantUML plantuml.PlantUML
}

func (c *AppContext) TransferHandler(w http.ResponseWriter, r *http.Request) {
	params := context.Get(r, "params").(httprouter.Params)

	str, err := inflate(params.ByName("data"))
	if err != nil {
		panic(err)
	}

	b, err := c.PlantUML.Transfer(str, plantuml.PNG)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(b)
}

func inflate(input string) (string, error) {
	str := strings.Replace(input, "_", "/", -1)

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	ret, err := ioutil.ReadAll(flate.NewReader(bytes.NewReader(data)))
	if err != nil {
		return "", err
	}

	return string(ret), nil
}

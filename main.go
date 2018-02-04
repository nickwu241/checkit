package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type Checklist struct {
	Name  string      `json:"name"`
	Items []CheckItem `json:"items"`
}

type CheckItem struct {
	Checked bool   `json:"checked"`
	Text    string `json:"text"`
}

type ChecklistRequestBody struct {
	Checklist []CheckItem `json:"checklist"`
	Item      CheckItem   `json:"item"`
	Index     int         `json:"index"`
	Name      string      `json:"name"`
}

var defaultItems = []CheckItem{
	CheckItem{
		false,
		"dishes",
	},
	CheckItem{
		false,
		"laundry",
	},
}

var masterList = Checklist{
	Name: "default",
	Items: []CheckItem{},
	// Items: defaultItems,
}

func main() {
	r := mux.NewRouter()
	n := negroni.New()
	l := negroni.NewLogger()
	l.SetDateFormat(time.Stamp)
	n.Use(l)
	n.Use(negroni.NewRecovery())
	n.UseHandler(r)

	r.HandleFunc("/checklist", checklistHandler).
		Methods("PUT", "GET", "POST", "DELETE")
	r.HandleFunc("/checklist/share", shareHandler).
		Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", n)
}

func checklistHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Println(r.Method)
	fmt.Println(r.URL.Path)
	switch r.Method {
	case http.MethodPut:
		reqBody, err := decode(r.Body)
		if err != nil {
			return
		}
		i := reqBody.Index
		masterList.Items[i].Checked = !masterList.Items[i].Checked
	case http.MethodGet:
		index := r.URL.Query().Get("index")
		name := r.URL.Query().Get("name")
		d := r.URL.Query().Get("default")
		if d != "" {
			masterList = firebaseGet("default")
			w.Write(encode(masterList))
			return
		}
		if name != "" {
			masterList = firebaseGet(name)
			w.Write(encode(masterList))
			return
		}

		i, err := strconv.Atoi(index)
		if index == "" || err != nil {
			w.Write(encode(masterList))
			return
		}
		w.Write(encode(masterList.Items[i]))
	case http.MethodPost:
		reqBody, err := decode(r.Body)
		if err != nil {
			w.Write([]byte("ERR"))
			return
		}
		if len(reqBody.Checklist) > 0 {
			masterList.Items = reqBody.Checklist
		} else {
			masterList.Items = append(masterList.Items, reqBody.Item)
		}
		w.Write([]byte("OK"))
	case http.MethodDelete:
		reqBody, err := decode(r.Body)
		if err != nil {
			return
		}
		i := reqBody.Index
		masterList.Items = masterList.Items[:i+copy(masterList.Items[i:], masterList.Items[i+1:])]
	default:
		fmt.Printf("method not supported %q\n", r.Method)
	}
}

func shareHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	type shareBody struct {
		Name string
	}
	b := shareBody{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&b)
	if err != nil {
		fmt.Println("cannot decode body")
	}
	masterList.Name = b.Name
	firebasePut("templates/"+b.Name, encode(masterList.Items))
	w.Write([]byte("OK"))
}

func encode(data interface{}) []byte {
	out, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("error encoding %v\n", data)
		return []byte{}
	}
	return out
}

func decode(reader io.ReadCloser) (ChecklistRequestBody, error) {
	var b ChecklistRequestBody
	d := json.NewDecoder(reader)
	err := d.Decode(&b)
	if err != nil {
		return ChecklistRequestBody{}, fmt.Errorf("cannot decode")
	}
	return b, nil
}

func firebasePut(fullName string, data []byte) {
	firebaseURL := fmt.Sprintf("https://checkit-bafbc.firebaseio.com/%s.json", fullName)
	req, err := http.NewRequest("PUT", firebaseURL, bytes.NewReader(data))
	if err != nil {
		fmt.Println("cannot create PUT request for firebase")
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error sending PUT request for firebase")
		return
	}
	fmt.Printf("PUT %q=%d\n", req.URL.Path, resp.StatusCode)
}

func firebaseGet(name string) Checklist {
	firebaseURL := fmt.Sprintf("https://checkit-bafbc.firebaseio.com/templates/%s.json", name)
	req, err := http.NewRequest("GET", firebaseURL, nil)
	if err != nil {
		fmt.Println("cannot create GET request for firebase")
		return Checklist{}
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error sending GET request for firebase")
		return Checklist{}
	}
	fmt.Printf("GET %q=%d\n", req.URL.Path, resp.StatusCode)
	var items []CheckItem
	d := json.NewDecoder(resp.Body)
	err = d.Decode(&items)
	if err != nil {
		fmt.Println("cannot decode")
		return Checklist{}
	}
	return Checklist{
		Name:  name,
		Items: items,
	}
}

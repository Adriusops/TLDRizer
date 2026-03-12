package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ollama_json struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollama_response struct {
	Response string `json:"response"`
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("POST /roast", getRoast)

	err := http.ListenAndServe(":7070", router)
	if err != nil {
		log.Fatal(err)
	}
}

func getRoast(w http.ResponseWriter, r *http.Request) {
	// save input
	input, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	if input == nil {
		http.Error(w, "Failed to call Ollama", http.StatusBadRequest)
		return
	}

	instance := new(ollama_json)
	instance.Model = "mistral"
	instance.Prompt = string(input)
	instance.Stream = false

	// convert with json.Marshal
	bs, _ := json.Marshal(instance)

	req, err := http.NewRequest("POST", "http://localhost:11434/api/generate", bytes.NewBuffer(bs))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var result ollama_response

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		http.Error(w, "Failed to Unmarshal", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(result.Response))
}

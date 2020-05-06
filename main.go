package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	channelTwelve = `#EXTINF:0 tvg-rec="0",12 Канал
#EXTGRP:новости
https://12channel.bonus-tv.ru/cdn/12channel/playlist.m3u8
`
	link         = "yourlinktoedem"
	updateDomain = "https://12channel.bonus-tv.ru/cdn/12channel/playlist.m3u8"
)

func updateDomainIP() {
	resp, err := http.Get(updateDomain)
	if err != nil {
		log.Println(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
}

func getList(returned *string) error {

	go updateDomainIP()

	resp, err := http.Get(link)
	if err != nil {
		return nil
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			*returned = ""
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	*returned = string(body) + channelTwelve
	return nil
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	var myList string
	err := getList(&myList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal_error"))
		return
	}
	fmt.Fprintf(w, myList)
	return
}

func main() {
	http.HandleFunc("/tv.m3u8", listHandler)
	http.ListenAndServe(":4040", nil)
}

package main

import (
	"fmt"
	"net/http"
)

func authenticationCheck(request http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		if contains(r.RemoteAddr) {
			fmt.Println("already blocked")
			forbiddenAuth(w)
			return
		}

		if auth != "" {
			if auth == authKey {
				request(w, r)
			} else {
				nonAuthRequest(r.RemoteAddr)
				fmt.Println("missmatch of key")
				forbiddenAuth(w)
			}
		} else {
			fmt.Println("no key available")
			forbiddenAuth(w)
		}
	})
}

func nonAuthRequest(ip string) {
	returnedIP := selectedIPStruct{}

	rows, err := db.Queryx("SELECT ip, tries, blocked FROM denylist WHERE ip = ?", ip)
	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		err := rows.StructScan(&returnedIP)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	rows.Close()

	if returnedIP.IP == "" {
		_, err := db.Query("INSERT INTO denylist (ip, tries, blocked) VALUES (?, 1, 0)", ip)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if returnedIP.Tries != 4 {
		_, err := db.Query("UPDATE denylist SET tries = ? WHERE ip = ?", returnedIP.Tries+1, ip)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		_, err := db.Query("UPDATE denylist SET tries = ?, blocked = 1 WHERE ip = ?", returnedIP.Tries+1, ip)
		if err != nil {
			fmt.Println(err)
			return
		}
		blocked = append(blocked, ip)
	}

}

func initBlockedIPs() {
	returnedIPS := []selectedIPSStruct{}
	err := db.Select(&returnedIPS, "SELECT ip FROM denylist WHERE blocked = 1")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, ip := range returnedIPS {
		blocked = append(blocked, ip.IP)
	}
}

func contains(value string) bool {
	for _, item := range blocked {
		if item == value {
			fmt.Println("blocked")
			return true
		}
	}
	return false
}

type selectedIPStruct struct {
	IP      string `db:"ip"`
	Tries   int    `db:"tries"`
	Blocked int    `db:"blocked"`
}

type selectedIPSStruct struct {
	IP string `db:"ip"`
}

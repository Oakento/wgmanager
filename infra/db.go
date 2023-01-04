package infra

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"wgmanager/util"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func InitDb(dbFile string) {
	_, err := os.Stat(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(dbFile)
			if err != nil {
				HandlePermissionOrFileSystemError(err)
			}
			defer f.Close()
			err = os.Chmod(dbFile, 0600)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			HandlePermissionOrFileSystemError(err)
		}
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		fmt.Println(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	Db = db
}

func CreateDb() {

	sqlStmt := `
	CREATE TABLE peer (
		username VARCHAR(128) PRIMARY KEY,
		public_key CHAR(44),
		private_key CHAR(44),
		ip_address CHAR(18)
	)
	`
	_, err := Db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func InsertDbPeer(peer *PeerEntity) {

	sqlStmt := "INSERT INTO peer(username, public_key, private_key, ip_address) values(?, ?, ?, ?)"

	_, err := Db.Exec(sqlStmt, peer.Username, peer.PublicKey, peer.PrivateKey, peer.IpAddress)
	if err != nil {
		fmt.Println("Error inserting data, probably permission not granted")
		os.Exit(1)
	}

}

func RemoveDbPeer(username string) string {

	sqlStmt := "DELETE FROM peer WHERE username=?"

	_, err := Db.Exec(sqlStmt, username)
	if err != nil {
		fmt.Println("Error operating data, probably permission not granted")
		os.Exit(1)
	}

	return username
}

func UpdateDbPeer(peer *PeerEntity) {
	if peer.PublicKey == "" && peer.PrivateKey == "" && peer.IpAddress == "" {
		return
	}
	params := make([]any, 0, 4)
	sqlStmt := "UPDATE peer SET "
	if peer.PublicKey != "" && peer.PrivateKey != "" {
		sqlStmt = util.StringConcat(sqlStmt, "public_key=?, private_key=?")
		params = append(params, peer.PublicKey, peer.PublicKey)
	}
	if peer.IpAddress != "" {
		sqlStmt = util.StringConcat(sqlStmt, ", ip_address=?")
		params = append(params, peer.IpAddress)
	}
	sqlStmt = util.StringConcat(sqlStmt, " WHERE username=?")
	params = append(params, peer.Username)
	_, err := Db.Exec(sqlStmt, params...)
	if err != nil {
		log.Fatal(err)
	}
}

func SelectDbPeer(username string) *PeerEntity {
	sqlStmt := "SELECT * FROM peer WHERE username=?"
	rows, err := Db.Query(sqlStmt, username)
	if err != nil {
		fmt.Println("Error querying data, probably permission not granted")
		os.Exit(1)
	}
	peer := &PeerEntity{}
	for rows.Next() {
		var uname string
		var pubkey string
		var prikey string
		var ipaddr string
		err = rows.Scan(&uname, &pubkey, &prikey, &ipaddr)
		if err != nil {
			fmt.Println(err)
		}
		peer.Username = uname
		peer.PublicKey = pubkey
		peer.PrivateKey = prikey
		peer.IpAddress = ipaddr
	}
	return peer
}

func SelectAllDbPeer() []*PeerEntity {
	sqlStmt := "SELECT * FROM peer"
	rows, err := Db.Query(sqlStmt)
	if err != nil {
		fmt.Println("Error querying data, probably permission not granted")
		os.Exit(1)
	}
	peers := make([]*PeerEntity, 0, 10)
	for rows.Next() {
		var uname string
		var pubkey string
		var prikey string
		var ipaddr string
		err = rows.Scan(&uname, &pubkey, &prikey, &ipaddr)
		if err != nil {
			fmt.Println(err)
		}
		peers = append(peers, &PeerEntity{
			Username:   uname,
			PublicKey:  pubkey,
			PrivateKey: prikey,
			IpAddress:  ipaddr,
		})
	}
	return peers
}

func SelectAllDbPeerIpAddress() []string {
	sqlStmt := "SELECT ip_address FROM peer"
	rows, err := Db.Query(sqlStmt)
	if err != nil {
		fmt.Println("Error querying data, probably permission not granted")
		os.Exit(1)
	}
	addrs := make([]string, 0, 10)
	for rows.Next() {
		var ipaddr string
		err = rows.Scan(&ipaddr)
		if err != nil {
			fmt.Println(err)
		}
		addrs = append(addrs, ipaddr)
	}
	return addrs
}

func SelectDbPeerByIpAddress(addr string) *PeerEntity {
	sqlStmt := "SELECT * FROM peer where ip_address=? limit 1"
	rows, err := Db.Query(sqlStmt, addr)
	if err != nil {
		fmt.Println("Error querying data, probably permission not granted")
		os.Exit(1)
	}
	for rows.Next() {
		var uname string
		var pubkey string
		var prikey string
		var ipaddr string
		err = rows.Scan(&uname, &pubkey, &prikey, &ipaddr)
		if err != nil {
			fmt.Println(err)
			return &PeerEntity{}
		}
		return &PeerEntity{
			Username:   uname,
			PublicKey:  pubkey,
			PrivateKey: prikey,
			IpAddress:  ipaddr,
		}
	}
	return &PeerEntity{}
}

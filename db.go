/*
 *
 *  * Copyright (c) Mike Johnson 2024.
 *  *
 *  * Redistribution and use in source and binary forms, with or without
 *  * modification, are permitted provided that the following conditions
 *  * are met:
 *  *
 *  * 1. Redistributions of source code must retain the above copyright
 *  *    notice, this list of conditions and the following disclaimer.
 *  *
 *  * 2. Redistributions in binary form must reproduce the above copyright
 *  *    notice, this list of conditions and the following disclaimer in
 *  *    the documentation and/or other materials provided with the
 *  *    distribution.
 *  *
 *  * 3. Neither the name of the copyright holder nor the names of its
 *  *    contributors may be used to endorse or promote products derived
 *  *    from this software without specific prior written permission.
 *  *
 *  * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 *  * “AS IS” AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 *  * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 *  * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 *  * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 *  * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
 *  * TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
 *  * PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 *  * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 *  * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 *  * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package main

import (
	"database/sql"
	"encoding/csv"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func CreateAndPopulateSQLiteDatabaseAll() {
	db, err := sql.Open("sqlite3", "./body.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStatement := `
	DROP TABLE IF EXISTS body;
create table body
(
    name                      TEXT,
    short_name                TEXT,
    url_name                  TEXT,
    tags                      TEXT,
    home_page                 TEXT,
    publication_scheme        TEXT,
    disclosure_log            TEXT,
    notes                     TEXT,
    created_at                TEXT,
    updated_at                TEXT,
    version                   INT,
    defunct                   INT,
    categories                TEXT,
    top_level_categories      TEXT,
    single_top_level_category TEXT
);
	`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStatement)
		return
	}

	//txn, err := db.Begin()
	//if err != nil {
	//	log.Fatal(err)
	//}

	// #############################################################################################

	var csvFile, _ = os.Open("all-authorities.csv")
	//var columnNames[] string

	reader := csv.NewReader(csvFile)
	headers, _ := reader.Read()
	println("headers: ", headers)
	records, err := reader.ReadAll()
	println(records[0])
	if err != nil {
		return
	}

	//for _, row := range rows {
	//
	//	tagsList := strings.Split(row[3], " ")
	//	var name = row[1]
	//	var shortName = row[2]
	//	var urlName = row[3]
	//	var tags = row[4]
	//	var homePage = row[5]
	//	var publicationScheme = row[6]
	//	var disclosureLog = row[7]
	//	var notes = row[8]
	//	var createdAt = row[9]
	//	var updatedAt = row[10]
	//	var version = row[11]
	//}

	queryStringBuilder := strings.Builder{}
	//queryStringBuilder.WriteString("insert into body")
	queryStringBuilder.WriteString("short_name, name, url_name, tags, ")
	queryStringBuilder.WriteString("home_page, publication_scheme, disclosure_log,")
	queryStringBuilder.WriteString(" notes, created_at, updated_at, version, defunct, ")
	queryStringBuilder.WriteString("categories, top_level_categories, single_top_level_category")
	queryStringBuilder.WriteString(" values(")
	for _, record := range records[1:] {
		row := make(map[string]string)
		for i, value := range record {
			if i < len(headers) {
				row[headers[i]] = value
				queryStringBuilder.WriteString(value)
			}
		}

	}
	//statement, err := txn.Prepare(" values(""?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer statement.Close()
	//
	//for i := 0; i < 100; i++ {
	//	_, err = statement.Exec(i, fmt.Sprintf("%03d", i))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
	//
	//err = txn.Commit()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//rows, err := db.Query("select url_name, name from body")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var id int
	//	var name string
	//	err = rows.Scan(&id, &name)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(id, name)
	//}
	//err = rows.Err()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//statement, err = db.Prepare("select name from body where short_name = ?")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer statement.Close()
	//
	//var name string
	//err = statement.QueryRow("3").Scan(&name)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(name)
	//
	//_, err = db.Exec("delete from body")
	//if err != nil {
	//	log.Fatal(err)
	//}

}

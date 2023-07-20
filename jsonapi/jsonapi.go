package jsonapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"gomailinglist/mdb"
	"io"
	"log"
	"net/http"
)

func setJsonHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")


}

func fromJson[T any](body io.Reader, target T) {
	buff := new(bytes.Buffer)
	buff.ReadFrom(body)

	json.Unmarshal(buff.Bytes(), &target)

}

func returnJson[t any](w http.ResponseWriter, withData func() (t, error)){
	setJsonHeader(w)

	data, serverErr := withData()

	if serverErr != nil {
		w.WriteHeader(500)
		serverErrJson, err := json.Marshal(&serverErr)
		if err != nil {
			log.Println(err)
			return

		}

		w.Write(serverErrJson)
		return

	}

	dataJson, err := json.Marshal(&data)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		return
	}
	w.Write(dataJson)


}

func returnErr(w http.ResponseWriter, err error, code int) {
	returnJson(w, func() (interface{}, error){
		errorMessage := struct{
			Err string
		}{
			Err: err.Error(),
		}
		w.WriteHeader(code)
		return errorMessage, nil
	})
}

func CreateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
		if req.Method != "POST" {
			return 
		}

		entry := mdb.EmailEntry{}
		fromJson(req.Body, &entry)

		if err := mdb.CreateEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)

		}

		returnJson(w, func() (interface{}, error){
			log.Println("Json createEmail")
			return mdb.GetEmail(db, entry.Email)
		})
	})
}

func GetEmail(db *sql.DB) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
		if req.Method != "GET" {
			return 
		}

		entry := mdb.EmailEntry{}
		fromJson(req.Body, &entry)

		returnJson(w, func() (interface{}, error){
			log.Println("Json GetEmail")
			return mdb.GetEmail(db, entry.Email)
		})
	})
}
func GetBatchEmail(db *sql.DB) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
		if req.Method != "GET" {
			return 
		}

		queryOptions := mdb.GetEmailBatchQueryParams{}
		fromJson(req.Body, &queryOptions)

		if queryOptions.Count <= 0 || queryOptions.Page <= 0 {
			returnErr(w, errors.New("Page and count fields required"), 400)

		}

		returnJson(w, func() (interface{}, error){
			log.Println("JSON getEMail batch")
			return mdb.GetEmailBatch(db, queryOptions)
		} )
			
		
	})
}

func UpdateEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
		if req.Method != "PUT" {
			return 
		}

		entry := mdb.EmailEntry{}
		fromJson(req.Body, &entry)

		if err := mdb.UpdateEmail(db, entry); err != nil {
			returnErr(w, err, 400)

		}

		returnJson(w, func() (interface{}, error){
			log.Println("Json updateEmail")
			return mdb.GetEmail(db, entry.Email)
		})
	})
}
func DeleteEmail(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request){
		if req.Method != "POST" {
			return 
		}

		entry := mdb.EmailEntry{}
		fromJson(req.Body, &entry)

		if err := mdb.DeleteEmail(db, entry.Email); err != nil {
			returnErr(w, err, 400)

		}

		returnJson(w, func() (interface{}, error){
			log.Println("Json deleted email")
			return mdb.GetEmail(db, entry.Email)
		})
	})
}




func Serve(db *sql.DB, bind string) {
	http.Handle("/email/create", CreateEmail(db))
	http.Handle("/email/get", GetEmail(db))
	http.Handle("/email/get_batch", GetBatchEmail(db))
	http.Handle("/email/update", UpdateEmail(db))
	http.Handle("/email/delete", DeleteEmail(db))
	
	log.Printf("JSON API Server Listening on %v \n", bind)
	err := http.ListenAndServe(bind, nil)
	if err != nil {
		log.Fatal("JSon Server Error")
	}
}


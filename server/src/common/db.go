package common

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)
type DBContext struct{
	DB *sql.DB
	SelectStmt 	map[string]*sql.Stmt
	InsertStmt 	map[string]*sql.Stmt
	UpdateStmt	map[string]*sql.Stmt
}
func Dbconnect(dsn string, min int, max int) (*DBContext, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error in DB connection.", err)
		return nil, err
	} else {
		log.Println("DB connection established.")
	}
	db_cntx := &DBContext{
		DB: db,
	}
	db.SetMaxOpenConns(max)
	db.SetMaxIdleConns(min)

	return db_cntx, err
}
func (db *DBContext) ExecutePreparedSelectStatement(table string, filters map[string]string, cols []string) (result []map[string]string, db_error error){
	prepared_stmt_key := table+"_"
	if len(cols) > 0{
		prepared_stmt_key += strings.Join(cols, "_") + "_"
	}
	var args []interface{}
	for k, v := range filters{
		prepared_stmt_key += k+"_"
		args = append(args, v)
	}
	prepared_stmt_key = strings.Trim(prepared_stmt_key, "_")
	if stmt, ok := db.SelectStmt[prepared_stmt_key]; ok{
		res, err := stmt.Query(args...)
		if err != nil{
			log.Println("DB Error: error in executing statement. ", err)
			db_error = err
			return
		}
		defer res.Close()
		col_names, err := res.Columns()
		if err != nil {
			log.Println("DB Error: error in fetching column names.", err)
			return nil, err
		}
		values := make([]sql.RawBytes, len(col_names))
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		for res.Next() {
			err = res.Scan(scanArgs...)
			if err != nil {
				return nil, err
			}
			var value string
			new_row := make(map[string]string)
			for i, col := range values {
				if col == nil {
					value = "NULL"
				} else {
					value = string(col)
				}
				new_row[col_names[i]] = string(value)
			}
			result = append(result, new_row)
		}
	}else{
		log.Println("statement does not exist for the query: ", prepared_stmt_key)
		db_error = errors.New("statement does not exist for the query..." + prepared_stmt_key)
		return 
	}
	return
}

func (db *DBContext) DbSelect(table string, filters map[string]string, columns []string) (result []map[string]string, db_error error) {
	var query, all_col string
	temp , err := db.ExecutePreparedSelectStatement(table, filters, columns)
	if err == nil{
		log.Println("result is fetched from prepared statment.  ")
		return temp, nil
	}
	
	log.Println("fetching result using db connection...")
	if len(columns) == 0 {
		all_col = "*"
	} else {
		all_col = strings.Join(columns, ",")
	}
	var where_keys []string

	var args []interface{}
	for k, v := range filters {
		where_keys = append(where_keys, fmt.Sprintf(" %s = ? ", k))
		args = append(args, v)
		//where_values = append(where_values, v)
	}

	query = fmt.Sprintf("select %s from %s where %s ", all_col, table, strings.Join(where_keys, " and "))
	log.Println("query:", query, args)
	
	rows, err := db.DB.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	col_names, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([]sql.RawBytes, len(col_names))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// var result []map[string]string
	result = make([]map[string]string, 0)

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		var value string
		new_row := make(map[string]string)
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			new_row[col_names[i]] = string(value)
		}
		result = append(result, new_row)
	}
	return result, nil
}

func (db *DBContext) DbInsert(table string, values map[string]interface{}) (result sql.Result, db_error error) {

	var query = "insert into " + table + " "
	var columns []string
	var args_value []interface{}
	var args []string

	for col := range values {
		columns = append(columns, col)
		args = append(args, "?")
		args_value = append(args_value, values[col])
	}
	if len(columns) != len(args_value) {
		db_error = errors.New("db error: invalid column length")
		return
	}
	query = query + " (" + strings.Join(columns, ",") + ") values ( " + strings.Join(args, ",") + " )"
	log.Println("Query: ", query, " Args: ",args_value)
	//log.Println(query, args_value)
	stmt, err := db.DB.Prepare(query)
	if err != nil{
		log.Println("DB Error: error in preparing statement. ", err)
		db_error = err
		return 
	}
	defer stmt.Close()
	result, db_error = stmt.Exec(args_value...)
	if err != nil{
		log.Println("DB Error: error in executing statement. ", err)
	}
	return
}
func (db *DBContext) DbUpdate(table string, filters map[string]string, values map[string]interface{}) (result sql.Result, db_error error) {

	var query = "update " + table + " set "
	var args []interface{}
	var where_keys []string

	for col := range values {
		query = query + col+"=?, "
		args = append(args, values[col])
	}
	query = strings.Trim(query, ", ")
	log.Println("Query: ", query, " Args: ",args)
	
	for k, v := range filters {
		where_keys = append(where_keys, fmt.Sprintf(" %s = ? ", k))
		args = append(args, v)
	}
	if len(where_keys) <= 0{
		log.Println("DB Error: where clause is compulsory")
		db_error = errors.New("where clause is compulsory in select query")
		return
	}
	where_clause := strings.Join(where_keys, " and ")
	query = query + " where " + where_clause
	log.Println("Query: ", query, " Args: ",args)

	stmt, err := db.DB.Prepare(query)
	if err != nil{
		log.Println("DB Error: error in preparing statement. ", err)
		db_error = err
		return 
	}
	defer stmt.Close()
	result, db_error = stmt.Exec(args...)
	if err != nil{
		log.Println("DB Error: error in executing statement. ", err)
	}
	return
}

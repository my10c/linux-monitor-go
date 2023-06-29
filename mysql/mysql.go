// Copyright (c) 2017 - 2017 badassops
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//	* Redistributions of source code must retain the above copyright
//	notice, this list of conditions and the following disclaimer.
//	* Redistributions in binary form must reproduce the above copyright
//	notice, this list of conditions and the following disclaimer in the
//	documentation and/or other materials provided with the distribution.
//	* Neither the name of the <organization> nor the
//	names of its contributors may be used to endorse or promote products
//	derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSEcw
// ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
// DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Version		:	0.1
//
// Date			:	June 4, 2017
//
// History	:
// 	Date:			Author:		Info:
//	June 4, 2017	LIS			First Go release
//
// TODO:

package mysql

import (
	"fmt"
	"strconv"

	myGlobal	"github.com/my10c/linux-monitor-go/global"
	myUtils		"github.com/my10c/linux-monitor-go/utils"
	myThreshold	"github.com/my10c/linux-monitor-go/threshold"

	_ "github.com/go-sql-driver/mysql"
	"database/sql"
)

type dbMysql struct {
	*sql.DB
}

// Function to create the mysql object and connects to the mysql
func New(mysqlCfg map[string]string) *dbMysql {
	// set the username and password
	mysql_user := fmt.Sprintf("%s:%s@", mysqlCfg["username"], mysqlCfg["password"])
	// need to create the string tcp(fqdn:port) according to the docs, set host, port, database and option
	mysql_host_db := fmt.Sprintf("tcp(%s:%s)/%s?parseTime=true", mysqlCfg["hostname"], mysqlCfg["port"], mysqlCfg["database"])
	// set the full authentication string
	auth_string := mysql_user + mysql_host_db
	db, err := sql.Open("mysql", auth_string)
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// check we can make the connection
	err = db.Ping()
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// make sure the connection get close
	return &dbMysql{db}
}

// Function to check write == insert
func (db *dbMysql) CheckWrite(table string, field string, data string) error {
	// Prepare statement for inserting data
	stmt, err := db.Prepare(fmt.Sprintf("INSERT INTO %s (%s) VALUE ('%s')", table, field, data))
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// make sure to close the statement
	defer stmt.Close()
	_, err = stmt.Exec()
	return err
}

// Function to check read == select
func (db *dbMysql) CheckRead(table string, field string, data string) error {
	// Prepare statement for reading the data
	stmt, err := db.Prepare(fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s'", table, field, data))
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// make sure to close the statement
	defer stmt.Close()
	_, err = stmt.Exec()
	return err
}

// Function to check delete == delete
func (db *dbMysql) CheckDelete(table string, field string, data string) error {
	// Prepare statement for delete the data
	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE %s = '%s'", table, field, data))
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// make sure to close the statement
	defer stmt.Close()
	_, err = stmt.Exec()
	return err
}

// Function on of the check the table space, can we create?
func (db *dbMysql) CreateTable(table string) error {
	// Prepare statement to create a table
	stmt, err := db.Prepare(fmt.Sprintf("CREATE TABLE %s (timestamp varchar(128))", table))
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// make sure to close the statement
	defer stmt.Close()
	_, err = stmt.Exec()
	return err
}

// Function on of the check the table space, can we delete?
func (db *dbMysql) DropTable(table string) error {
	// Prepare statement to drop a table
	stmt, err := db.Prepare(fmt.Sprintf("DROP TABLE %s", table))
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// make sure to close the statement
	defer stmt.Close()
	_, err = stmt.Exec()
	return err
}

// Get the Slavestatus this will create a map with the current values
func (db *dbMysql) getSlaveStatus() (map[string]interface{}, error) {
	rows, err := db.Query("SHOW SLAVE STATUS")
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// if slave is set we should get exacly 1 row!
	if !rows.Next() {
		err := fmt.Errorf("Server is not set as slave!")
		return nil, err
	}
	defer rows.Close()
	// since layout could change, we get the columns info/names
	columns, err := rows.Columns()
	values := make([]interface{}, len(columns))
	for cnt := range values {
		var value sql.RawBytes
		values[cnt] = &value
	}
	if err = rows.Scan(values...); err != nil {
		return nil, err
	}
	// now get the values and set these into the map
	slaveInfo := make(map[string]interface{})
	for cnt, name := range columns {
		currValueStr := string(*values[cnt].(*sql.RawBytes))
		// if if Seconds_Behind_Master is empty then slave is not running!
		// we break and return an error since there is no point to go futher
		if name == "Seconds_Behind_Master" && len(currValueStr) == 0 {
			err := fmt.Errorf("Slave is not running!")
			return nil, err
		}
		// convert value string to int, if its a int (digits)
		currValueInt, err := strconv.ParseInt(currValueStr, 10, 64)
		if err == nil {
			slaveInfo[name] = currValueInt
		} else {
			// in case of string we need to handle mysql NULL's
			// we set these to the string "NULL"
			if len(currValueStr) == 0 {
				slaveInfo[name] = "NULL"
			} else {
				slaveInfo[name] = currValueStr
			}
		}
	}
	return slaveInfo, nil
}

// Check Slave_IO_Running and Slave_SQL_Running
func (db *dbMysql) SlaveStatusCheck() (int, error) {
	currStatus, err := db.getSlaveStatus()
	if err != nil {
		return myGlobal.CRITICAL, err
	}
	// Slave_IO_Running and Slave_SQL_Running must be Yes
	slaveIO := currStatus["Slave_IO_Running"].(string)
	slaveSQL := currStatus["Slave_SQL_Running"].(string)
	if slaveIO != "Yes" || slaveSQL != "Yes" {
		err := fmt.Errorf("Slave issues; Slave_IO_Running %s, Slave_SQL_Running %s", slaveIO, slaveSQL)
		return myGlobal.CRITICAL, err
	}
	err = fmt.Errorf("Status Slave_IO_Running %s, Slave_SQL_Running %s", slaveIO, slaveSQL)
	return myGlobal.OK, err
}

// Check the Seconds_Behind_Master value
func (db *dbMysql) SlaveLagCheck(warning uint64, critical uint64) (int, error) {
	currStatus, err := db.getSlaveStatus()
	if err != nil {
		return myGlobal.CRITICAL, err
	}
	currStatusInt := currStatus["Seconds_Behind_Master"].(uint64)
	currTholdStatus := myThreshold.CalculateUsage(false, false, warning, critical, currStatusInt, 0)
	if currTholdStatus != 0 {
		err := fmt.Errorf("Slave is behind by %d", currStatusInt)
		return currTholdStatus, err
	}
	err = fmt.Errorf("Slave lag %d", currStatusInt)
	return currTholdStatus, err
}

// Get current process count
func (db *dbMysql) ProcessStatusCheck(warning uint64, critical uint64) (int, error) {
	var totalRows uint64
	err := db.QueryRow("SELECT COUNT(*) FROM information_schema.PROCESSLIST").Scan(&totalRows)
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	defer db.Close()
	myUtils.ExitWithNagiosCode(myGlobal.UNKNOWN, err)
	// we need to subtract 1 as we should not count our self
	realRows := totalRows - 1
	currTholdStatus := myThreshold.CalculateUsage(false, false, warning, critical, realRows , 0)
	if currTholdStatus != 0 {
		err := fmt.Errorf("Process count issue, current processes count %d", realRows)
		return currTholdStatus, err
	}
	err = fmt.Errorf("Current proces count %d", realRows)
	return currTholdStatus, err
}

// basic check: read, write and delete
func (db *dbMysql) BasisCheck(table string, field string, data string) (int, error) {
	if err := db.CheckWrite(table, field, data); err != nil {
		db.Close()
		return myGlobal.CRITICAL, err
	}
	if err := db.CheckRead(table, field, data); err != nil {
		db.Close()
		return myGlobal.CRITICAL, err
	}
	if err := db.CheckDelete(table, field, data); err != nil {
		db.Close()
		return myGlobal.CRITICAL, err
	}
	err := fmt.Errorf("Passed, INSERT, SELECT and DELETE")
	return myGlobal.OK, err
}

// basic table space check: create and drop
func (db *dbMysql) DropCreateCheck(tablename string) (int, error) {
	if err := db.CreateTable(tablename); err != nil {
		db.Close()
		return myGlobal.CRITICAL, err
	}
	if err := db.DropTable(tablename); err != nil {
		db.Close()
		return myGlobal.CRITICAL, err
	}
	err := fmt.Errorf("Passed, create and drop table %s", tablename)
	return myGlobal.OK, err
}

// read check
func (db *dbMysql) ReadCheck(table string, field string) (int, error) {
	if err := db.CheckRead(table, field, "%"); err != nil {
		db.Close()
		return myGlobal.CRITICAL, err
	}
	err := fmt.Errorf("Passed, SELECT")
	return myGlobal.OK, err
}

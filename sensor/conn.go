
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azuresql "github.com/microsoft/go-mssqldb/azuread"
)

// # Connection Function (conn.go)
// Input  : server and database strings from azure
// Output : *sql.DB, a live database connection pointer

/*
Creates a database connection for my specific instance of azure. It first
creates the credential from the entraID on the machines (THIS NEEDS TO BE SET
UP ON THE MACHINE TO WORK, or else be running on azure itself), before crafting
a connection string from the currently hardcoded values. Once the connection 
has been made, the function then pings it, and returns the pointer if there
have been no errors. 
*/

var db *sql.DB

// ## Private/Helper Functions ------------------------------------------------

func getAzureID() (*azidentity.DefaultAzureCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, errors.New("Error obtaining azure credential: " + err.Error())
	}

	return cred, nil
}

func getConnectionString(server string, database string) string {
	connString := fmt.Sprintf(
		"sqlserver://%s?database=%s&Encrypt=true&fedauth=ActiveDirectoryDefault",
    server, database,
	)
	return connString
}

func testDBConnection(db *sql.DB) (error) {
	if err := db.PingContext(context.Background()); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}
	return nil
}

// ## Public/Principle Function ===============================================

func getDBConnection(server string, database string) (*sql.DB, error) {
	_, err := getAzureID()
	if err != nil { return nil, err }

	db, err := sql.Open(
		azuresql.DriverName,
		getConnectionString(server, database),
	)
	if err != nil { return nil, fmt.Errorf("error connecting to DB: %w", err) }
	if err := testDBConnection(db); err != nil { return nil, err }

	return db, nil
}
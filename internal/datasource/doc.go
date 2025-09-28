// Package datasource provides interfaces and implementations
// for interacting with various types of data sources (SQL, NoSQL, etc.).
//
// This package defines a generic IDataSource interface along with concrete
// implementations such as PostgresDataSource and MongoDataSource.
//
// Example usage:
//
//	ds, err := factories.CreateDataSource("postgres")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer ds.Close()
//
//	if err := ds.Connect(); err != nil {
//	    log.Fatal(err)
//	}
//	result := ds.Query("SELECT * FROM users")
//	fmt.Println(result)
package datasource

handlersocket-go
================

Go library for connecting to HandlerSocket Mysql plugin.  See github.com/ahiguti/HandlerSocket-Plugin-for-MySQL/


## Installation

```bash
$ go get github.com/azhai/handlersocket-go

$ #Install Mariadb and HandlerSocket in CentOS7
$ su -
# yum install mariadb-server
# tee /etc/my.cnf.d/innodb.cnf <<EOD
[mysqld]
default-time-zone       = "+08:00"
innodb_file_per_table   = 1
innodb_buffer_pool_size = 402653184 #384MB

loose_handlersocket_address         = "127.0.0.1"
loose_handlersocket_plain_secret    = ""
loose_handlersocket_plain_secret_wr = ""
loose_handlersocket_port            = 9998
loose_handlersocket_port_wr         = 9999
EOD
# mysql -u root -p
MariaDB [(none)]> INSTALL PLUGIN handlersocket SONAME 'handlersocket.so';
MariaDB [(none)]> use `test`;
MariaDB [(none)]> CREATE TABLE `people` ( 
  `id` INTEGER UNSIGNED NOT NULL AUTO_INCREMENT, 
  `name` VARCHAR(60) NOT NULL DEFAULT '',
  `age` INTEGER UNSIGNED NOT NULL DEFAULT 0,
  `dob` DATE DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `name_key` (`name`) );
MariaDB [(none)]> INSERT INTO `people` (`name`, `dob`) VALUES
    ('Joe', '1985-02-23'), ('Mary', '1982-03-12'), ('Gordon', '1978-09-02');
MariaDB [(none)]> UPDATE `people` SET `age`=ROUND(DATEDIFF(NOW(),`dob`)/365) WHERE `age`=0;
```


## HandlerSocketWrapper

```go
package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/azhai/handlersocket-go"
)


func CalcAge(birthday string) int {
	t, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return -1
	}
	h := time.Since(t).Hours()
	return int(math.Round(h/365/24))
}


func main() {
	var index *handlersocket.HandlerSocketIndex
	hs := handlersocket.NewWrapper("127.0.0.1", 9998, 9999)
	defer hs.Close()
	
	//Send: P\t1\ttest\tpeople\tPRIMARY\tid,name,age,dob\n
	columns4 := []string{"id", "name", "age", "dob"}
	index = hs.WrapIndex("test", "people", "", columns4...)
	
	//Send: 1\t<=\t1\t3\t2\t0\n
	rows, _ := index.FindAll(2, 0, "<=", "3")
	for i := range rows {
		fmt.Println(rows[i].Data)
	}
	
	columns3 := []string{"name", "age", "dob"}
	//Send: P\t2\ttest\tpeople\tname_key\tname,age,dob\n
	index = hs.WrapIndex("test", "people", "name_key", columns3...)
	//Send: 2\t=\t1\tGordon\t1\t0\n
	row, _ := index.FindOne("=", "Gordon")
	fmt.Printf("%s %s\n", row.Data["name"], row.Data["dob"])
	
	//Send: 2\t+\t3\tFred\t0\t1965-07-07\n
	//index.Insert("Fred","0","1965-07-07")
	
	//Send: 2\t=\t1\tFred\t1\t0\tU\Fred,53,1965-07-07\n
	age := strconv.Itoa(CalcAge("1965-07-07"))
	index.Update(1, "=", []string{"Fred"}, "Fred", age, "1965-07-07")
	
	//Send: 2\t=\t1\tFred\t1\t0\tD\n
	index.Delete(1, "=", []string{"Fred"})
}
```


## (Old) Read Example  - Best examples are in the TEST file.

	hs := New()

	// Connect to database
	hs.Connect("127.0.0.1", 9998, 9999)
	defer hs.Close()
	hs.OpenIndex(1, "gotesting", "kvs", "PRIMARY", "id", "content")

	found, _ := hs.Find(1, "=", 1, 0, "brian")

	for i := range found {
			fmt.Println(found[i].Data) 
		}

	fmt.Println(len(found), "rows returned")


## (Old) Write Example

	hs := New()
	hs.Connect("127.0.0.1", 9998, 9999) // host, read port, write port
	defer hs.Close()

	// id is varchar(255), content is text
	hs.OpenIndex(3, "gotesting", "kvs", "PRIMARY", "id", "content")

	err := hs.Insert(3,"mykey1","a quick brown fox jumped over a lazy dog")


## (Old) Modify Example

	var keys, newvals []string
	keys = make([]string,1)
	newvals = make([]string,2)
	keys[0] = "blue3"
	newvals[0] = "blue7"
	newvals[1] = "some new thing"
	count, err := hs.Modify(3, "=", 1, 0, "U", keys, newvals)
	if err != nil {
		t.Error(err)
		}
	fmt.Println("modified", count, "records")


## Copyright and licensing

Licensed under **Apache License, version 2.0**.  
See file LICENSE.


## Contact

Brian Ketelsen - bketelsen@gmail.com

## Known bugs

No known bugs, but testing is far from comprehensive.

Working:  OpenIndex, Find, Insert,  Update/Delete


## Todo

Provide a layer of abstraction from the wire-level implementation of HandlerSocket to make a more intuitive interface.




## Credits and acknowledgments


Took some inspiration from the original GoMySQL implementation, although I've backed much of that out in this initial release.
https://github.com/Philio/GoMySQL
I can see how it would be extremely useful for GoMySQL or GoDBI to use HandlerSocket in the background for simple finds, inserts, etc.


## ChangeLog
13/8/2018 (v0.0.5)
	Add HandlerSocketWrapper by azhai(https://github.com/azhai/handlersocket-go)

1/10/2015 (v0.0.4)
	Fix bugs by nisboo(https://github.com/nisboo/handlersocket-go), for Go v1.8+

1/20/2011
	Updated library extensively
	now working OpenIndex and Find commands
	
1/21/2011
	Insert works now
	
3/14/2011
	Modify and Delete work now - need more tests!



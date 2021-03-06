package handlersocket

import "errors"

const DefaultIndexName = "PRIMARY"

type HandlerSocketIndex struct {
	Socket    *HandlerSocket
	indexNo   int      //1-base
	indexName string
	dbName    string
	table     string
	columns   []string
}

type HandlerSocketWrapper struct {
	Socket    *HandlerSocket
	lastNo    int
}

func NewWrapper(host string, rport, wport int) *HandlerSocketWrapper {
	auth := &HandlerSocketAuth{}
	auth.host = host
	auth.readPort = DefaultReadPort
	auth.writePort = DefaultWritePort
	if rport > 0 {
		auth.readPort = rport
	}
	if wport > 0 {
		auth.writePort = wport
	}
	obj := &HandlerSocketWrapper{lastNo: 0}
	obj.Socket = New()
	obj.Socket.auth = auth
	return obj
}

func (this *HandlerSocketWrapper) Close() error {
	if this.Socket.connected {
		return this.Socket.Close()
	}
	return nil
}

func (this *HandlerSocketWrapper) WrapIndex(dbName, table, indexName string, columns ...string) *HandlerSocketIndex {
	this.lastNo++
	if indexName == "" {
		indexName = DefaultIndexName
	}
	index := &HandlerSocketIndex{
		dbName: dbName, table: table, columns: columns, indexName: indexName,
	}
	index.Socket = this.Socket
	index.indexNo = this.lastNo
	index.Socket.OpenIndex(index.indexNo, dbName, table, indexName, columns...)
	return index
}

func (this *HandlerSocketIndex) FindAll(limit int, offset int, oper string, where ...string) ([]HandlerSocketRow, error) {
	rows, err := this.Socket.Find(this.indexNo, oper, limit, offset, where...)
	if err != nil {
		panic(err)
	}
	return rows, err
}

func (this *HandlerSocketIndex) FindOne(oper string, where ...string) (HandlerSocketRow, error) {
	rows, err := this.FindAll(1, 0, oper, where...)
	if rows == nil || len(rows) == 0 {
		err = errors.New("Nothing found")
		return HandlerSocketRow{}, err
	}
	return rows[0], err
}

func (this *HandlerSocketIndex) Insert(vals ...string) error {
	return this.Socket.Insert(this.indexNo, vals...)
}

func (this *HandlerSocketIndex) Delete(limit int, oper string, where []string) (int, error) {
	return this.Socket.Modify(this.indexNo, oper, limit, 0, "D", where, nil)
}

func (this *HandlerSocketIndex) Update(limit int, oper string, where []string, vals ...string) (int, error) {
	return this.Socket.Modify(this.indexNo, oper, limit, 0, "U", where, vals)
}

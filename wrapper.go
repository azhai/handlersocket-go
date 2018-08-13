package handlersocket

const DefaultIndexName = "PRIMARY"

type HandlerSocketIndex struct {
	hs        *HandlerSocket
	dbname    string
	table     string
	columns   []string
	indexName string
	indexNo   int //1-base
}

type HandlerSocketWrapper struct {
	hs     *HandlerSocket
	lastNo int
}

func NewWrapper(host string, rport, wport int) *HandlerSocketWrapper {
	auth := &HandlerSocketAuth{}
	auth.host = host
	if rport > 0 {
		auth.readPort = rport
	} else {
		auth.readPort = DefaultReadPort
	}
	if wport > 0 {
		auth.writePort = wport
	} else {
		auth.writePort = DefaultWritePort
	}
	obj := &HandlerSocketWrapper{lastNo: 0}
	obj.hs = New()
	obj.hs.auth = auth
	return obj
}

func (this *HandlerSocketWrapper) Close() error {
	if this.hs.connected {
		return this.hs.Close()
	}
	return nil
}

func (this *HandlerSocketWrapper) Unwrap() *HandlerSocket {
	return this.hs
}

func (this *HandlerSocketWrapper) WrapIndex(dbname, table, name string, columns ...string) *HandlerSocketIndex {
	this.lastNo++
	index := &HandlerSocketIndex{
		dbname: dbname, table: table, columns: columns, indexName: name,
	}
	index.hs = this.hs
	index.indexNo = this.lastNo
	return index
}

func (this *HandlerSocketIndex) Unwrap() *HandlerSocket {
	return this.hs
}

func (this *HandlerSocketIndex) FindAll(limit int, offset int, oper string, vals ...string) ([]HandlerSocketRow, error) {
	rows, err := this.hs.Find(this.indexNo, oper, limit, offset, vals...)
	if err != nil {
		panic(err)
	}
	return rows, err
}

func (this *HandlerSocketIndex) FindOne(oper string, vals ...string) (HandlerSocketRow, error) {
	rows, err := this.FindAll(1, 0, oper, vals...)
	return rows[0], err
}

func (this *HandlerSocketIndex) Insert(vals ...string) error {
	return this.hs.Insert(this.indexNo, vals...)
}

func (this *HandlerSocketIndex) Delete(limit int, oper string, keys []string) (int, error) {
	return this.hs.Modify(this.indexNo, oper, limit, 0, "D", keys, nil)
}

func (this *HandlerSocketIndex) Update(limit int, oper string, keys []string, vals ...string) (int, error) {
	return this.hs.Modify(this.indexNo, oper, limit, 0, "U", keys, vals)
}

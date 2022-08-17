package scanner

type IScanner interface {
	Scan(dest ...interface{}) error
}

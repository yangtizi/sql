package scanner

// TResult 结果
type TResult struct {
	LastInsertId int64
	RowsAffected int64
}

// NewResult 新结果
func NewResult(LastInsertId int64, RowsAffected int64) *TResult {
	r := &TResult{
		LastInsertId: LastInsertId,
		RowsAffected: RowsAffected,
	}

	return r
}

// // RowsAffected 数量
// func (m *TResult) RowsAffected() (int64, error) {
// 	return m.r.RowsAffected()
// }

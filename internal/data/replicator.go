package data

type Replicator interface {
	Replicate(tableName string) error
}

type RowReplicator struct {
	reader Reader
	writer Writer
}

func NewRowReplicator(reader Reader, writer Writer) *RowReplicator {
	return &RowReplicator{reader: reader, writer: writer}
}

func (r *RowReplicator) Replicate(tableName string) error {
	page := 0
	for {
		dml, err := r.reader.Read(tableName, page)
		if err != nil {
			return err
		}
		if dml == "" {
			break
		}
		if err = r.writer.Write(dml, tableName); err != nil {
			return err
		}
		page++
	}
	return nil
}

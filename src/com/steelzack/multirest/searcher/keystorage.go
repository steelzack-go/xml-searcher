package searcher

import (
	"github.com/gocql/gocql"
	"log"
	"bytes"
	"container/list"
)


type KeyStorage struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session

}

type KeyValues struct {
	value1uuid string
	value2uuid string
}

func (keystorage *KeyStorage) OpenDatabase(host string, port int) {
	cluster := gocql.NewCluster(host)
	cluster.ProtoVersion = 3
	cluster.Keyspace = "xml-searcher"
	cluster.CQLVersion = "3.3.1"
	cluster.Consistency = gocql.One
	cluster.Port = port

	session, err := cluster.CreateSession()
	if err != nil {
		log.Println(err)
		return
	}
	keystorage.cluster = cluster
	keystorage.session = session
}

func (keystorage *KeyStorage) Init() {
	if err := (*keystorage).session.Query("CREATE TABLE IF NOT EXISTS value_mapping ( value1UUID UUID, value2UUID UUID,  PRIMARY KEY (value1UUID, value2UUID));").Exec(); err != nil {
		log.Fatal(err)
	}
}

func (keystorage *KeyStorage) InsertKeys(value1UUID gocql.UUID, value2UUID gocql.UUID) {
	buffer := bytes.Buffer{}
	buffer.WriteString("INSERT INTO value_mapping (value1UUID, value2UUID) values (")
	buffer.WriteString(value1UUID.String())
	buffer.WriteString(",")
	buffer.WriteString(value2UUID.String())
	buffer.WriteString(");")
	queryinsert := buffer.String()
	if err := (*keystorage).session.Query(queryinsert).Exec(); err != nil {
		log.Fatal(err)
	}
}

func (keystorage *KeyStorage) CloseDatabase() {
	(*keystorage).session.Close()
}

// Simple value retrieval example
func (keystorage *KeyStorage) GetValue1FromValue2(value2Id string) string {
	buffer := bytes.Buffer{}
	buffer.WriteString("SELECT value1UUID FROM value_mapping WHERE value2UUID=")
	buffer.WriteString(value2Id)
	buffer.WriteString(" ALLOW FILTERING;")
	var value1uuid string
	if err := (*keystorage).session.Query(buffer.String()).Scan(&value1uuid); err != nil {
		log.Println(err)
	}
	return value1uuid
}

// List value retrieval example
func (keystorage *KeyStorage) GetValue2sFromValue1(value1Id string) list.List {
	buffer := bytes.Buffer{}
	buffer.WriteString("SELECT value2UUID FROM value_mapping WHERE value1UUID=")
	buffer.WriteString(value1Id)
	buffer.WriteString(" ALLOW FILTERING;")
	var value2uuid string
	var value2list list.List
	iter := (*keystorage).session.Query(buffer.String()).Iter()
	for (iter.Scan(&value2uuid)) {
		value2list.PushBack(value2uuid)
	}
	if err := iter.Close(); err !=nil {
		log.Println(err)
	}
	return value2list;
}
package config

type MysqlDB struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *MysqlDB) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

func (m *MysqlDB) GetLogMode() string {
	return m.LogMode
}

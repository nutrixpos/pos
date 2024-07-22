package globals

var DBHost string

func Init(SWP_DB_HOST string) {
	DBHost = SWP_DB_HOST
}

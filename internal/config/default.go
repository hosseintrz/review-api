package config

var defaultConf = []byte(`
DATABASE:
  DRIVER: postgres
  SOURCE: postgresql://root:secret@localhost:5432/feedback_db?sslmode=disable
SERVER_ADDRESS: 0.0.0.0:80
`)

package config

const PORT_ENVAR string = "PORT"
const DEFAULT_PORT string = "3000"

const PG_DATABASE_CONNECTION_URL_ENVAR string = "POSTGRESQL_URL"

const MAX_IMAGE_FILE_SIZE int = 16 * 1024 * 1024
const MAX_IMAGE_FILE_SIZE_MB int = MAX_IMAGE_FILE_SIZE / 1024 / 1024
const FILE_DOWNLOAD_USERAGENT string = "Mozilla/5.0 (Windows NT 10.0; rv:108.0) Gecko/20100101 Firefox/108.0"

const ZMQ_IMAGE_EMBEDDING_DAEMON_PORT_ENVAR string = "ZMQ_IMAGE_PORT"
const ZMQ_IMAGE_EMBEDDING_DAEMON_DEFAULT_PORT string = "5554"
const ZMQ_TEXT_EMBEDDING_DAEMON_PORT_ENVAR string = "ZMQ_TEXT_PORT"
const ZMQ_TEXT_EMBEDDING_DAEMON_DEFAULT_PORT string = "5553"

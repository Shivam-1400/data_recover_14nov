package model

type RedisConn struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type RedisQueue struct {
	QueueLvl1 string `json:"queue_lvl1"`
	QueueLvl2 string `json:"queue_lvl2"`
	QueueLvl3 string `json:"queue_lvl3"`
}

type Application struct {
	Worker   int    `json:"worker"`
	ReadPath string `json:"read_path"`
}

type Config struct {
	RedisConn   RedisConn   `json:"redis_conn"`
	RedisQueue  RedisQueue  `json:"redis_queue"`
	Application Application `json:"application"`
}

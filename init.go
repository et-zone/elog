package elog

var log *Log

func Init(cfg *Config)error{
	if cfg==nil{
		return errCfgNil
	}
	log = NewLog(cfg)
	if log ==nil{
		return errLogNil
	}
	return nil
}
func Close(){
	if log !=nil{
		log.Sync()
	}
}


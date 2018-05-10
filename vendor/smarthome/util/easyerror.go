package util

import (
	"log"
)

func PanicIf(status interface{}, msg ...interface{}) {
	switch v := status.(type) {
	case error:
		log.Panicln(v, msg)
	case bool:
		if v {
			log.Panicf("... unknown panic: %v\n\t ... \n", msg)
		}
	default:
		if v != nil {
			log.Panicf("... unknown panic: %v\n\t ... since %T:%v\n", msg, v, v)
		}
	}
}

// func WarnIf(status interface{}, msg ...interface{}) {
// 	switch v := status.(type) {
// 	case error:
// 		log.Warningln(v, msg)
// 	case bool:
// 		if v {
// 			log.Warningf("... unknown warning: %v\n\t ... \n", msg)
// 		}
// 	default:
// 		if v != nil {
// 			log.Warningf("... unknown warning: %v\n\t ... since %T:%v\n", msg, v, v)
// 		}
// 	}
// }

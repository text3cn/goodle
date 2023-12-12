package cryptokit

//import (
//	"statistics/kit/fn"
//
//	"github.com/dgrijalva/jwt-go"
//)
//
//func JWTparseToken(token string) map[string]interface{} {
//	token = fn.Str_replace("Bearer ", "", token, 1)
//	parseAuth, _ := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
//		return secret(), nil
//	})
//	// 将 token 中的内容存入 _map
//	claim := parseAuth.Claims.(jwt.MapClaims)
//	var _map map[string]interface{}
//	_map = make(map[string]interface{})
//	for key, val := range claim {
//		_map[key] = val
//	}
//	return _map
//}
//
//func secret() jwt.Keyfunc {
//	return func(token *jwt.Token) (interface{}, error) {
//		// 这是测试环境的密钥，目前的 token 没有加密，不用密钥就能解出来
//		return []byte("lrgcTAt^w4neJzrJfvwecY&DaE"), nil
//	}
//}

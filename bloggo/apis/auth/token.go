package auth

import (
	"encoding/json"
	"fmt"

	//"go/token"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(user_id uint32) (string, error) { //CREATE TOKEN WITH CLAIMS!

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()       //TOKEN EXPIRES IN 1HR
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //CREATING NEW TOKEN
	return token.SignedString([]byte(os.Getenv("API_SECRET"))) //RETURN SIGNEDSTRING(COMBO OF HEADER + PAYLOAD)
}

func TokenValid(r *http.Request) error { //TO CHECK WHETHER TOKEN CREATED IS VALID OR NOT
	tokenString := ExtractToken(r)                                                     //CALL TO EXTRACT() TO EXTRACT EXACT JWT FROM REQUEST MADE TO SERVER
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { //WE NEED TO PARSE THE TOKEN,IT RETURNS INTERFACE AND ERR
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method:%v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil //IF EVERYTHING OK,WE RETURN THE ORIGINAL KEY
	})

	if err != nil {
		return err
	}

	//EXTRACTING CLAIMS FROM TOKEN AND CHECKING ITS VALIDATION

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		Pretty(claims)
	}

	return nil
}

func ExtractToken(r *http.Request) string { //EXTRACTING TOKENS
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}

	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 { //SPLIT STRING AND RETURN DESIRED PART OF TOKEN AMONG 3 PARTS
		return strings.Split(bearerToken, " ")[1]
	}

	return "" //ELSE RETURNS EMPTY STRING

}

func ExtractTokenID(r *http.Request) (uint32, error) { //FUNC TO RETURN ONLY TOKENID
	tokenString := ExtractToken(r)                                                     //FIRST STORE IN TOKENSTRING THE ORIGINAL TOKEN
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) { //PARSE IT AS DONE BEFORE
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method:%v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_SECRET")), nil //IF SUCCESFULL PARSE FUNC RETURNS API SECRET KEY
	})

	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims) //AGAIN EXTRACTING CLAIMS,IF VALID CLAIMS THEN PROCEED!
	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["user_id"]), 10, 32)

		if err != nil {
			return 0, err
		}
		return uint32(uid), nil //IF ALL OK,RETURNS TOKEN ID
	}
	return 0, nil //ELSE RETURNS 0 OR NIL
}

func Pretty(data interface{}) {

	/*MarshalIndent is like Marshal(WHICH RETURNS JSON ENCODING OF ANY VALUE) but applies Indent to format the output.
	Each JSON element in the output will begin on a new
	line beginning with prefix followed by one or more copies of indent
	according to the indentation nesting.*/

	b, err := json.MarshalIndent(data, "", " ") //MARSHAL AND INDENT THE CLAIM DATA PASSED
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(string(b)) //IF ALL OK,RETURN MARSHAL DATA
}

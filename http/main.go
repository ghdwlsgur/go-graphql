package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var data map[string]user

/*
GraphQLObjectTypeConfig를 사용하여 "id" 및 "name" 필드가 있는 사용자 개체 유형을 만듭니다.
- 이름: 객체 유형의 이름
- 필드: GraphQLFields를 사용한 필드 맵
GraphQLFieldConfig 사용 필드 유형 설정
*/
var userType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

/*
GraphQLObjectTypeConfig를 사용하여 "user" 필드가 있는 쿼리가 개체 유형 생성: [userType] 유형:
- 이름: 객체 유형의 이름
- 필드: GraphQLFields를 사용한 필드 맵
필드의 설정 유형은 GraphQLFieldConifg를 사용하여 다음을 정의합니다.
- 유형: 필드의 유형
- Args: 현재 필드로 쿼리할 인수
- Resolve: [Args]의 매개변수를 사용하여 데이터를 쿼리하고 현재 유형의 값을 반환하는 함수
*/
var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					idQuery, isOK := p.Args["id"].(string)
					if isOK {
						return data[idQuery], nil
					}
					return nil, nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: queryType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	_ = importJSONDataFromFile("data.json", &data)

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={user(id:\"1\"){name}}'")
	http.ListenAndServe(":8080", nil)
}

func importJSONDataFromFile(fileName string, result interface{}) (isOK bool) {
	isOK = true
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Print("Error:", err)
		isOK = false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		isOK = false
		fmt.Print("Error:", err)
	}
	return
}

/*
1. go run main.go
2. execute command that curl -g http://localhost:8080/graphql?query={user(id:"1"){name}}

! result
{
	"data":{
		"user":{
				"name":"Dan"
		}
	}
}
*/

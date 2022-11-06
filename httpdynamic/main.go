package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/graphql-go/graphql"
)

// ! Schema-First Development
// 스키마 우선주의는 디자인 방법론의 일종이다. 개발시 스키마를 우선 개발하는 것이다. 여기서 스키마(Schema)란 데이터 타입의 집합이다. 이를 미리 정의해두면,
// 스키마 정의는 API 문서 같은 역할을 하며, 프론트엔드 개발자와 백엔드 개발자가 많은 의사소통에 대한 비용을 줄이고 빠른 개발을 할 수 있다는 장점이 있다.
// 백엔드 개발자는 어떤 데이터를 전달해야 하는지, 프론트엔드 개발자는 인터페이스 작업을 할 때 필요한 데이터를 정의할 수 있는 것이다.
var schema graphql.Schema

const jsonDataFile = "data.json"

func handleSIGUSR1(c chan os.Signal) {
	for {
		<-c
		fmt.Printf("Caught SIGUSR1. Reloading %s\n", jsonDataFile)
		err := importJSONDataFromFile(jsonDataFile)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
	}
}

func filterUser(data []map[string]interface{}, args map[string]interface{}) map[string]interface{} {
	for _, user := range data {
		for k, v := range args {
			if user[k] != v {
				// 이동할 레이블 지정
				goto nextuser
			}
			return user
		}
	nextuser:
		// 레이블 [실행할 코드 작성]
	}
	return nil
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v\n", result.Errors)
	}
	return result
}

func importJSONDataFromFile(fileName string) error {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var data []map[string]interface{}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return err
	}

	fields := make(graphql.Fields)
	args := make(graphql.FieldConfigArgument)

	for _, item := range data {
		for k := range item {
			fields[k] = &graphql.Field{
				Type: graphql.String,
			}
			args[k] = &graphql.ArgumentConfig{
				Type: graphql.String,
			}
		}
	}

	var userType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "User",
			Fields: fields,
		},
	)

	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"user": &graphql.Field{
					Type: userType,
					Args: args,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						return filterUser(data, p.Args), nil
					},
				},
			},
		})

	schema, _ = graphql.NewSchema(
		graphql.SchemaConfig{
			Query: queryType,
		},
	)

	return nil
}

func main() {

	c := make(chan os.Signal, 1)
	// SIGUSR1: 유저 정의 시그널1으로 log rotate시의 파일 제어에 사용
	signal.Notify(c, syscall.SIGUSR1)
	go handleSIGUSR1(c)

	err := importJSONDataFromFile(jsonDataFile)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/graphql?query={user(name:\"Dan\"){id,surname}}'")
	fmt.Printf("Reload json file: file -SIGUSR1 %s\n", strconv.Itoa(os.Getpid()))
	http.ListenAndServe(":8080", nil)

}

/* Result

1. go run main.go
2. bash ===============================================
Now server is running on port 8080
Test with Get: curl -g 'http://localhost:8080/graphql?query={user(name:"Dan"){id,surname}}'
Reload json file: file -SIGUSR1 13328
=======================================================
3. execute command that kill -SIGUSR1 13328
4. bash ===============================================
Caught SIGUSR1. Reloading data.json
=======================================================
5. execute command that curl -g ...
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

func main() {

	/* GraphQL이란 ?
	GraphQL은 페이스북에서 만든 쿼리 언어입니다. 국내에서 GraphQL API를 Open API로 공개한 곳은 드뭅니다. 또한, 해외의 경우,
	Github 사례(Github v4 GraphQL)를 찾을 수는 있지만, 전반적으로 GraphQL API를 Open API로 공개한 곳은 많지 않습니다.

	SQL은 데이터베이스 시스템에 저장된 데이터를 효율적으로 가져오는 것이 목적이고, GQL은 웹 클라이언트가 데이터를 서버로부터 효율적으로
	가져오는 것이 목적입니다. GQL의 문장은 주로 클라이언트 시스템에서 작성하고 호출합니다.

	GraphQL은 애플리케이션 프로그래밍 인터페이스(API)를 위한 쿼리 언어이자 서버측 런타임으로 클라이언트에게 요청한 만큼의 데이터를
	제공하는데 우선 순위를 둡니다. GraphQL은 API를 더욱 빠르고 유연하며 개발자 친화적으로 만들기 위해 설계되었습니다.

	서버사이드 gql 어플리케이션은 gql로 작성된 쿼리를 입력으로 받아 쿼리를 처리한 결과를 다시 클라이언트로 돌려줍니다.
	HTTP API 자체가 특정 데이터베이스나 플랫폼에 종속적이지 않은것 처럼 마찬가지로 gql 역시 어떠한 특정 데이터베이스나 플랫폼에 종속적이지 않습니다.
	심지어 네트워크 방식에도 종속적이지 않습니다. 일반적으로 gql의 인터페이스간 송수신은 네트워크 레이어 L7의 HTTP POST 메서드와 웹 소켓 프로토콜을
	활용합니다. 필요에 따라서는 얼마든지 L4의 TCP/UDP를 활용하거나 심지어 L2 형식의 이더넷 프레임을 활용할 수도 있습니다.

	REST API는 URL, METHOD 등을 조합하기 때문에 다양한 Endpoint가 존재합니다. 반면, gql은 단 하나의 Endpoint가 존재합니다.
	또한, gql API에서는 불러오는 데이터의 종류를 쿼리 조합을 통해서 결정합니다. 예를 들면, REST API에서는 각 Endpoint마다 데이터베이스 SQL 쿼리가 달라지는
	반면, gql API는 gql 스키마의 타입마다 데이터베이스 SQL 쿼리가 달라집니다.

	gql API를 사용하면 여러번 네트워크 호출을 할 필요 없이, 한번의 네트워크 호출로 처리할 수 있습니다.
	GraphiQL이라고 하는 통합 개발 환경(IDE) 내에 배포될 수도 있습니다. REST를 대체할 수 있는 GraphQL은 개발자가 단일 API 호출로
	다양한 데이터 소스에서 데이터를 끌어오는 요청을 구성할 수 있도록 지원합니다. 또한 GraphQL은 API 유지 관리자에게 기존 쿼리에 영향을 미치지 않고
	필드를 추가하거나 폐기할 수 있는 유연성을 부여합니다. 개발자는 자신이 선호하는 방식으로 API를 빌드할 수 있으며, GraphQL 사양은 이러한 API가 예측
	가능한 방식으로 작동하도록 보장합니다.
	*/

	/* GraphQL의 구조

	쿼리 / 뮤테이션 (query / mutation)
	쿼리와 뮤테이션 그리고 응답 내용의 구조는 상당히 직관적입니다. 요청하는 쿼리문의 구조와 응답 내용의 구조는 거의 일치합니다.

	! GraphQL 쿼리문
	{
		hero {
			name
		}
	}

	! 응답데이터 형식
	{
		"data": {
			"hero": {
				"name": "R2-D2"
			}
		}
	}

	gql에서는 굳이 쿼리와 뮤테이션을 나누는데 내부적으로 들어가면 사실상 이 둘은 별 차이가 없습니다. 쿼리는 데이터를 읽는데 사용하고,
	뮤테이션은 데이터를 변조(CUD)하는데 사용한다는 개념적인 규약을 정해 놓은 것 뿐입니다.

	오퍼레이션 네임 쿼리는 매우 편리합니다. 굳이 비유하자면 쿼리용 함수입니다. 데이터베이스에서의 프로시저(procedure) 개념과 유사하다고 생각하면
	됩니다. 이 개념 덕분에 여러분은 REST API를 호출할 때와 다르게, 한번의 인터넷 네트워크 왕복으로 여러분이 원하는 모든 데이터를 가져올 수 있습니다.
	데이터베이스의 프로시져는 DBA 혹은 백엔드 프로그래머가 작성하고 관리하였지만 gql 오퍼레이션 네임 쿼리는 클라이언트 프로그래머가 작성하고 관리합니다.

	gql이 제공하는 추가 기능 덕분에 백엔드 프로그래머와 프론트엔드 프로그래머의 협업 방식에도 영향을 줍니다. 이전 협업 방식(REST API)에서는
	프론트엔드 프로그래머는 백엔드 프로그래머가 작성하여 전달하는 API의 requset / response의 형식에 의존하게 됩니다. 그러나, gql을 사용한 방식에서는
	이러한 의존도가 많이 사라집니다. 다만 여전히 데이터 schema에 대한 협업 의존성은 존재합니다.
	*/

	/*
		! 오브젝트 타입과 필드

		type Character {
			name: String!
			appearsIn: [Episode!]!
		}

		- 오브젝트 타입: Character
		- 필드: name, appearsIn
		- 스칼라 타입: String, ID, Int 등
		- 느낌표(!): 필수 값을 의미 (non-nullable)
		- 대괄호([, ]): 배열을 의미 (array)
	*/

	/*
		! 리졸버 (resolver)

		데이터베이스 사용시, 데이터를 가져오기 위해서 sql을 작성했습니다. 또한 데이터베이스에는 데이터베이스 어플리케이션을 사용하여
		데이터를 가져오는 구체적인 과정이 구현되어 있습니다. 그러나 gql에서는 데이터를 가져오는 구체적인 과정을 직접 구현해야 합니다. gql 쿼리문
		파싱은 대부분의 gql 라이브러리에서 처리를 하지만, gql에서 데이터를 가져오는 구체적인 과정은 resolver(이하 리졸버)가 담당하고, 이를
		직접 구현해야 합니다. 프로그래머는 리졸버를 직접 구현해야하는 부담이 있지만, 이를 통해서 데이터 source의 종류에 상관업싱 구현이 가능합니다.
		예를 들어서, 리졸버를 통해 데이터를 데이터베이스에서 가져올 수 있고, 일반 파일에서 가져올 수 있고, 심지어 http, SOAP와 같은 네트워크 프로토콜을
		활용해서 원격 데이터를 가져올 수 있습니다. 덧붙이면, 이러한 특성을 이용하면 legacy 시스템을 gql 기반으로 바꾸는데 활용할 수 있습니다.

		gql 쿼리에서는 각각의 필드마다 함수가 하나씩 존재한다고 생각하면 됩니다. 이 함수는 다음 타입을 반환합니다. 이러한 각각의 함수를
		리졸버(resolver)라고 합니다. 만약 필드가 스칼라값(문자열이나 숫자와 같은 primitive 타입)인 경우에는 실행이 종료됩니다.
		즉 더 이상의 연쇄적인 리졸버 호출이 일어나지 않습니다. 하지만 필드의 타입이 스칼라 타입이 아닌 우리가 정의한 타입이라면 해당 타입의 리졸버를 호출되게 됩니다.

		이러한 연쇄적 리졸버 호출은 DFS(Depth First Search)로 구현되어 있을 것으로 추측합니다. 이점이 바로 gql이 Graph라는 단어를 쓴 이유가
		아닐까 생각합니다. 연쇄 리졸버 호출은 여러모로 장점이 있습니다. 연쇄 리졸버 특성을 잘 활용하면 DBMS의 관계에 대한 쿼리를 매우 쉽고, 효율적으로 처리할 수 있습니다.
		예를 들어, gql의 query에서 어떤 타입의 필드 중 하나가 해당 타입과 1:n 관계를 맺고 있다고 가정해보겠습니다.

		type Query {
			users: [User]
			user(id: ID): User
			limits: [Limit]
			limit(UserId: ID): Limit
			paymentsByUser(userId: ID): [Payment]
		}

		여기에서는 User와 Limit의 관계는 1:1이고 User와 Payment는 1:n 관계입니다.

		type User {
			id: ID!
			name: String!
			sex: SEX!
			birthDay: String!
			phoneNumber: String!
		}

		type Limit {
			id: ID!
			UserId: ID
			max: Int!
			amount: Int
			user: User
		}

		type Payment {
			id: ID!
			limit: Limit!
			user: User!
			pg: PaymentGateway!
			productName: String!
			amount: Int!
			ref: String
			createdAt: String!
			updateAt: String!
		}
	*/

	/*
		! Request (요청) - 쿼리

		{
			paymentsByUser(userId: 10) {
				id
				amount
			}
		}

		{
			paymentsByUser(userId: 10) {
				id
				amount
				user {
					name
					phoneNumber
				}
			}
		}

		위 두 쿼리는 동일한 쿼리명을 가지고 있지만, 호출되는 리좁버 함수의 갯수는 아래가 더 많습니다. 각각의 리졸버 함수에는 내부적으로
		데이터베이스 쿼리가 존재합니다. 이 말인즉, 쿼리에 맞게 필요한 만큼만 최적화하여 호출할 수 있다는 의미입니다. 내부적으로 로직 설계를 어떻게 하느냐에
		따라서 달라질 수 있겠지만 이러한 재귀형의 리졸버 체인을 잘 활용한다면 효율적인 설계가 가능합니다. (기존에 REST API 시대에는 정해진 쿼리는
		무조건 전부 호출되었습니다.)
	*/

	/* 리졸버 함수는 다음과 같이 총 4개의 인자를 받습니다.

	Query: {
		paymentsByUser: async (parent, { userId }, context, info) => {
			const limit = await Limit.findOne({ where: { UserId: userId }})
			const payments = await Payment.findAll({ where: { LimitId: limit.id }})
			return payments
		}
	},
	Payment: {
		limit: async (payment, args, context, info) => {
			return await Limit.findOne({ where: { id: payment.LimitId }})
		}
	}

	- 첫번째 인자는 parent로 연쇄적 리졸버 호출에서 부모 리졸버가 리턴한 객체입니다. 이 객체를 활용해서 현재 리졸버가 내보낼 값을 조절합니다.
	- 두번째 인자는 args로 쿼리에서 입력으로 넣은 인자입니다.
	- 세번째 인자는 context로 모든 리졸버에게 전달이 됩니다. 주로 미들웨어를 통해 입력된 값들이 들어 있습니다. 로그인 정보 혹은 권한과 같이 주요
	컨텍스트 관련 정보를 가지고 있습니다.
	- 네번째 인자는 info로 스키마 정보와 더불어 현재 쿼리의 특정 필드 정보를 가지고 있습니다. 잘 사용하지 않는 필드입니다.
	*/

	/*
		! 인트로스펙션 (introspection)

		기존 서버-클라이언트 협업 방식에서는 연동 규격서라고 하는 API 명세서를 주고 받는 절차가 반드시 필요했습니다. 프로젝트 관리 측면에서 관리해야 할 대상의
		증가는 작업의 복잡성 및 효율성 저해를 의미합니다.. 이 API 명세서는 때때로 관리가 제대로 되지 않아, 인터페이스 변경 사항을 제때 문서에 반영하지 못하기도 하고, 제 타이밍에 전달 못하곤 합니다.

		이러한 REST의 API 명세서 공유와 같은 문제를 해결하는 것이 gql의 인트로스펙션 기능입니다. gql의 인트로스펙션은 서버 자체에서 현재 서버에 정의된 스키마의 실시간 정보를 공유할 수 있게
		합니다. 이 스키마 정보만 알고 있으면 클라이언트 사이드에서는 따로 연동규격서를 요청할 필요가 없게 됩니다. 클라이언트 사이드에서는 실시간으로 현재 서버에서 정의하고 있는 스키마를 의심할 필요 없이 받아들이고,
		그에 맞게 쿼리문을 작성하면 됩니다.

		이러한 인트로스펙션용 쿼리가 따로 존재합니다. 일반 gql 쿼리문을 작성하면 됩니다. 다만 실제로는 굳이 스키마 인트로스펙션을 위해 gql 쿼리문을 작성할 필요가 없습니다. 대부분의 서버용 gql 라이브러리에는 쿼리용 IDE
		를 제공합니다.
	*/

	/*
		! GraphQL을 활용할 수 있게 도와주는 다양한 라이브러리들

		gql 자체는 쿼리 언어입니다. 이것만으로 할 수 있는 것이 없습니다. gql을 실제 구체적으로 활용할 수 있도록 도와주는 라이브러들이 몇가지 존재합니다. gql 자체는 개발 언어와 사용 네트워크에
		완전히 독립적입니다. 이를 어떻게 활용할지는 여러분에게 달려 있습니다.

		대표적인 gql 라이브러리 셋에 대한 링크는 2개를 소개합ㄴ디ㅏ.

		릴레이는 GraphQL의 어머니인 Facebook이 만들었습니다. 하지만 개인적인 의견으로는 현재(2019년 7월)버전의 릴레이는 사용하기
		매우 번거롭게 디자인되어 있다고 생각합니다. 개인적으로는 아폴로가 사용하기 편했습니다.

		- 릴레이 (Relay)
		- 아폴로 (Apollo GraphQL)
	*/

	/* 스키마, 리졸버를 비롯한 일반적인 GraphQL 용어

	API 개발자는 GraphQL을 사용해 클라이언트가 서비스를 통해 쿼리할 가능성이 있는 모든 데이터를 설명하는 스키마를 생성합니다.
	GraphQL 스키마는 개체 유형으로 구성되어 어떤 종류의 개체를 요청할 수 있으며 어떠한 필드가 있는지 정의합니다. 쿼리가 수신되면
	GraphQL은 스키마에 대해 쿼리를 검증하고 그 다음 검증된 쿼리를 실행합니다.

	API 개발자는 스키마의 각 필드를 리졸버라고 불리는 기능에 첨부합니다. 실행 중 값을 생산하기 위해 리졸버가 호출됩니다.
	*/

	// 객체 필드는 해당 객체에 대한 데이터나 다른 객체와의 연결에 대한 데이터를 보여준다.
	fields := graphql.Fields{
		"hello": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "world", nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}

	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	query := `
					{
								hello
					}
	`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("failed to execute graphql operation, errors: %+v", r.Errors)
	}
	// fmt.Printf("%s \n", r)
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}

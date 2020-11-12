package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/machinebox/graphql"
	"gitlab.com/shikho/buffers/protobuff"
	"io"
	"log"
	"os"
)

func main() {

	// create a client (safe to share across requests)
	client := graphql.NewClient("http://localhost:4000/graphql")

	// make a request
	req := graphql.NewRequest(`
		mutation MyMutation($c : ClassEnum!, $n: String!, $g: StudyGroupTypeEnum, $a: SubjectAttributeTypeEnum) {
		  insertSubjectCode(class: $c, name: $n, group: $g, attribute: $a) {
			code
			display
			ref
			attr
		  }
		}
	`)

	// set header fields
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIyODg1ODAiLCJlbWFpbCI6IiIsImV4cCI6MTU5NDE0ODY2MCwiaWF0IjoxNTkxNTU2NjYwLCJqdGkiOiJrWmh3MHdMZ1VLN2Q3R3E2Y0xrcTZzVmloMXp5bzQ1WiIsInN1YiI6IjI4ODU4MCJ9.G-1TzMmJcodaZO9ifZxov54yabxHN4UcGHd5gwsaXqIGiNYGgck7FyS3dv3sLTUfLzfPPSuJPXlEq5wQkFYm9bxphGRa0Xqq-GV4t3RV4eA7VSGPFx0TiF5RP64Pm4M02soDLHTBPdOVdLXf0o_C-D03ut397rpvBQlPJL80UYd2ZDPlFGM_hAUGGuCcRA4i_SetGjTQ8DuyAofdozdkJxmUvpqBNgQr6gJj_A9WorQ2gZA8uTzfAxjsIw6bF7dI-cfdqstMi_QqRdZWWIatkPl3wWnZmNFyEn2_sZ6h0uau8kRHR5BvYnPkdIIw66Z-J6cUB_No4fsxeKFOnBjcXA")

	// define a Context for the request
	ctx := context.Background()

	csvFile, err := os.Open("subjects.csv")
	if err != nil {
		panic(err.Error())
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {

		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		// set any variables
		req.Var("c", fmt.Sprintf(`C%s`, line[0]))
		req.Var("n", line[2])
		req.Var("g", line[1])
		req.Var("a", line[3])

		// run it and capture the response
		respData := map[string]protobuff.CodeSystem{
			"insertSubjectCode": protobuff.CodeSystem{},
		}
		if err := client.Run(ctx, req, &respData); err != nil {
			fmt.Println(err)
		}

		fmt.Println(respData)

	}

}

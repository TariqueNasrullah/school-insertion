package main

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/machinebox/graphql"
	"gitlab.com/shikho/buffers/protobuff"
)

func main() {

	// create a client (safe to share across requests)
	//client := graphql.NewClient("http://localhost:5000/graphql")
	client := graphql.NewClient("https://api.shikho.net/graphql")

	// make a request
	req := graphql.NewRequest(`
      mutation AreaInsert($area: String, $district: String, $division: String) {
	  insertAddressCode(area: $area, district: $district, division: $division) {
		division {
		  code
		  display
		  ref
		}
		district {
		  code
		  display
		  ref
		}
		area {
		  code
		  display
		  ref
		}
	  }
	}`)

	// make a request
	req2 := graphql.NewRequest(`
    mutation SchoolInsert($area_id: String!, $eiin: String, $name: String) {
	  upsertSchool(area_id: $area_id, eiin: $eiin, name: $name) {
		eiin_no
		id
		name
	  }
	}`)

	// set header fields
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIyODg1ODAiLCJlbWFpbCI6ImFkbWluQHNoaWtoby50ZWNoIiwiZXhwIjoxNjA2NTU0MDA0LCJpYXQiOjE2MDM5NjIwMDQsImp0aSI6IlpNcWw2RVBsTjh3MU1QTFpPWE9CUU5zR0FHaXJHTFRXIiwicm9sZSI6ImFkbWluIiwic3ViIjoiMjg4NTgwIiwidXNlcl9pZCI6IjI4ODU4MCIsInVzZXJfcGhvbmUiOiIwMTc2MDAwMDAwMCJ9.g46hYiMLoJ0AvrBuR5H9HWTJm_jUVRSFdF0RYZwGbDAzqY5oeDJ0Y-pJePeXW_kotTbdbUz9XFt6Lf-UTPFpyJmqQOtamgoNvX_zYbla-j-m_5MhrbvgSxBOHPkxDgXoJ8bdT9nymTt0mBu-c1WZJVQZE1OkXLy6foUOnkBNkgrnDf8Hxtbu4GKVApqzAeV9wAdoXTt4aCe_ksTGBLJLOPHUbvNG5fvpgKnZOhqG250wgI2-m6Zn7gDaOYxOEiu8UElB1kcqwgr8jhLO-Xu-lW15fSwK9IE9DDkw36h2_doP5j0yCBS1FTNGXY0xfYfC6JzkZ5rKMeX24UIgOYSWxg")

	req2.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIyODg1ODAiLCJlbWFpbCI6ImFkbWluQHNoaWtoby50ZWNoIiwiZXhwIjoxNjA2NTU0MDA0LCJpYXQiOjE2MDM5NjIwMDQsImp0aSI6IlpNcWw2RVBsTjh3MU1QTFpPWE9CUU5zR0FHaXJHTFRXIiwicm9sZSI6ImFkbWluIiwic3ViIjoiMjg4NTgwIiwidXNlcl9pZCI6IjI4ODU4MCIsInVzZXJfcGhvbmUiOiIwMTc2MDAwMDAwMCJ9.g46hYiMLoJ0AvrBuR5H9HWTJm_jUVRSFdF0RYZwGbDAzqY5oeDJ0Y-pJePeXW_kotTbdbUz9XFt6Lf-UTPFpyJmqQOtamgoNvX_zYbla-j-m_5MhrbvgSxBOHPkxDgXoJ8bdT9nymTt0mBu-c1WZJVQZE1OkXLy6foUOnkBNkgrnDf8Hxtbu4GKVApqzAeV9wAdoXTt4aCe_ksTGBLJLOPHUbvNG5fvpgKnZOhqG250wgI2-m6Zn7gDaOYxOEiu8UElB1kcqwgr8jhLO-Xu-lW15fSwK9IE9DDkw36h2_doP5j0yCBS1FTNGXY0xfYfC6JzkZ5rKMeX24UIgOYSWxg")

	// define a Context for the request
	ctx := context.Background()

	csvFile, err := os.Open("final_school - formated.csv")
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
		req.Var("area", line[2])
		req.Var("district", line[1])
		req.Var("division", line[0])

		// run it and capture the response
		respData := map[string]protobuff.Address{
			"insertAddressCode": {},
		}
		if err := client.Run(ctx, req, &respData); err != nil {
			fmt.Println(err)
		}

		// set any variables
		req2.Var("area_id", respData["insertAddressCode"].Area.Code)
		req2.Var("name", line[3])
		req2.Var("eiin", line[4])

		// run it and capture the response
		var resp interface{}
		if err := client.Run(ctx, req2, &resp); err != nil {
			fmt.Println(err)
		}

		fmt.Println(resp)

	}

}

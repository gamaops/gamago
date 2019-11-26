/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

type OperationItem struct {
	Method string           `json:"method"`
	Path   string           `json:"path"`
	Body   *json.RawMessage `json:"body"`
}

type OperationsContent struct {
	Operations []*OperationItem `json:"operations"`
}

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "This command will apply an Elasticsearch operations file",
	Long:  `This command will apply an Elasticsearch operations file`,
	Run: func(cmd *cobra.Command, args []string) {
		jsonFile, err := os.Open(args[0])
		if err != nil {
			log.Fatal(err)
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatal(err)
		}

		var opsContent OperationsContent

		err = json.Unmarshal(byteValue, &opsContent)
		if err != nil {
			log.Fatal(err)
		}

		for opIndex, opItem := range opsContent.Operations {
			log.Infof("Performing operation #%v", opIndex)
			reqReader := bytes.NewReader(*opItem.Body)
			req, err := http.NewRequest(opItem.Method, opItem.Path, reqReader)
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Add("Content-Type", "application/json")
			res, err := ElasticClient.Perform(req)
			if err != nil {
				log.Fatal(err)
			}
			if res.StatusCode >= 300 {
				bodyBytes, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Errorf("Error reading Elasticsearch response: %v", err)
				}
				log.Fatalf("Elasticsearch response error for operation #%v\n%v", opIndex, string(bodyBytes))
			}
			res.Body.Close()
		}

	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

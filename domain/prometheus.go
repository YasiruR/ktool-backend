package domain

type PromResponse struct {
	Status 		string		`json:"status"`
	Data 		struct {
		ResultType 		string		`json:"resultType"`
		Result 			[] struct {
			Metric 			struct {
				Name 			string		`json:"__name__"`
				Instance 		string		`json:"instance"`
				Job 			string		`json:"job"`
			}	`json:"metric"`
			Value 		[2]interface{}		`json:"value"`
			Values 		[][2]interface{}	`json:"values"`
		}
	} `json:"data"`
}

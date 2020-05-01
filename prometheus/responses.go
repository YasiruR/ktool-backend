package prometheus

type BrokerBytes struct {
	Status 		string 	`json:"status"`
	Data   		struct {
		ResultType 	string 	`json:"resultType"`
		Result     	[]struct {
			Metric 		struct {
				Instance 	string 	`json:"instance"`
				Job      	string 	`json:"job"`
			} `json:"metric"`
			Value 		[]interface{} `json:"value"`	//consists of ts (float) and value (string)
		} `json:"result"`
	} `json:"data"`
	ErrorType 	string 	`json:"error_type"`
	Error 		string	`json:"error"`
	Warnings 	string	`json:"warnings"`
}

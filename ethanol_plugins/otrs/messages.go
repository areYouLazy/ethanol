package main

type AuthenticationBody struct {
	UserLogin string `json:"UserLogin"`
	Password  string `json:"Password"`
}

// {
// 		"TicketID": [
// 			"1",
// 			"2",
// 			"3"
// 		]
// }

type TicketSearchResponse struct {
	TicketID []string `json:"TicketID"`
}

// {
//     "Ticket": [
//         {
//             "Age": 12022,
//             "PriorityID": "3",
//             "ServiceID": "",
//             "Type": "Unclassified",
//             "Responsible": "root@localhost",
//             "StateID": "1",
//             "ResponsibleID": "1",
//             "ChangeBy": "1",
//             "EscalationTime": "0",
//             "OwnerID": "1",
//             "Changed": "2024-08-08 14:52:23",
//             "TimeUnit": 0,
//             "RealTillTimeNotUsed": "0",
//             "GroupID": "1",
//             "Owner": "root@localhost",
//             "CustomerID": null,
//             "TypeID": 1,
//             "Created": "2024-08-08 14:52:23",
//             "Priority": "3 normal",
//             "UntilTime": 0,
//             "EscalationUpdateTime": "0",
//             "Queue": "Raw",
//             "QueueID": "2",
//             "State": "new",
//             "Title": "Znuny says hi!",
//             "CreateBy": "1",
//             "TicketID": "1",
//             "StateType": "new",
//             "UnlockTimeout": "0",
//             "EscalationResponseTime": "0",
//             "EscalationSolutionTime": "0",
//             "LockID": "1",
//             "TicketNumber": "2021012710123456",
//             "ArchiveFlag": "n",
//             "Lock": "unlock",
//             "SLAID": "",
//             "CustomerUserID": null
//         }
//     ]
// }

type TicketResponse struct {
	Ticket []struct {
		Age            string `json:"Age"`
		PriorityID     string `json:"PriorityID"`
		ServiceID      string `json:"ServiceID"`
		Type           string `json:"Unclassified"`
		Responsible    string `json:"Responsible"`
		StateID        string `json:"StateID"`
		ResponsibleID  string `json:"ResponsibleID"`
		ChangeBy       string `json:"ChangeBy"`
		EscalationTime string `json:"EscalationTime"`
		OwnerID        string `json:"OwnerID"`
		Changed        string `json:"Changed"`
		Owner          string `json:"Owner"`
		State          string `json:"State"`
		Title          string `json:"Title"`
		TicketID       string `json:"TicketID"`
		StateType      string `json:"StateType"`
		TicketNumber   string `json:"TicketNumber"`
	}
}

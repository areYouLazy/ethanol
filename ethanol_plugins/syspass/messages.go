package main

// rawResponse maps response from syspass backend
type rawResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		ItemID int `json:"itemId"`
		Result []struct {
			ID                        int    `json:"id"`
			UserID                    int    `json:"userId"`
			UserGroupID               int    `json:"userGroupId"`
			UserEditId                int    `json:"userEditId"`
			Name                      string `json:"name"`
			ClientID                  int    `json:"clientId"`
			CategoryID                int    `json:"categoryId"`
			Login                     string `json:"login"`
			URL                       string `json:"url"`
			Notes                     string `json:"notes"`
			OtherUserEdit             int    `json:"otherUserEdit"`
			OtherUserGroupEdit        int    `json:"otherUserGroupEdit"`
			IsPrivate                 int    `json:"isPrivate"`
			IsPrivateGroup            int    `json:"isPrivateGroup"`
			DateEdit                  int    `json:"dateEdit"`
			PassDate                  int    `json:"passDate"`
			PassDateChange            int    `json:"passDateChange"`
			ParentID                  int    `json:"parentId"`
			CategoryName              string `json:"categoryName"`
			ClientName                string `json:"clientName"`
			UserGroupName             string `json:"userGroupname"`
			UserName                  string `json:"userName"`
			UserLogin                 string `json:"userLogin"`
			UserEditName              string `json:"userEditName"`
			UserEditLogin             string `json:"userEditLogin"`
			NumFiles                  int    `json:"num_files"`
			PublicLinkHash            string `json:"publicLinkHash"`
			PublicLinkDateExpire      int    `json:"publicLinkDateExpire"`
			PublicLinkTotalCountViews int    `json:"publicLinkTotalCountViews"`
			CountView                 int    `json:"countView"`
		}
		ResultCode    int    `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
		Count         int    `json:"count"`
	}
	ID int `json:"id"`
}

// rawPasswordResponse maps response from backend for account/viewPass query
type rawPasswordResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  struct {
		ItemID int `json:"itemId"`
		Result struct {
			Password string `json:"password"`
		}
		ResultCode    int    `json:"resultCode"`
		ResultMessage string `json:"resultMessage"`
		Count         int    `json:"count"`
	}
	ID int `json:"id"`
}

// queryBodyParams structure to map body parameters for query
type queryBodyParams struct {
	AuthToken string `json:"authToken"` // api key
	TokenPass string `json:"tokenPass"` // password used to generate api key
	Text      string `json:"text"`      // text to search for
	Count     int    `json:"count"`     // number of results
	ID        int    `json:"id"`        // id of a syspass account, for account/view action
	Details   int    `json:"details"`   // 1 to get account details
}

// queryBody structure to map body for query
type queryBody struct {
	JSONRPC string          `json:"jsonrpc"` // jsonrpc version, default to 2.0
	Method  string          `json:"method"`  // action to take, default to account/search
	Params  queryBodyParams `json:"params"`  // query params
	ID      int             `json:"id"`      // query ID
}

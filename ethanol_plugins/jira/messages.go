package main

// {
//     "expand": "names,schema",
//     "startAt": 0,
//     "maxResults": 50,
//     "total": 1,
//     "issues": [
//         {
//             "expand": "operations,versionedRepresentations,editmeta,changelog,renderedFields",
//             "id": "10016",
//             "self": "http://192.168.56.104:32778/rest/api/latest/issue/10016",
//             "key": "ET-17",
//             "fields": {
//                 "issuetype": {
//                     "self": "http://192.168.56.104:32778/rest/api/2/issuetype/10002",
//                     "id": "10002",
//                     "description": "A task that needs to be done.",
//                     "iconUrl": "http://192.168.56.104:32778/secure/viewavatar?size=xsmall&avatarId=10318&avatarType=issuetype",
//                     "name": "Task",
//                     "subtask": false,
//                     "avatarId": 10318
//                 },
//                 "components": [],
//                 "timespent": null,
//                 "timeoriginalestimate": null,
//                 "description": "This is Dummy Host on Jira!",
//                 "project": {
//                     "self": "http://192.168.56.104:32778/rest/api/2/project/10000",
//                     "id": "10000",
//                     "key": "ET",
//                     "name": "ethanol",
//                     "projectTypeKey": "software",
//                     "avatarUrls": {
//                         "48x48": "http://192.168.56.104:32778/secure/projectavatar?avatarId=10324",
//                         "24x24": "http://192.168.56.104:32778/secure/projectavatar?size=small&avatarId=10324",
//                         "16x16": "http://192.168.56.104:32778/secure/projectavatar?size=xsmall&avatarId=10324",
//                         "32x32": "http://192.168.56.104:32778/secure/projectavatar?size=medium&avatarId=10324"
//                     }
//                 },
//                 "fixVersions": [],
//                 "customfield_10110": null,
//                 "customfield_10111": null,
//                 "aggregatetimespent": null,
//                 "resolution": null,
//                 "customfield_10104": null,
//                 "customfield_10105": "0|i0003j:",
//                 "customfield_10107": null,
//                 "customfield_10108": null,
//                 "aggregatetimeestimate": null,
//                 "customfield_10109": null,
//                 "resolutiondate": null,
//                 "workratio": -1,
//                 "summary": "Dummy Host",
//                 "lastViewed": null,
//                 "watches": {
//                     "self": "http://192.168.56.104:32778/rest/api/2/issue/ET-17/watchers",
//                     "watchCount": 1,
//                     "isWatching": false
//                 },
//                 "creator": {
//                     "self": "http://192.168.56.104:32778/rest/api/2/user?username=emil",
//                     "name": "emil",
//                     "key": "JIRAUSER10000",
//                     "emailAddress": "emil@pandocchi.it",
//                     "avatarUrls": {
//                         "48x48": "https://www.gravatar.com/avatar/b45f3c9dfdf20c5dd69cdb6cf529abcf?d=mm&s=48",
//                         "24x24": "https://www.gravatar.com/avatar/b45f3c9dfdf20c5dd69cdb6cf529abcf?d=mm&s=24",
//                         "16x16": "https://www.gravatar.com/avatar/b45f3c9dfdf20c5dd69cdb6cf529abcf?d=mm&s=16",
//                         "32x32": "https://www.gravatar.com/avatar/b45f3c9dfdf20c5dd69cdb6cf529abcf?d=mm&s=32"
//                     },
//                     "displayName": "Emil",
//                     "active": true,
//                     "timeZone": "GMT"
//                 },
//                 "subtasks": [],
//                 "created": "2024-08-09T13:42:08.000+0000",
//                 "reporter": {
//                     "self": "http://192.168.56.104:32778/rest/api/2/user?username=emil",
//                     "name": "emil",
//                     "key": "JIRAUSER10000",
//                     "emailAddress": "emil@pandocchi.it",
//                     "avatarUrls": {
//                         "48x48": "https://www.gravatar.com/avatar/b45f3c9dfdf20c5dd69cdb6cf529abcf?d=mm&s=48",
//                         "24x24": "https://www.gravatar.com/avatar/b45f3c9dfdf20c5dd69cdb6cf529abcf?d=mm&s=24",
//                         "16x16": "https://www.gravatar.com/avatar/b45f3c9dfdf20c5dd69cdb6cf529abcf?d=mm&s=16",
//                         "32x32": "https://www.gravatar.com/avatar/b45f3c9dfdf20c5dd69cdb6cf529abcf?d=mm&s=32"
//                     },
//                     "displayName": "Emil",
//                     "active": true,
//                     "timeZone": "GMT"
//                 },
//                 "customfield_10000": "{summaryBean=com.atlassian.jira.plugin.devstatus.rest.SummaryBean@95d79902[summary={pullrequest=com.atlassian.jira.plugin.devstatus.rest.SummaryItemBean@e707cad3[overall=PullRequestOverallBean{stateCount=0, state='OPEN', details=PullRequestOverallDetails{openCount=0, mergedCount=0, declinedCount=0}},byInstanceType={}], build=com.atlassian.jira.plugin.devstatus.rest.SummaryItemBean@b1a06565[overall=com.atlassian.jira.plugin.devstatus.summary.beans.BuildOverallBean@d8dedc1b[failedBuildCount=0,successfulBuildCount=0,unknownBuildCount=0,count=0,lastUpdated=<null>,lastUpdatedTimestamp=<null>],byInstanceType={}], review=com.atlassian.jira.plugin.devstatus.rest.SummaryItemBean@e1f79304[overall=com.atlassian.jira.plugin.devstatus.summary.beans.ReviewsOverallBean@d495e744[stateCount=0,state=<null>,dueDate=<null>,overDue=false,count=0,lastUpdated=<null>,lastUpdatedTimestamp=<null>],byInstanceType={}], deployment-environment=com.atlassian.jira.plugin.devstatus.rest.SummaryItemBean@a91a83c5[overall=com.atlassian.jira.plugin.devstatus.summary.beans.DeploymentOverallBean@38062fc7[topEnvironments=[],showProjects=false,successfulCount=0,count=0,lastUpdated=<null>,lastUpdatedTimestamp=<null>],byInstanceType={}], repository=com.atlassian.jira.plugin.devstatus.rest.SummaryItemBean@af16ce7a[overall=com.atlassian.jira.plugin.devstatus.summary.beans.CommitOverallBean@c2c53102[count=0,lastUpdated=<null>,lastUpdatedTimestamp=<null>],byInstanceType={}], branch=com.atlassian.jira.plugin.devstatus.rest.SummaryItemBean@e2026490[overall=com.atlassian.jira.plugin.devstatus.summary.beans.BranchOverallBean@87667d17[count=0,lastUpdated=<null>,lastUpdatedTimestamp=<null>],byInstanceType={}]},errors=[],configErrors=[]], devSummaryJson={\"cachedValue\":{\"errors\":[],\"configErrors\":[],\"summary\":{\"pullrequest\":{\"overall\":{\"count\":0,\"lastUpdated\":null,\"stateCount\":0,\"state\":\"OPEN\",\"details\":{\"openCount\":0,\"mergedCount\":0,\"declinedCount\":0,\"total\":0},\"open\":true},\"byInstanceType\":{}},\"build\":{\"overall\":{\"count\":0,\"lastUpdated\":null,\"failedBuildCount\":0,\"successfulBuildCount\":0,\"unknownBuildCount\":0},\"byInstanceType\":{}},\"review\":{\"overall\":{\"count\":0,\"lastUpdated\":null,\"stateCount\":0,\"state\":null,\"dueDate\":null,\"overDue\":false,\"completed\":false},\"byInstanceType\":{}},\"deployment-environment\":{\"overall\":{\"count\":0,\"lastUpdated\":null,\"topEnvironments\":[],\"showProjects\":false,\"successfulCount\":0},\"byInstanceType\":{}},\"repository\":{\"overall\":{\"count\":0,\"lastUpdated\":null},\"byInstanceType\":{}},\"branch\":{\"overall\":{\"count\":0,\"lastUpdated\":null},\"byInstanceType\":{}}}},\"isStale\":false}}",
//                 "aggregateprogress": {
//                     "progress": 0,
//                     "total": 0
//                 },
//                 "priority": {
//                     "self": "http://192.168.56.104:32778/rest/api/2/priority/3",
//                     "iconUrl": "http://192.168.56.104:32778/images/icons/priorities/medium.svg",
//                     "name": "Medium",
//                     "id": "3"
//                 },
//                 "customfield_10100": null,
//                 "labels": [],
//                 "environment": null,
//                 "timeestimate": null,
//                 "aggregatetimeoriginalestimate": null,
//                 "versions": [],
//                 "duedate": null,
//                 "progress": {
//                     "progress": 0,
//                     "total": 0
//                 },
//                 "issuelinks": [],
//                 "votes": {
//                     "self": "http://192.168.56.104:32778/rest/api/2/issue/ET-17/votes",
//                     "votes": 0,
//                     "hasVoted": false
//                 },
//                 "assignee": null,
//                 "updated": "2024-08-09T13:42:08.000+0000",
//                 "status": {
//                     "self": "http://192.168.56.104:32778/rest/api/2/status/10000",
//                     "description": "",
//                     "iconUrl": "http://192.168.56.104:32778/",
//                     "name": "Backlog",
//                     "id": "10000",
//                     "statusCategory": {
//                         "self": "http://192.168.56.104:32778/rest/api/2/statuscategory/2",
//                         "id": 2,
//                         "key": "new",
//                         "colorName": "blue-gray",
//                         "name": "To Do"
//                     }
//                 }
//             }
//         }
//     ]
// }

type RawResponse struct {
	Issues []struct {
		Key    string `json:"key"`
		Fields struct {
			Description string `json:"description"`
			Summary     string `json:"summary"`
			Creator     struct {
				Name         string `json:"name"`
				EmailAddress string `json:"emailAddress"`
			}
			Created string `json:"created"`
		}
	}
}

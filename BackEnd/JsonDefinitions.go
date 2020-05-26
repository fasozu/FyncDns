package BackEnd

import (
)

//Output Json definition for CheckServer endpoint
type OutputServer struct {
    Address string `json:"address"`
    Ssl_grade string `json:"ssl_grade"`
    Country string `json:"country"`
    Owner string `json:"owner"`     
}

type OutputResponse struct {
	Success bool `json:"success"`
	ErrorMessage string `json:"errorMessage"`
	Servers []OutputServer `json:"servers"`
	Servers_changed bool `json:"servers_changed"`
	Ssl_grade string `json:"ssl_grade"`
	Logo string `json:"logo"`
	Is_down bool `json:"is_down"`
	StatusCode string `json:"statusCode"`
	Title string `json:"title"`
}


//Input Json definition por api.sllabs
type InputEndpoint struct {
	IpAddress string `json:"ipAddress"`
	ServerName string `json:"serverName"`
	StatusMessage string `json:"statusMessage"`
	Grade string `json:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored"`
	HasWarnings bool `json:"hasWarnings"`
	IsExceptional bool `json:"isExceptional"`
	Progress int `json:"progress"`
	Duration int `json:"duration"`
	Delegation int `json:"delegation"`
}

type InputResponse struct {
	Host string `json:"host"`
	Port int `json:"port"`
	Protocol string `json:"protocol"`
	IsPublic bool `json:"isPublic"`
	Status string `json:"status"`
	StartTime int `json:"startTime"`
	TestTime int `json:"testTime"`
	EngineVersion string `json:"engineVersion"`
	CriteriaVersion string `json:"criteriaVersion"`
	Servers []OutputServer `json:"servers"`
	Endpoints []InputEndpoint `json:"Endpoints"`
}



	


package models

type AvbDemand struct {
	PmpAdspaceId     int
	DemandAdspaceId  int
	PmpAdspaceKey    string
	DemandAdspaceKey string
	PlanImp          int
	PlanClk          int
	ActualImp        int
	ActualClk        int
}

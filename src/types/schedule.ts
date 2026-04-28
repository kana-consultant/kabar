export interface ScheduleRequest {
	Title          : string   
	Topic          : string   
	Article        : string   
	ImageURL?       : string   
	ImagePrompt    : string   
	ScheduledFor   : string   
	TargetProducts : string[] 
	HasImage?       : boolean     
}

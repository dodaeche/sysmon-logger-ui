package main

import (
	"time"
	"html/template"
)

// Base fields that all database tables support
type Base struct {
	Id      	int64 		`db:"id"`
	Domain  	string 		`db:"domain"`
	Host    	string 		`db:"host"`
	UtcTime		time.Time	`db:"utc_time"`
	UtcTimeStr	string		`db:"-"`
}

// Represents an "event" record
type Event struct {
	Base
	Type		string			`db:"type"`
	Message		string			`db:"message"`
	MessageHtml	template.HTML	`db:"_"`
}

// Represents an "process_create" record
type ProcessCreate struct {
	Base
	ProcessId       	int64 			`db:"process_id"`
	Image 				string 			`db:"image"`
	CommandLine			string 			`db:"command_line"`
	CurrentDirectory 	string 			`db:"current_directory"`
	Md5 				string 			`db:"md5"`
	Sha256 				string 			`db:"sha256"`
	ParentProcessId 	int64 			`db:"parent_process_id"`
	ParentImage 		string 			`db:"parent_image"`
	ParentCommandLine	string 			`db:"parent_command_line"`
	ProcessUser			string 			`db:"process_user"`
	OtherData			template.HTML	`db:"_"`
}

// Represents an "process_terminate" record
type ProcessTerminate struct {
	Base
	ProcessId	int64	`db:"process_id"`
	Image     	string  `db:"image"`
}

// Represents an "file_creation_time" record
type FileCreationTime struct {
	Base
	ProcessId					int64    		`db:"process_id"`
	Image     					string     		`db:"image"`
	TargetFileName  			string     		`db:"target_file_name"`
	CreationUtcTime  			time.Time		`db:"creation_utc_time"`
	PreviousCreationUtcTime 	time.Time   	`db:"previous_creation_utc_time"`
	OtherData					template.HTML	`db:"_"`
}

// Represents an "network_connection" record
type NetworkConnection struct {
	Base
	ProcessId       	int64 			`db:"process_id"`
	Image 				string 			`db:"image"`
	ProcessUser			string 			`db:"process_user"`
	Protocol			string 			`db:"protocol"`
	Initiated 			bool 			`db:"initiated"`
	SourceIp 			string 			`db:"source_ip"`
	SourceHostName 		string 			`db:"source_host_name"`
	SourcePort 			int32 			`db:"source_port"`				// Matches values for postgres (integer) e.g. Range: -2147483648 through 2147483647.
	SourcePortName 		string 			`db:"source_port_name"`
	DestinationIp		string 			`db:"destination_ip"`
	DestinationHostName	string 			`db:"destination_host_name"`
	DestinationPort		int32 			`db:"destination_port"`			// Matches values for postgres (integer) e.g. Range: -2147483648 through 2147483647.
	DestinationPortName	string 			`db:"destination_port_name"`
	OtherData			template.HTML	`db:"_"`
}

// Represents an "driver_loaded" record
type DriverLoaded struct {
	Base
	ImageLoaded	string  		`db:"image_loaded"`
	Md5			string			`db:"md5"`
	Sha256		string			`db:"sha256"`
	Signed		bool			`db:"signed"`
	Signature	string			`db:"signature"`
	OtherData	template.HTML	`db:"_"`
}

// Represents an "image_loaded" record
type ImageLoaded struct {
	Base
	ProcessId   int64  			`db:"process_id"`
	Image       string  		`db:"image"`
	ImageLoaded string  		`db:"image_loaded"`
	Md5         string    		`db:"md5"`
	Sha256      string    		`db:"sha256"`
	Signed      bool    		`db:"signed"`
	Signature   string    		`db:"signature"`
	OtherData	template.HTML	`db:"_"`
}

// Represents an "raw_access" record
type RawAccessRead struct {
	Base
	ProcessId	int64 	`db:"process_id"`
	Image     	string  `db:"image"`
	Device     	string  `db:"device"`
}

// Represents an "create_remote_thread" record
type CreateRemoteThread struct {
	Base
	SourceProcessId	int64  			`db:"source_process_id"`
	SourceImage		string  		`db:"source_image"`
	TargetProcessId	int64  			`db:"target_process_id"`
	TargetImage		string			`db:"target_image"`
	NewThreadId		int64  			`db:"new_thread_id"`
	StartAddress	string  		`db:"start_address"`
	StartModule		string			`db:"start_module"`
	StartFunction	string			`db:"start_function"`
	OtherData		template.HTML	`db:"_"`
}
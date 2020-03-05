# c2m
## Introduction
c2m -abberviation for csv to mqtt- is a general-purpose configurable command-line tool (CLI) for publishing data 
from CSV files to a MQTT message broker

## Usage
	c2m --config={path} {--dry-run|--health} 

E.g.

    mp.go dry-run -c config.yml -s ./files/LabView --pretty --watch

	
Command options

**config**	The full path to a configuration file

	# MQTT Settings
	MQTT_HOST={string}
	MQTT_USER=
	MQTT_PASSWORD=
	MQTT_TOPIC=

	# CSV Settings
	fileExt=csv
	fieldMappingsToJson={index or header title:jsonAttribute} E.g., 0:drumTest|1:name
	headerDetectionStart=name	
	report.heading = line no., separator, index of field, name 
				E.g. 1, ,3,serialNumber|2, ,3,partNumber
	report.heading.includeFields = comma separated field name
				E.g. serialNumber,partNumber

	report.data.heading.detection
	report.data.fieldMapping = 


**dry-run**	Processes the files and dumps the output to the console without publishing to the broker

**health**	Reports the status of the message broker
run	Processes the files and publishes to the message broker

canProcessFile
- Not a directory
- Matches the provided extension
- Has not already been processed

# Libraries
- https://fsnotify.org/
- https://github.com/ilyakaznacheev/cleanenv
- https://github.com/teris-io/cli

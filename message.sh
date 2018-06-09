#!/bin/bash

count=200
value=1
limit=$value+$count
limit=$(( count + value ))

while [ $value -lt $limit ]; do
	echo $value
	curl -X POST \
		http://localhost:8101/messages \
		-H 'Content-Type: application/json' \
		-H 'X-Sub: 616215fe-f43f-460e-ac33-7d47a7af6601' \
		-d '{"messages":
		[
			{
				"roomId":"6196136d-6ae7-4638-a3b0-3df081baa97b",
				"userId":"616215fe-f43f-460e-ac33-7d47a7af6601",
				"type":"text",
				"eventName":"message",
				"payload":{
					"text":"'$value\\n$value\\n$value'"
				}
			}
		]
	}'
	value=$(( value + 1 ))
done

# Golang Open AI Chat BOT 

The bot can retrieve photos using a username and if username is not available then it can create one also and save that in database.

## Installation

To run the project, navigate to the `cmd` directory and execute the following command:

bash
cd cmd
go run .



## Socket chat bot api 
ws://localhost:8765/v1/ws/user_chat

## Uploading Photos
To upload photos using the API, you can use cURL. Here's an example command:

curl --location 'http://localhost:8765/v1/upload_photos' \
--form 'images=@"/Users/username/Downloads/6935d6b06fee3002f712f852b48f3c95-original.jpeg"' \
--form 'images=@"/Users/username/Downloads/8f9a92fe241b9530ae8701eb9f5bb9ce-original.jpeg"'


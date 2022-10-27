## Description

VK bot which gets info about last played `dota 2` match of specified user and sends info about  
played match to VK chat. Bot keeps listening for new matches and if specified user finishes match,  
bot will notify about this match.  
Used repos: `https://github.com/l2x/dota2api` and `https://github.com/SevereCloud/vksdk/v2`

## Running the project

First of all you need to create VK group and generate TOKEN in group settings after export token  
to OS ENV as:
~~~
export TOKEN=your_token_value
~~~ 
You can add Bot from VK group which you created to Group chat or just in personal messages.  
Specify steamIds variable, just replace value in code and run the program.
~~~
go run ./cmd/api
~~~
Now if you write `start` in the chat with bot, it will start listening for new matches. 
# Secret Handshake Server
The code base for the reddit storage solution Secret Handshake server.

It is written in Go and is extremely messy as it was written during a 36hr Hackathon. We will probably continue to make changes as a learning exercise in go.


The website that should work with any local instance of the server (see below)

  https://luis-gzz.github.io/SecretHandshake/

The "database" we used during HackGT
  
   https://www.reddit.com/r/SecretHandshakeVault/
   
### Usage

To use you will need to replace all instances of /r/SecretHandshakeVault in the server code with the subreddit you would like to use. You will also need to setup a reddit script and .agent file (the code requires one named bot.agent). Since this project uses graw (Go Reddit API Wrapper) the instruction about registration on their [gitbook](https://turnage.gitbooks.io/graw/content/) will be helpful.

Place all files from this repository into a directory along with the .agent file you created. And from within that folder's directory run 
  `go run server.go encrypt.go` within the command line. The [SecretHandshake website](https://luis-gzz.github.io/SecretHandshake/) should be set up to point to localhost:3000 which the go script should automatically run in. 
  
From here you can upload text and images(no large image file support) and you can retrieve text; image retrieval sort of works but still needs some work.

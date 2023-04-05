# TCP server with POW DDoS protection

This is a simple implementation of a tcp server with Proof-of-Work protection(challenge-response protocol).

To build and start this project use default Makefile target, it will build client and server docker images and start them in foregraund with docker-compose.

Server useses the POW alghoritm inspired by HashCash. POW complexity is increased with the number of current active requests. 

There is a lot to improve here, but as a proof of concept, I think it's enough.
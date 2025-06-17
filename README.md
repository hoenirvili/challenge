
### Problem

A "peering relationship" is a way to keep track of debt between two parties. The peering relationship
balance starts at 0 and both parties independently track their views of the balance. 
If Alice owes Bob $10, then Alice sees a balance of -10 and Bob sees a balance of 10. 
If Alice sends Bob $10 more, then her balance would be -20, and Bob’s balance would be 20.

Write a program that implements a peering relationship and exposes an interactive Command Line
Interface. Once both users start the interactive prompt, they should be able to send money to the
other user and view their own balance.

#### Constraints

- Each user keeps track of their own balance
- Assume the users are on different computers
- State does not have to persist between sessions
- State should not be tracked or stored remotely (i.e. on a server)

```bash 
# Example Terminal Output
# Alice
# (plus connection options) 
$ ./start-peer --user=Alice
Welcome to your peering relationship!
> balance
0
> pay 10
Sent
> balance
-10 # (other balance is now 10) 
> exit
Goodbye.

# Bob
# (plus connection options) 
$ ./start-peer --user=Bob
# Welcome to your peering relationship!
> balance
0
You were paid 10!
> balance
10
> exit
Goodbye.
```



### Opinions


There are a couple of caveats that this implementation has.

1. No unit tests.
2. No pipeline to automate linting, runing tests, building the process.
2. No graceful cancelation of the entire process.
3. No way of informing other peers that somebody has leaved. 
4. This will not work if both peers are on different networks, example

```
[Peer] <- Router  -----   Router -> [Peer]
```

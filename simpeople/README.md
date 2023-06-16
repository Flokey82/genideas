# simpeople

This is a simple package to generate random people with random personalities. They are able to decide what to do from a list of pre-defined actions. The actions are based on their personality and the current situation, as well as the personality of the person they are interacting with. 

Actions might have a certain requirement or threshold to be considered. For example, a person might only steal from another person if they have a certain level of greed, or if they are hungry and/or do not like the other person.

- Step is to set up a random person with a random personality.
- Establish how compatible two people are... if they find each other likable.
- Select an action based of the predisposition of the person wrt the other person and the compatibility with the other person.

## TODO

This should use for example a GOAP (Goal Oriented Action Planning) system to decide what to do (any other suitable solution) where actions are weighed additionally based on the personality of the person and the situation. Right now, it's just a hardcoded list of actions. So we should calculate the utility as well as the tendency or desire to do something...
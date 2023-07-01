# Simpeople2

This package implements a sims-like simulation of people interacting with the world and each other.

NOTE: This implementation is loosely based on 
The Genius AI Behind The Sims - Game Maker's Toolkit
https://www.youtube.com/watch?v=9gf2MT-IOsg


## How does it work?

Like the original sims, the people have a set of needs (motives) that they try to satisfy. Each motive is represented as a number between -100 and 100 and have each a rate of decay, depending on what they represent. For example, sleep should require us to go to bed every 16 hours for 8 hours or so. 

Each motive value is then used to calculate a multiplier for the utility of actions that can be taken. For example, if we are very hungry, we should be more likely to eat something before we starve to death.
These multipliers are represented as a curve specific to each motive. Life threatening motives should have a very steep curve, while other motives should have a more shallow curve (like social interactions).

Personality traits are used as an additional weight to the utility of actions. A tidy person should value cleaning more than a slob, for example.

The utility of an action is calculated by (and this is speculation on my part), by adding the trait modifier (a value from 1 to 10) to the utility prior to multiplying it with the motive multiplier. This is then used to select an action from a list of possible actions.

Once all possible actions are ranked, one of the top actions is selected at random (possibly with a bias towards the top actions).

## How do we know where to do what to achieve what?

In the Sims, objects within the world space advertise their utility for certain actions. For example, a bed would advertise that it can be used to sleep in for a +50 to sleep. A toilet would advertise that it can be used to pee in for a +20 to bladder. This is also slightly modified by the distance to the object, so a bed that is far away would have a lower utility than a bed that is close by.

## Step by step implementation

- [X] Create motive type
    - [X] Decay
        - [ ] Move from ticks to delta time
    - [X] Multiplier
        - [X] Curves (linear, exponential, logarithmic, etc) for multipliers
        - [ ] Curve modifiers (steepness, etc)
    - [X] Implement two basic motives
        - [X] Food
        - [X] Sleep
        - [X] Bladder
        - [X] Hygiene
        - [ ] Social
        - [X] Fun
- [X] Create a person
    - [ ] Movement and pathfinding
    - [X] Add a list of motives
    - [X] Implement a basic action
        - [X] Eat
        - [X] Sleep
        - [X] Pee
        - [X] Shower
        - [X] Watch TV
        - [ ] Socialize
    - [X] Implement a basic action selection
        - [X] ... based on utility
        - [ ] ... with weighted randomization
    - [ ] Create personality type
        - [ ] Implement modifiers for traits
        - [ ] Implement a basic trait
            - [ ] Tidy
            - [ ] Slob
    - [ ] Implement a basic action selection (continued)
        - [ ] ... based on utility and personality
    - [X] Implement execution of action
        - [-] Actions should be interruptible
        - [ ] Actions should have a duration
- [X] Create object types
    - [X] Advertise utility for motives
    - [X] Implement a bed
    - [X] Implement a fridge
    - [X] Implement a toilet
    - [X] Implement a shower
    - [ ] Implement a couch
- [ ] Export to webp animation
- [ ] Tweak values!
# simsettlers

This is a WIP little procgen sim that simulates the growth of a small village. It is very crude at the moment and just operates on a square grid. In the future, building buildings will be tied to individuals and their resources, and plots can be expanded to include more than one tile. I just upload this current state since it does do something :)

![alt text](/simsettlers/images/test.png "a village around a river system")

## TODO

- [ ] Add a way to expand plots
- [X] Tie buildings to individual residents
- [-] Custom fitness functions for non-social individuals
    - [ ] Ranking system for purchasing a house for non-social individuals
- [X] Add a way to buy/sell houses
- [ ] Buildings can be upgraded or abandoned
- [ ] Add more building types
- [ ] Resource types
- [ ] Merge with simvillagesimple


## Notes to self

What do I want to achieve with this project?

I want to simulate the growth, and expansion of a village. 

Each settler should be able to independently decide on their actions:

- pick a hobby
- pick a job, change jobs (or be unemployed)
- gather resources, save money (or not)
- family
    - pick a partner, marry (or not)
    - have children (or not)
- realestate
    - build a house
    - buy/sell a house
    - reclaim/abandon a house (or not)
    - occupy a house / live in a house (or not)
    - maintain a house (or not)

Settlers should have:

- name
- age
- personality
- job
- hobby
- parents
- partner
- children
- home / residence (or not)
- money (or resources)
- inventory
- property (land, buildings)

Buying / selling a house:

Buying a house can be cheaper than building one, depending on the condition of the building. In any way, buying a house is also "cheaper" in terms of time. Building a house takes time, and the settler will be homeless (or will have to live with their parents) until the house is finished. Buying a house is instant, but the settler will have to pay for it.

Since a person can only have one home, we should consider selling excess real estate by default. Later we can introduce the concept of renting out a house.

Death:

If a person dies, all their property is inherited by their heir(s). Initially, we will consider the following relatives in order of inheritance:

- spouse
- children
- parents
- siblings
- grandparents
- uncles/aunts
- cousins

How do we let people plan?

I'd say, that every person can have a list of goals, which are ranked by priority. Depending on how many goals are fulfilled, the person will be happy or unhappy. Some people might not have any goals at all, and just live their life opportunistically. So if a person has a life goal of having a job, short term needs like buying food won't be a problem. Someone who doesn't want a job will do other things to get food, like begging, stealing, or growing their own food, which will be interesing to see.

I'd say we encode the personal life goals as bit flags in an uint32.
Some life goals might develop later in life, like finding a partner will start with puberty. Some life goals might end, like having children will end with old age.

If, we encode the goals as uint32, we could divide the bits into 4 categories:

- childhood goals (0-7)
    - socialize
- teenage goals (8-15)
    - hobby
    - education
    - skills
- adult goals (16-X)
    - job
    - partner
    - children
    - home
    - wealth
- old age goals (X->)
    - retirement
    - health
    - family
    - wealth

Maybe we also have one byte that encodes "vices", like addiction, laziness, insencereity, etc.

These goals can be altered through life events. For example, if a social child is traumatized, it might lose interest in socializing. If a child goes hungry when the parents are poor, it might develop a goal of wealth. 

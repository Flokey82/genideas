# gamestrategy

This is an experiment to create a simple (for now, zero-player) strategy game in order to experiment with different game mechanics and AI. Currently, this is just really bare bones and doesn't do much, but it's a start.

![alt text](/gamestrategy/images/expansion.webp "Screenshot")

## Game Mechanics

Here is the elevator pitch that ChatGPT came up with after much discussion and arguments:

Welcome to the medieval strategy game where you can experience the growth and conflicts of kingdoms! In this game, players will take on the role of rulers, strategically expanding their territories and engaging in wars. Let's explore the rules:

1. Game Setup:
   - Each player represents a kingdom and starts with a single territory.
   - The game is played on a grid divided into different cell types, such as meadows, forests, mountains, rivers, deserts, and coastal areas.
   - Each cell type has a cost associated with occupying and maintaining control over it.
   - Players can
      - ... expand their territories by occupying neighboring cells.     
      - ... only see the cell types of their own territories and the neighboring cells.
      - ... develop their territories by building and upgrading structures, such as farms, mines, barracks, walls, and markets which will provide various benefits like reduced costs, increased resource production, or improved defense.
      - ... recruit and train troops in their territories.
      - ... spend one turn and resources to explore ruins for a chance to find valuable artifacts.

2. Turn-based Gameplay:
   - The game progresses in turns, with each player taking actions one after another.
   - Players can perform a limited number of actions during their turn.

3. Actions:
   - Expanding Territories: Players can choose to attack or occupy neighboring territories to expand their kingdom. The success of an attack can depend on factors such as troop strength and terrain advantages.
   - Building and Upgrading: Players can use resources to construct buildings, such as farms, mines, barracks, walls, and markets, in their territories. Upgrading buildings can provide additional benefits, such as increased resource production or improved defense.
   - Resource Management: Players need to manage their resources effectively to sustain their kingdom and support their troops. Resources can be obtained from controlled territories or through trade with other players.
   - Diplomacy and Alliances: Players can negotiate with each other, form alliances, or engage in diplomacy. Alliances can provide mutual defense or trade benefits.

4. Territory Costs and Resources:
   - Each cell type has a cost associated with occupying and maintaining control over it.
   - Deserts are difficult to hold, so they have a higher cost per turn, but they provide limited resources in return.
   - Meadows are relatively cheap to occupy and can generate moderate resources per turn, making them valuable for sustaining a kingdom's economy.
   - Forests have a moderate cost and offer a bonus to resource production, such as timber or forage.
   - Mountains have a higher cost to occupy, but they may provide access to valuable resources like minerals or precious metals.
   - Rivers have a moderate cost and can provide a bonus to food production or act as trade routes for additional resources.
   - Coastal areas have a moderate cost and offer access to fishing grounds or trade routes, providing bonuses to food production or trade-related resources.

5. Troop Management and Victory Conditions:
   - Players can recruit and train troops in their territories.
   - The game can have various victory conditions, such as conquest (controlling the most territories), dominion (accumulating the most resources), or fame (achieving specific objectives).

By considering the costs and benefits associated with different cell types, players must carefully choose their expansion strategies, manage their resources efficiently, and adapt their tactics based on the terrain to succeed in the medieval world of kingdoms and wars.

## TODO

- [X] Add a simple map using noise for generating tile types
- [X] Add a simple player entity
   - [X] Mock up a simple AI
   - [ ] Add AI for the player entity
- [X] Add game loop
- [ ] Add more resources
- [ ] Add more buildings
- [ ] Flesh out the game mechanics
   - [ ] Improve game balance
- [X] Add webp animation export
- [ ] Add a simple GUI

https://www.gamedeveloper.com/design/designing-ai-algorithms-for-turn-based-strategy-games
https://catlikecoding.com/unity/tutorials/hex-map/
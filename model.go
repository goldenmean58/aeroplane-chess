//package areoplane_chess
package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

type Areoplane struct {
    index int
	team  int
	place int
}

type Player struct {
    index int
	aps       [4]Areoplane
	finishNum int
	isBot     bool
}

func makeDice() int64 {
	max := big.NewInt(6)
	n, _ := rand.Int(rand.Reader, max)
	m := n.Int64()
	m += 1
	return m
}

func (p *Areoplane) getGlobalLoc() (globalLoc int) {
	offset := p.team * 13
	if p.place > 51 || p.place <= 1 {
		globalLoc = -1
	} else {
		globalLoc = (p.place - 2 + offset) % 54
	}
	return
}

func (p *Areoplane) printStatus() {
	fmt.Printf("\tPlane %d\tlocal place %d\tglobal place %d\n", p.index, p.place, p.getGlobalLoc())
}

func (p *Player) printStatus() {
    fmt.Printf("Team %d:\n", p.index)
    for i := 0; i<4 ; i++ {
        p.aps[i].printStatus()
    }
	fmt.Printf("Finish: %d\n", p.finishNum)
}
func (p *Areoplane) move(step int) (err int) {
    return 0
}
func start() {
	fmt.Println("Game Start")
	fmt.Println("Init data")
	var players [4]Player
	var winner int
    score := make([]int, 0, 4)
	turn := 0
	for i := 0; i < 4; i++ {
		players[i].finishNum = 0
		players[i].isBot = false
        players[i].index = i
		for j := 0; j < 4; j++ {
			players[i].aps[j].team = i
			players[i].aps[j].place = 0
			players[i].aps[j].index = j
		}
	}
	fmt.Println("Init Over")
	for {
        if len(score) == 4 {
            fmt.Println("Match finished!")
            fmt.Println("Score:")
            for i := 0; i < 4; i++ {
                fmt.Printf("No.%d is Player %d\n", i+1, score[i])
            }
            break
        }
        if players[turn].finishNum == 4 {
            turn++
            if turn >=4 {
                turn = 0
            }
            continue
        }
		fmt.Printf("It's player %d's turn, Please input roll to roll dice! Type help to get help\n", turn)
		step := int(makeDice())
		if players[turn].isBot {
			// AI
			fmt.Printf("AI running\n")
		} else {
			var cmd string
			fmt.Scanf("%s", &cmd)
			switch cmd {
			case "cheat":
				fmt.Print("Step: ")
				fmt.Scanf("%d", &step)
				fallthrough
			case "roll", "r":
				fmt.Printf("Dice is %d\n", step)
				var moveableChess = make([]int, 0, 4)
				for i := 0; i < 4; i++ {
					if players[turn].aps[i].place > 0 {
						moveableChess = append(moveableChess, i)
					}
				}
				if step < 5 && len(moveableChess) == 0 {
					fmt.Println("No moveable plane, go next turn")
				} else {
					var moveId int
					for {
                        players[turn].printStatus()
						fmt.Print("Input plane's id to move:")
						fmt.Scanf("%d", &moveId)
                        status := players[turn].aps[moveId].move(step)
                        status += 0
						if moveId < 0 || moveId > 3 {
							fmt.Println("Please input correct id(0 to 3)")
							continue
						}
						if players[turn].aps[moveId].place < 0 {
							fmt.Println("That plane has finished!")
							continue
						}
						var is_moveable = false
						for i := 0; i < len(moveableChess); i++ {
							if moveableChess[i] == moveId {
								is_moveable = true
								break
							}
						}
						if is_moveable || step >= 5 {
							if players[turn].aps[moveId].place == 0 {
								players[turn].aps[moveId].place = 1
							} else {
								players[turn].aps[moveId].place += step
							}
							if players[turn].aps[moveId].place < 51 && players[turn].aps[moveId].place%4 == 3 {
								if players[turn].aps[moveId].place == 19 {
									players[turn].aps[moveId].place = 31
									fmt.Println("Long jump")
								} else {
									players[turn].aps[moveId].place += 4
									fmt.Println("Short jump")
								}
							}
							// kill
							// (local + offset) % one circle = global
							// local > one circle - 3 = enter final time
							curGL := players[turn].aps[moveId].getGlobalLoc()
							for i := 0; i < 4; i++ {
								if i == turn {
									continue
								}
								for j := 0; j < 4; j++ {
									tmpGL := players[i].aps[j].getGlobalLoc()
									if tmpGL >= 0 && curGL >= 0 && tmpGL == curGL {
										players[i].aps[j].place = 0
										fmt.Printf("A plane owned by team %d has been sent home\n", i)
									}
								}
							}
							if players[turn].aps[moveId].place >= 57 {
							    if players[turn].aps[moveId].place == 57 {
								    players[turn].aps[moveId].place = -1
								    players[turn].finishNum++
                                } else {
                                    players[turn].aps[moveId].place = 2 * 57 - players[turn].aps[moveId].place
                                }
                            }
							if step == 6 {
								fmt.Println("Dice is 6, gain an extra chance")
								turn--
								if turn < 0 {
									turn = 4
								}
							}
							break
						} else {
							fmt.Println("Failed to move that plane!")
						}
					}
				}
				turn++
				if turn >= 4 {
					turn = 0
				}
				break
			case "p":
				for i := 0; i < 4; i++ {
					players[i].printStatus()
				}
				break
			default:
				fmt.Printf("p for printing info\nr/roll for rolling the dice\n")
			}
		}
		for winner = 0; winner < 4; winner++ {
			if players[winner].finishNum == 4 {
	            fmt.Printf("Player %d finished the goal\n", winner)
			}
		}
	}
}

func main() {
	start()
}

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

var players [4]Player
var winner int
var score = make([]int, 0, 4)
var turn = 0

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
	fmt.Printf("\tPlane %d\tlocal place (%d / 57)\tglobal place %d\n", p.index, p.place, p.getGlobalLoc())
}

func (p *Player) printStatus() {
    fmt.Printf("Team %d:\n", p.index)
    for i := 0; i<4 ; i++ {
        p.aps[i].printStatus()
    }
	fmt.Printf("Finish: %d\n", p.finishNum)
}

func (p *Player) getMoveableChess() []int {
    var moveableChess = make([]int, 0, 4)
	for i := 0; i < 4; i++ {
		if p.aps[i].place > 0 {
			moveableChess = append(moveableChess, i)
		}
	}
    return moveableChess
}

func (p *Areoplane) isMoveable(moveableChess []int) bool {
    for i := 0; i < len(moveableChess); i++ {
        if moveableChess[i] == p.index {
            return true
        }
    }
    return false
}

func (p *Areoplane) predictMoveDis(moveableChess []int, step int) (int, bool) {
    //预测行动步长(主要是最终跳跃)
    var finalPos = p.place
    if finalPos < 0 {
        //fmt.Println("That plane has finished!")
        return 0, false
    }
    var is_moveable = p.isMoveable(moveableChess)
    if is_moveable || step >= 5 {
        if finalPos == 0 {
            finalPos = 1
        } else {
            finalPos += step
        }
        if finalPos < 51 && finalPos%4 == 3 {
            if finalPos == 19 {
                finalPos = 31
            } else {
                finalPos += 4
            }
        }
        if finalPos >= 57 {
            if finalPos == 57 {
                finalPos = -1
            } else {
                finalPos = 2 * 57 - finalPos
            }
        }
        return finalPos - p.place, true
    }
    return 0, false
}

func (p *Areoplane) kill(isPredict bool) int {
    count := 0
    curGL := p.getGlobalLoc()
    for i := 0; i < 4; i++ {
        if i == p.team {
            continue
        }
        for j := 0; j < 4; j++ {
            tmpGL := players[i].aps[j].getGlobalLoc()
            if tmpGL >= 0 && curGL >= 0 && tmpGL == curGL {
                count++
                if !isPredict {
                    players[i].aps[j].place = 0
                    fmt.Printf("Plane %d owned by team %d has been sent home\n", j, i)
                }
            }
        }
    }
    return count
}

func (p *Areoplane) move(moveableChess []int, step int) (err int) {
    dis, ok := p.predictMoveDis(moveableChess, step)
    if !ok {
        fmt.Println("Failed to move that plane!")
        return -1
    }
    p.place += dis
    if p.place == -1 {
        players[turn].finishNum++
    }
    p.kill(false)
    return 0
}

func initGame() {
	fmt.Println("Init data")
    for i := 0; i < 4; i++ {
		players[i].finishNum = 0
		players[i].isBot = true
        players[i].index = i
		for j := 0; j < 4; j++ {
			players[i].aps[j].team = i
			players[i].aps[j].place = 0
			players[i].aps[j].index = j
		}
	}
    players[0].isBot = false
	fmt.Println("Init Over")
}

func turnInc() {
    turn++
    if turn >= 4 {
        turn = 0
    }
}
func turnDec() {
    turn--
    if turn < 0 {
        turn = 3
    }
}

func (p *Player) doAI(moveableChess []int, step int) int {
    // 决策优先级:
    // 1.起飞
    // 2.踩人
    // 3.移动距离
    // 4.在终点区进入终点
    var i = 0
    var max = 0
    var maxIndex = -1
    //起飞
    if step >= 5 {
        for i=0;i<4;i++{
            if p.aps[i].place == 0 {
                p.aps[i].move(moveableChess, step)
                return i
            }
        }
    }
    // 踩
    for i=0;i<4;i++{
        oldPos := p.aps[i].place
        maxDis, ok := p.aps[i].predictMoveDis(moveableChess, step)
        if !ok {
            continue
        }
        p.aps[i].place+=maxDis
        killCount := p.aps[i].kill(true)
        if killCount > max {
            maxIndex = i
            max = killCount
        }
        p.aps[i].place = oldPos
    }
    if maxIndex >= 0 {
        p.aps[maxIndex].move(moveableChess, step)
        return maxIndex
    }
    //移动距离最大
    max = 0
    maxIndex = -1
    for i=0;i<4;i++{
        maxDis, ok :=p.aps[i].predictMoveDis(moveableChess, step)
        if p.aps[i].place > 51 && p.aps[i].place + maxDis == -1 && ok { //到达终点
            maxIndex = i
            break
        }
        if maxDis > max {
            maxIndex = i
            max = maxDis
        }
    }
    if maxIndex >= 0 {
        p.aps[maxIndex].move(moveableChess, step)
        return maxIndex
    }
    p.aps[moveableChess[0]].move(moveableChess, step)
    return moveableChess[0]
}

func start() {
	fmt.Println("Game Start")
    initGame()
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
            turnInc()
            continue
        }
		fmt.Printf("It's player %d's turn, Please input roll to roll dice! Type help to get help\n", turn)
		step := int(makeDice())
		var moveableChess = players[turn].getMoveableChess()
		if players[turn].isBot {
			// AI
			//fmt.Printf("AI running\n")
			fmt.Printf("Dice is %d\n", step)
            if step < 5 && len(moveableChess) == 0 {
				fmt.Println("No moveable plane, go next turn")
			    turnInc()
                continue
			}
            index := players[turn].doAI(moveableChess, step)
            fmt.Printf("Moved plane id: %d\n", index)
            players[turn].aps[index].printStatus()
            turnInc()
            if step == 6 {
                fmt.Println("Dice is 6, gain an extra chance")
                turnDec()
            }
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
				if step < 5 && len(moveableChess) == 0 {
					fmt.Println("No moveable plane, go next turn")
				    turnInc()
				} else {
                    if len(moveableChess) == 1 {
                        players[turn].aps[moveableChess[0]].move(moveableChess, step)
                        break
                    }
					var moveId int
					for {
                        players[turn].printStatus()
						fmt.Print("Input plane's id to move:")
						fmt.Scanf("%d", &moveId)
                        if moveId < 0 || moveId > 3 {
                            fmt.Println("Please input correct id(0 to 3)")
                            continue
                        }
                        status := players[turn].aps[moveId].move(moveableChess, step)
                        if status == 0 {
                            players[turn].aps[moveId].printStatus()
				            turnInc()
                            if step == 6 {
                                fmt.Println("Dice is 6, gain an extra chance")
                                turnDec()
                            }
                            break
                        }
                    }
				}
				break
			case "p":
				for i := 0; i < 4; i++ {
					players[i].printStatus()
				}
				break
            case "ai":
                fallthrough
            case "a":
                players[turn].isBot = true
			default:
				fmt.Printf("p for printing info\nr/roll for rolling the dice\nai for control shifting\n")
			}
		}
		for winner = 0; winner < 4; winner++ {
            for i:=0;i<len(score);i++{
                if score[i] == winner {
                    goto next
                }
            }
			if players[winner].finishNum == 4 {
	            fmt.Printf("Player %d finished the goal\n", winner)
                score = append(score, winner)
			}
            next:
		}
	}
}

func main() {
	start()
}

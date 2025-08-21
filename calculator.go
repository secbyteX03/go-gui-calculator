package main

import (
    "fmt"
    "strconv"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

var input string
var result *widget.Label

func main() {
    myApp := app.New()
    myWindow := myApp.NewWindow("Go Calculator")

    result = widget.NewLabel("0")

    buttons := []string{
        "7", "8", "9", "/",
        "4", "5", "6", "*",
        "1", "2", "3", "-",
        "0", ".", "=", "+",
        "C",
    }

    buttonWidgets := make([]fyne.CanvasObject, len(buttons))
    for i, label := range buttons {
        btn := widget.NewButton(label, func(l string) func() {
            return func() {
                handleInput(l)
            }
        }(label))
        buttonWidgets[i] = btn
    }

    grid := container.NewGridWithColumns(4, buttonWidgets...)
    content := container.NewVBox(result, grid)

    myWindow.SetContent(content)
    myWindow.Resize(fyne.NewSize(250, 300))
    myWindow.ShowAndRun()
}

func handleInput(value string) {
    if value == "C" {
        input = ""
        result.SetText("0")
        return
    }

    if value == "=" {
        eval, err := evalExpression(input)
        if err != nil {
            result.SetText("Error")
        } else {
            result.SetText(fmt.Sprintf("%v", eval))
            input = fmt.Sprintf("%v", eval)
        }
        return
    }

    input += value
    result.SetText(input)
}

// A simple evaluator (supports +, -, *, /)
func evalExpression(expr string) (float64, error) {
    var num float64
    var lastOp byte = '+'
    var current string
    stack := []float64{}

    for i := 0; i < len(expr); i++ {
        ch := expr[i]
        if (ch >= '0' && ch <= '9') || ch == '.' {
            current += string(ch)
        }
        if ch < '0' || ch > '9' || i == len(expr)-1 {
            if current != "" {
                n, _ := strconv.ParseFloat(current, 64)
                switch lastOp {
                case '+':
                    stack = append(stack, n)
                case '-':
                    stack = append(stack, -n)
                case '*':
                    stack[len(stack)-1] *= n
                case '/':
                    stack[len(stack)-1] /= n
                }
                current = ""
            }
            if ch == '+' || ch == '-' || ch == '*' || ch == '/' {
                lastOp = ch
            }
        }
    }

    // sum stack
    sum := 0.0
    for _, v := range stack {
        sum += v
    }
    return sum, nil
}

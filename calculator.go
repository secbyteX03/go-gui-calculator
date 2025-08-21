package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Calculator struct {
	app        fyne.App
	window     fyne.Window
	input      string
	result     *widget.Label
	history    []string
	historyBox *widget.List
	memory     float64
	themeDark  bool
}

func NewCalculator() *Calculator {
	c := &Calculator{
		app:    app.New(),
		memory: 0,
	}
	c.window = c.app.NewWindow("Go Calculator")
	c.window.Resize(fyne.NewSize(350, 500))
	c.themeDark = true
	c.app.Settings().SetTheme(theme.DarkTheme())

	return c
}

func (c *Calculator) buildUI() {
	// Result display
	c.result = widget.NewLabel("0")
	c.result.TextStyle = fyne.TextStyle{Monospace: true}
	c.result.Alignment = fyne.TextAlignTrailing

	// History list
	c.historyBox = widget.NewList(
		func() int { return len(c.history) },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(c.history[i])
		},
	)

	// Buttons
	buttons := []struct {
		label    string
		action   func()
		size     fyne.Size
		style    fyne.TextStyle
		icon      fyne.Resource
		disabled bool
	}{
		{label: "MC", action: c.clearMemory},
		{label: "MR", action: c.recallMemory},
		{label: "M-", action: c.subtractFromMemory},
		{label: "M+", action: c.addToMemory, style: fyne.TextStyle{Bold: true}},
		{label: "√", action: func() { c.appendToInput("sqrt(") }},
		{label: "(", action: func() { c.appendToInput("(") }},
		{label: ")", action: func() { c.appendToInput(")") }},
		{label: "C", action: c.clearInput, style: fyne.TextStyle{Bold: true}},
		{label: "7", action: func() { c.appendToInput("7") }},
		{label: "8", action: func() { c.appendToInput("8") }},
		{label: "9", action: func() { c.appendToInput("9") }},
		{label: "÷", action: func() { c.appendToInput("/") }},
		{label: "4", action: func() { c.appendToInput("4") }},
		{label: "5", action: func() { c.appendToInput("5") }},
		{label: "6", action: func() { c.appendToInput("6") }},
		{label: "×", action: func() { c.appendToInput("*") }},
		{label: "1", action: func() { c.appendToInput("1") }},
		{label: "2", action: func() { c.appendToInput("2") }},
		{label: "3", action: func() { c.appendToInput("3") }},
		{label: "-", action: func() { c.appendToInput("-") }},
		{label: "0", action: func() { c.appendToInput("0") }},
		{label: ".", action: func() { c.appendToInput(".") }},
		{label: "±", action: c.toggleSign},
		{label: "+", action: func() { c.appendToInput("+") }},
		{label: "%", action: c.percentage, style: fyne.TextStyle{Bold: true}},
		{label: "⌫", action: c.backspace, style: fyne.TextStyle{Bold: true}},
		{label: "=", action: c.calculate, style: fyne.TextStyle{Bold: true}},
	}

	// Create buttons
	buttonGrid := container.NewGridWithColumns(4)
	for _, btn := range buttons {
		button := widget.NewButton(btn.label, btn.action)
		button.Importance = widget.HighImportance
		if btn.style.Bold {
			button.TextStyle = fyne.TextStyle{Bold: true}
		}
		buttonGrid.Add(button)
	}

	// Theme toggle
	themeToggle := widget.NewButton("🌙/☀️", c.toggleTheme)

	// Layout
	historyContainer := container.NewVBox(
		widget.NewLabel("History:"),
		container.NewVScroll(c.historyBox),
	)

	mainContent := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), themeToggle),
		c.result,
		buttonGrid,
	)

	split := container.NewHSplit(
		mainContent,
		historyContainer,
	)
	split.Offset = 0.7

	c.window.SetContent(split)
}

func (c *Calculator) Run() {
	c.buildUI()
	c.window.ShowAndRun()
}

func (c *Calculator) updateDisplay() {
	if c.input == "" {
		c.result.SetText("0")
	} else {
		c.result.SetText(c.input)
	}
}

func (c *Calculator) appendToInput(s string) {
	c.input += s
	c.updateDisplay()
}

func (c *Calculator) clearInput() {
	c.input = ""
	c.updateDisplay()
}

func (c *Calculator) backspace() {
	if len(c.input) > 0 {
		c.input = c.input[:len(c.input)-1]
		c.updateDisplay()
	}
}

func (c *Calculator) toggleSign() {
	if c.input != "" && c.input[0] == '-' {
		c.input = c.input[1:]
	} else if c.input != "" {
		c.input = "-" + c.input
	}
	c.updateDisplay()
}

func (c *Calculator) percentage() {
	result, err := c.evalExpression(c.input)
	if err == nil {
		c.input = fmt.Sprintf("%g", result/100)
		c.updateDisplay()
	}
}

func (c *Calculator) calculate() {
	if c.input == "" {
		return
	}

	result, err := c.evalExpression(c.input)
	if err != nil {
		dialog.ShowError(err, c.window)
		return
	}

	// Add to history
	historyEntry := fmt.Sprintf("%s = %g", c.input, result)
	c.history = append([]string{historyEntry}, c.history...)
	c.historyBox.Refresh()

	c.input = fmt.Sprintf("%g", result)
	c.updateDisplay()
}

func (c *Calculator) addToHistory(entry string) {
	c.history = append([]string{entry}, c.history...)
	// Keep only the last 50 history entries
	if len(c.history) > 50 {
		c.history = c.history[:50]
	}
	c.historyBox.Refresh()
}

func (c *Calculator) evalExpression(expr string) (float64, error) {
	// Replace × and ÷ with * and /
	expr = strings.ReplaceAll(expr, "×", "*")
	expr = strings.ReplaceAll(expr, "÷", "/")

	// Handle square root
	expr = strings.ReplaceAll(expr, "sqrt(", "math.Sqrt(")

	// Evaluate the expression using Go's parser
	result, err := evaluateExpression(expr)
	if err != nil {
		return 0, fmt.Errorf("invalid expression")
	}

	return result, nil
}

func evaluateExpression(expr string) (float64, error) {
	// This is a simplified evaluator. In a production app, you'd want to use
	// a proper expression evaluator or implement a more robust solution.
	// For now, we'll use a simple approach that works for basic expressions.

	// Handle empty expression
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return 0, nil
	}

	// Try to parse as a number first
	if val, err := strconv.ParseFloat(expr, 64); err == nil {
		return val, nil
	}

	// Handle parentheses
	if strings.Contains(expr, "(") || strings.Contains(expr, ")") {
		return evaluateWithParentheses(expr)
	}

	// Handle operators in order of precedence
	operators := []struct {
		symbol string
		apply  func(a, b float64) float64
	}{
		{"*", func(a, b float64) float64 { return a * b }},
		{"/", func(a, b float64) float64 {
			if b == 0 {
				return math.NaN()
			}
			return a / b
		}},
		{"+", func(a, b float64) float64 { return a + b }},
		{"-", func(a, b float64) float64 { return a - b }},
	}

	for _, op := range operators {
		if i := strings.LastIndex(expr, op.symbol); i > 0 {
			left, err1 := evaluateExpression(expr[:i])
			right, err2 := evaluateExpression(expr[i+1:])

			if err1 != nil || err2 != nil {
				return 0, fmt.Errorf("invalid expression")
			}

			result := op.apply(left, right)
			if math.IsNaN(result) {
				return 0, fmt.Errorf("division by zero")
			}
			return result, nil
		}
	}

	return 0, fmt.Errorf("invalid expression")
}

func evaluateWithParentheses(expr string) (float64, error) {
	var stack []int
	pairs := make(map[int]int)

	// Find matching parentheses
	for i, r := range expr {
		if r == '(' {
			stack = append(stack, i)
		} else if r == ')' {
			if len(stack) == 0 {
				return 0, fmt.Errorf("mismatched parentheses")
			}
			pairs[stack[len(stack)-1]] = i
			stack = stack[:len(stack)-1]
		}
	}

	if len(stack) > 0 {
		return 0, fmt.Errorf("mismatched parentheses")
	}

	// If no parentheses, evaluate normally
	if len(pairs) == 0 {
		return evaluateExpression(expr)
	}

	// Find innermost parentheses
	var start, end int
	for s, e := range pairs {
		// Check if these parentheses are inside any other parentheses
		nested := false
		for s2, e2 := range pairs {
			if s2 < s && e2 > e {
				nested = true
				break
			}
		}
		if !nested {
			start, end = s, e
			break
		}
	}

	// Evaluate the innermost expression
	innerResult, err := evaluateExpression(expr[start+1 : end])
	if err != nil {
		return 0, err
	}

	// Replace the parenthesized expression with its result
	newExpr := expr[:start] + fmt.Sprintf("%g", innerResult) + expr[end+1:]

	// Evaluate the new expression
	return evaluateExpression(newExpr)
}

func (c *Calculator) addToMemory() {
	if c.input == "" {
		return
	}
	val, err := strconv.ParseFloat(c.input, 64)
	if err == nil {
		c.memory += val
	}
}

func (c *Calculator) subtractFromMemory() {
	if c.input == "" {
		return
	}
	val, err := strconv.ParseFloat(c.input, 64)
	if err == nil {
		c.memory -= val
	}
}

func (c *Calculator) recallMemory() {
	c.input = fmt.Sprintf("%g", c.memory)
	c.updateDisplay()
}

func (c *Calculator) clearMemory() {
	c.memory = 0
}

func (c *Calculator) toggleTheme() {
	c.themeDark = !c.themeDark
	if c.themeDark {
		c.app.Settings().SetTheme(theme.DarkTheme())
	} else {
		c.app.Settings().SetTheme(theme.LightTheme())
	}
}

func main() {
	calc := NewCalculator()
	calc.Run()
}

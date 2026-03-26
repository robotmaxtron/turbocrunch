Output demo.gif

Set Shell "bash"
Set FontSize 16
Set Width 1200
Set Height 800
Set Padding 20
# Set WindowBar Fixed
Set Margin 20
Set BorderRadius 10
Set TypingSpeed 100ms

# Start the application
Type "./turbocrunch"
Enter
Sleep 1s

# 1. Expression Evaluation
Type "1+1"
Enter
Sleep 500ms
Type "ans * 3"
Enter
Sleep 1s

# 2. Angle Mode Toggle & Context Usage
# Define a variable for context
Type "a = 90"
Enter
Sleep 500ms
# Sin(a) in Rad (default)
Type "sin(a)"
Enter
Sleep 500ms
# Toggle to Deg (Ctrl+A)
Ctrl+A
Sleep 500ms
# Re-evaluate using ans or variable
Type "sin(a)"
Enter
Sleep 1s

# 3. Formula Book (Ctrl+F)
# First define needed variables
Type "r=5"
Enter
Sleep 500ms
Ctrl+F
Sleep 1s
Type "circle"
Sleep 500ms
Enter
Sleep 500ms
# Formula pi * r^2 is inserted, evaluate it
Enter
Sleep 1s

# 4. Constants (Ctrl+C)
Ctrl+C
Sleep 1s
Type "speed of light"
Sleep 500ms
Enter
Sleep 500ms
# Constant is inserted, evaluate it
Enter
Sleep 1s

# 5. Units (Ctrl+U)
Type "100"
Sleep 500ms
Ctrl+U
Sleep 1s
Type "meter"
Sleep 500ms
Enter
Sleep 500ms
# Unit is inserted, evaluate it
Enter
Sleep 1s

# 6. History Navigation (Ctrl+H)
Ctrl+H
Sleep 1s
# Move up to 1+1 (it's at the top of history table)
Up
Sleep 500ms
Up
Sleep 500ms
Up
Sleep 500ms
# Select it
Enter
Sleep 1s
# Re-evaluate it
Enter
Sleep 1s

# 7. Help Menu (?)
Type "?"
Sleep 2s
Type "?"
Sleep 500ms

# 8. Quit (Ctrl+Q)
Ctrl+Q
Sleep 1s

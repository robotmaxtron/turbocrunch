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

# 2. Angle Mode Toggle
# Sin(90) in Rad (default)
Type "sin(90)"
Enter
Sleep 500ms
# Toggle to Deg (Ctrl+A)
Ctrl+A
Sleep 500ms
Type "sin(90)"
Enter
Sleep 1s

# 3. Formula Book (Ctrl+F)
Ctrl+F
Sleep 1s
Type "circle"
Sleep 500ms
Enter
Sleep 500ms
# Formula is inserted, evaluate it (might need variable values)
# For simplicity, just browse and go back
Type "r=5"
Enter
Sleep 500ms

# 4. Constants (Ctrl+C)
Ctrl+C
Sleep 1s
Type "speed of light"
Sleep 500ms
Enter
Sleep 500ms
Enter
Sleep 1s

# 5. Units (Ctrl+U)
Ctrl+U
Sleep 1s
Type "meter"
Sleep 500ms
Enter
Sleep 500ms
Enter
Sleep 1s

# 6. History Navigation (Ctrl+H)
Ctrl+H
Sleep 1s
Down
Sleep 500ms
Down
Sleep 500ms
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

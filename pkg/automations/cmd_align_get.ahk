#Requires AutoHotkey v2.0+
Click(500, 975)
Sleep(2000)
TypeTextSlowly("/align get")
Sleep(1000)
Send("{Enter}")
Sleep(1000)
Send("{Enter}")

RandomDelay(min := 70, max := 90) {
    return Random(min, max)
}

TypeTextSlowly(text) {
    for char in StrSplit(text) {
        SendText(char)
        Sleep(RandomDelay())
    }
}

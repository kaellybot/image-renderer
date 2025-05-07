#Requires AutoHotkey v2.0+
SetTitleMatchMode(2)

; Step 1: Activate Discord window
if WinExist("Discord")
{
    WinActivate("Discord")
    WinWaitActive("Discord")
}
else
{
    MsgBox("Discord not found!")
    ExitApp()
}

; Click server tutorial
Click(40, 250)

; Click tutorial channel
Click(150, 220)
Sleep(1000)
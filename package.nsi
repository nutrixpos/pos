!include "MUI2.nsh"
!include "nsDialogs.nsh"

; Basic setup
Name "Nutrix POS"
OutFile "nutrixpos-installer.exe"
Unicode True

; Default installation folder
InstallDir "$LocalAppData\NutrixPOS"
InstallDirRegKey HKLM "Software\Nutrixpos" "Install_Dir"
RequestExecutionLevel admin

; Variables
Var StartMenuFolder

; Interface Settings
!define MUI_ABORTWARNING
!define MUI_ICON "icon.ico"
!define MUI_UNICON "icon.ico"

; Installer pages
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_LICENSE "license.txt"
!insertmacro MUI_PAGE_COMPONENTS
!insertmacro MUI_PAGE_DIRECTORY

!define MUI_STARTMENUPAGE_REGISTRY_ROOT "HKLM"
!define MUI_STARTMENUPAGE_REGISTRY_KEY "Software\Nutrixpos"
!define MUI_STARTMENUPAGE_REGISTRY_VALUENAME "Start Menu Folder"
!insertmacro MUI_PAGE_STARTMENU Application $StartMenuFolder

!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

; Languages
!insertmacro MUI_LANGUAGE "English"

; ---------------------------
; Installer Sections
; ---------------------------

Section "Main Application" SecMain
    SectionIn RO
    SetOutPath "$LocalAppData\NutrixPOS"
    File "nutrixpos64.exe"
    File "config.yaml"
    File /r "assets"
    File "icon.ico"
    
    CreateDirectory "$LocalAppData\NutrixPOS"
    CreateDirectory "$LocalAppData\NutrixPOS\mnt"
    CreateDirectory "$LocalAppData\NutrixPOS\mnt\public"
    CreateDirectory "$LocalAppData\NutrixPOS\mnt\frontend"

    SetOutPath "$LocalAppData\NutrixPOS\mnt\frontend"
    File /r "frontend\*.*"
    
    SetOutPath "$LocalAppData\NutrixPOS"

    WriteRegStr HKLM "Software\Nutrixpos" "Install_Dir" "$LocalAppData\NutrixPOS"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "DisplayName" "Nutrix POS"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "UninstallString" '"$LocalAppData\NutrixPOS\uninstall.exe"'
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "DisplayIcon" "$LocalAppData\NutrixPOS\icon.ico"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "Publisher" "Amr Elmawardy"
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "NoModify" 1
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "NoRepair" 1

    WriteUninstaller "$LocalAppData\NutrixPOS\uninstall.exe"

    !insertmacro MUI_STARTMENU_WRITE_BEGIN Application
        CreateDirectory "$SMPROGRAMS\$StartMenuFolder"
        CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Nutrixpos.lnk" "$LocalAppData\NutrixPOS\nutrixpos64.exe" "" "$LocalAppData\NutrixPOS\icon.ico" 0
        CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Uninstall.lnk" "$LocalAppData\NutrixPOS\uninstall.exe"
    !insertmacro MUI_STARTMENU_WRITE_END
SectionEnd

Section "Desktop Shortcut"
    CreateShortcut "$DESKTOP\Nutrixpos.lnk" "$LocalAppData\NutrixPOS\nutrixpos64.exe" "" "$LocalAppData\NutrixPOS\icon.ico" 0
SectionEnd


Section "Uninstall"
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos"
    DeleteRegKey HKLM "Software\Nutrixpos"

    Delete "$LocalAppData\NutrixPOS\license.txt"
    Delete "$LocalAppData\NutrixPOS\nutrixpos64.exe"
    Delete "$LocalAppData\NutrixPOS\nssm.exe"
    Delete "$LocalAppData\NutrixPOS\uninstall.exe"
    Delete "$LocalAppData\NutrixPOS\icon.ico"

    Delete "$LocalAppData\NutrixPOS\mnt\frontend\*.*"
    RMDir "$LocalAppData\NutrixPOS\mnt\frontend"
    RMDir "$LocalAppData\NutrixPOS\mnt\public"
    RMDir "$LocalAppData\NutrixPOS\mnt"
    RMDir /r "$LocalAppData\NutrixPOS\assets"
    RMDir "$LocalAppData\NutrixPOS"
    Delete "$LocalAppData\NutrixPOS\config.yaml"
    RMDir "$LocalAppData\NutrixPOS"

    ; Remove start menu shortcuts
    !insertmacro MUI_STARTMENU_GETFOLDER Application $StartMenuFolder
    Delete "$SMPROGRAMS\$StartMenuFolder\Nutrixpos.lnk"
    Delete "$SMPROGRAMS\$StartMenuFolder\Uninstall.lnk"
    RMDir "$SMPROGRAMS\$StartMenuFolder"

    ; Remove desktop shortcut
    Delete "$DESKTOP\Nutrixpos.lnk"

SectionEnd

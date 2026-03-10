!include "MUI2.nsh"
!include "nsDialogs.nsh"

; Basic setup
Name "Nutrix POS"
OutFile "nutrixpos-installer.exe"
Unicode True

; Default installation folder
InstallDir "$PROGRAMFILES\Nutrixpos"
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
    SetOutPath "$INSTDIR"
    File "nutrixpos64.exe"
    File "config.yaml"
    File /r "assets"
    File "icon.ico"
    
    CreateDirectory "$INSTDIR\mnt"
    CreateDirectory "$INSTDIR\mnt\public"
    CreateDirectory "$INSTDIR\mnt\frontend"
    CreateDirectory "$LocalAppData\NutrixPOS"

    SetOutPath "$INSTDIR\mnt\frontend"
    File /r "frontend\*.*"
    
    SetOutPath "$INSTDIR"

    WriteRegStr HKLM "Software\Nutrixpos" "Install_Dir" "$INSTDIR"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "DisplayName" "Nutrix POS"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "UninstallString" '"$INSTDIR\uninstall.exe"'
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "DisplayIcon" "$INSTDIR\icon.ico"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "Publisher" "Amr Elmawardy"
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "NoModify" 1
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "NoRepair" 1

    WriteUninstaller "$INSTDIR\uninstall.exe"

    !insertmacro MUI_STARTMENU_WRITE_BEGIN Application
        CreateDirectory "$SMPROGRAMS\$StartMenuFolder"
        CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Nutrixpos.lnk" "$INSTDIR\nutrixpos64.exe" "" "$INSTDIR\icon.ico" 0
        CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Uninstall.lnk" "$INSTDIR\uninstall.exe"
    !insertmacro MUI_STARTMENU_WRITE_END
SectionEnd


Section "Uninstall"
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos"
    DeleteRegKey HKLM "Software\Nutrixpos"

    Delete "$INSTDIR\license.txt"
    Delete "$INSTDIR\nutrixpos64.exe"
    Delete "$INSTDIR\nssm.exe"
    Delete "$INSTDIR\config.yaml"
    Delete "$INSTDIR\uninstall.exe"
    Delete "$INSTDIR\icon.ico"

    Delete "$INSTDIR\mnt\frontend\*.*"
    RMDir "$INSTDIR\mnt\frontend"
    RMDir "$INSTDIR\mnt\public"
    RMDir "$INSTDIR\mnt"
    RMDir /r "$INSTDIR\assets"
    RMDir "$INSTDIR"

    RMDir /r "$LocalAppData\NutrixPOS\uploads"
    RMDir "$LocalAppData\NutrixPOS"
    RMDir "$LocalAppData"

    ; Remove start menu shortcuts
    !insertmacro MUI_STARTMENU_GETFOLDER Application $StartMenuFolder
    Delete "$SMPROGRAMS\$StartMenuFolder\Nutrixpos.lnk"
    Delete "$SMPROGRAMS\$StartMenuFolder\Uninstall.lnk"
    RMDir "$SMPROGRAMS\$StartMenuFolder"

SectionEnd

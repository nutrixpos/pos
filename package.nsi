!include "MUI2.nsh"
!include "nsDialogs.nsh"

; Basic setup
Name "Nutrix POS"
OutFile "nutrixpos-installer.msi"
Unicode True

; Default installation folder
InstallDir "$PROGRAMFILES\Nutrixpos"
InstallDirRegKey HKLM "Software\Nutrixpos" "Install_Dir"
RequestExecutionLevel admin

; Variables
Var StartMenuFolder
Var UninstallMongoDBChecked
Var InstallMongoDBChecked

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

; Add custom options page with MongoDB checkbox
Page custom InstallOptionsPage InstallOptionsPageLeave

!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_CONFIRM
UninstPage custom un.UninstallOptionsPage un.UninstallOptionsPageLeave
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
    File "nutrixpos.exe"
    File "config.yaml"
    
    CreateDirectory "$INSTDIR\mnt"
    CreateDirectory "$INSTDIR\mnt\public"
    CreateDirectory "$INSTDIR\mnt\frontend"

    SetOutPath "$INSTDIR\mnt\frontend"
    File /r "frontend\*.*"
    
    SetOutPath "$INSTDIR"

    WriteRegStr HKLM "Software\Nutrixpos" "Install_Dir" "$INSTDIR"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "DisplayName" "Nutrix POS"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "UninstallString" '"$INSTDIR\uninstall.exe"'
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "DisplayIcon" "$INSTDIR\icon.ico"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "Publisher" "Elmawardy"
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "NoModify" 1
    WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos" "NoRepair" 1

    WriteUninstaller "$INSTDIR\uninstall.exe"

    !insertmacro MUI_STARTMENU_WRITE_BEGIN Application
        CreateDirectory "$SMPROGRAMS\$StartMenuFolder"
        CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Nutrixpos.lnk" "$INSTDIR\nutrixpos.exe"
        CreateShortcut "$SMPROGRAMS\$StartMenuFolder\Uninstall.lnk" "$INSTDIR\uninstall.exe"
    !insertmacro MUI_STARTMENU_WRITE_END
SectionEnd

Section "Install nutrixpos windows service"
    SetOutPath "$INSTDIR"
    File "nutrixpos.exe"
    File "nssm.exe"

    ;nsExec::ExecToLog 'sc create NutrixPosService binPath= "$INSTDIR\nutrixpos.exe" DisplayName= "Nutrix POS" start= auto obj= "NT AUTHORITY\NetworkService"'
    ;nsExec::ExecToLog 'sc start NutrixPosService'
    ; Install your app as a service via NSSM under NetworkService
    nsExec::ExecToLog '"$INSTDIR\nssm.exe" install NutrixPosService "$INSTDIR\nutrixpos.exe"'
    
    ; Set display name
    nsExec::ExecToLog '"$INSTDIR\nssm.exe" set NutrixPosService DisplayName "NutrixPOS Service"'


    nsExec::ExecToLog '"$INSTDIR\nssm.exe" set NutrixPosService ObjectName "NT AUTHORITY\NetworkService"'
    nsExec::ExecToLog 'sc start NutrixPosService'

    WriteUninstaller "$INSTDIR\uninstall.exe"
SectionEnd

Section "-Conditional MongoDB Install"
    ${If} $InstallMongoDBChecked == ${BST_CHECKED}
        Call InstallMongoDB
    ${EndIf}
SectionEnd

Function InstallMongoDB
    ; Extract MongoDB MSI to temp folder
    SetOutPath "$TEMP\MongoDB"
    File "mongodb-windows-x86_64-8.2.1-signed.msi"

    ; Install MongoDB silently from temp
    ExecWait 'msiexec /i "$TEMP\MongoDB\mongodb-windows-x86_64-8.2.1-signed.msi" INSTALLLOCATION="C:\MongoDB" ADDLOCAL=All'

    ; Clean up extracted MSI
    Delete "$TEMP\MongoDB\mongodb-windows-x86_64-8.2.1-signed.msi"
    RMDir "$TEMP\MongoDB"
FunctionEnd

; ---------------------------
; Install Options Page (with MongoDB checkbox)
; ---------------------------

Function InstallOptionsPage
    nsDialogs::Create 1018
    Pop $0
    ${If} $0 == error
        Abort
    ${EndIf}

    ${NSD_CreateCheckbox} 20u 20u 200u 12u "Install MongoDB"
    Pop $InstallMongoDBChecked
    ${NSD_SetState} $InstallMongoDBChecked ${BST_CHECKED}

    nsDialogs::Show
FunctionEnd

Function InstallOptionsPageLeave
    ${NSD_GetState} $InstallMongoDBChecked $0
    StrCpy $InstallMongoDBChecked $0
FunctionEnd

; ---------------------------
; Uninstall Custom Page
; ---------------------------

Function un.UninstallOptionsPage
    nsDialogs::Create 1018
    Pop $0
    ${If} $0 == error
        Abort
    ${EndIf}

    ${NSD_CreateCheckbox} 20u 20u 200u 12u "Uninstall MongoDB"
    Pop $UninstallMongoDBChecked
    ${NSD_SetState} $UninstallMongoDBChecked ${BST_CHECKED}

    nsDialogs::Show
FunctionEnd

Function un.UninstallOptionsPageLeave
    ${NSD_GetState} $UninstallMongoDBChecked $0
    StrCpy $UninstallMongoDBChecked $0
FunctionEnd

; ---------------------------
; Uninstaller Section
; ---------------------------

Section "Open frontend in browser"
  ; 🔥 Open browser after install
  ExecShell "open" "http://localhost:8080"
SectionEnd

Section "Uninstall"
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\Nutrixpos"
    DeleteRegKey HKLM "Software\Nutrixpos"

    nsExec::ExecToLog 'sc stop NutrixPosService'
    nsExec::ExecToLog 'sc delete NutrixPosService'

    Delete "$INSTDIR\license.txt"
    Delete "$INSTDIR\nutrixpos.exe"
    Delete "$INSTDIR\nssm.exe"
    Delete "$INSTDIR\config.yaml"
    Delete "$INSTDIR\uninstall.exe"

    Delete "$INSTDIR\mnt\frontend\*.*"
    RMDir "$INSTDIR\mnt\frontend"
    RMDir "$INSTDIR\mnt\public"
    RMDir "$INSTDIR\mnt"
    RMDir "$INSTDIR"

    ; Remove start menu shortcuts
    !insertmacro MUI_STARTMENU_GETFOLDER Application $StartMenuFolder
    Delete "$SMPROGRAMS\$StartMenuFolder\Nutrixpos.lnk"
    Delete "$SMPROGRAMS\$StartMenuFolder\Uninstall.lnk"
    RMDir "$SMPROGRAMS\$StartMenuFolder"

    ; MongoDB uninstall step
    ${If} $UninstallMongoDBChecked == ${BST_CHECKED}
        nsExec::ExecToLog 'sc stop MongoDB'
        nsExec::ExecToLog 'msiexec /x {DA66F0D9-2B0F-4DCB-BBA8-E540E020B162} /qn'

        RMDir /r "C:\MongoDB\data\db"
        RMDir /r "C:\MongoDB"
    ${EndIf}
SectionEnd

This program allows control of the original version of the HiTecAstro Mount Hub Pro - it looks like this:

<img src="https://i.imgur.com/6VYarDZ.jpeg" alt="MHP">

It creates an ASCOM Alpaca driver that can be selected from the Switch and Focuser tabs in <a href="https://nighttime-imaging.eu/">N.I.N.A.</a>. It allows for control of the 4 dew heaters, 8 switched power ports and the stepper motor controller.

Commands can then be sent from the N.I.N.A. to the Mount Hub Pro.

Setting can be customised in the settings.json file which is created when the program is first run. Focucer speed defaults to 50.

Example screen prints from N.I.N.A.

<img src="https://raw.githubusercontent.com/exploded/mhp-ascom-alpaca/refs/heads/main/NINASwitch.jpg" alt="Switch">


<img src="https://raw.githubusercontent.com/exploded/mhp-ascom-alpaca/refs/heads/main/NINAfocuser.jpg" alt="Focuser">

## Instructions
Plug the MHP into a Windows computer and it will be recognised as a USB HID device - no drivers are needed. Download and run the mhp.exe executable on the same computer the MHP is plugged into. Run N.I.N.A and it should be able to find and connect to the MHP as a switch device and a focuser device.
